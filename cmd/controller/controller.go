package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	check "github.com/Kurlabs/alerty/internal/check"
	event "github.com/Kurlabs/alerty/internal/event"
	conn "github.com/Kurlabs/alerty/shared/mongo"
	message "github.com/Kurlabs/alerty/shared/pubsub"

	"cloud.google.com/go/pubsub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	METRICRESPONSE   = "response"
	GTE              = ">="
	LT               = "<"
	DOWNTIME         = "downtime"
	UPTIME           = "uptime"
	EMAILCONTACT     = "Contact.EmailContact"
	PHONECONTACT     = "Contact.PhoneContact"
	SLACKINTEGRATION = "Integration.SlackIntegration"
	WEBSITEMONITOR   = "Monitor.WebsiteMonitor"
	SOCKETMONITOR    = "Monitor.SocketMonitor"
	APP              = "app"
	EMAIL            = "email"
	SMS              = "sms"
	SLACK            = "slack"
	INFO             = "info"
	WARNING          = "warning"
	CRITICAL         = "critical"
	DONE             = "done"
)

var pbClient *pubsub.Client

// ChannelsList contains Contacts array field to handle json parsing
type ChannelsList struct {
	EmailContacts     []event.Contact
	PhoneContacts     []event.Contact
	SlackIntegrations []event.Integration
}

func sendMessage(eventType string, monitor check.Monitor, sms, email, slack bool, emailContacts, phoneContacts []event.Contact, slackIntegrations []event.Integration) {
	log.Println(monitor.Name, ": Sending pubsub message!")
	objectTypes := map[string]string{
		WEBSITEMONITOR: "website",
		SOCKETMONITOR:  "socket",
	}
	attrs := map[string]string{
		"event":       eventType,
		"object":      monitor.URL,
		"object_type": objectTypes[monitor.CLS],
	}
	if sms {
		attrs["sms"] = "true"
	}
	if email {
		attrs["email"] = "true"
	}
	if slack {
		attrs["slack"] = "true"
	}
	channelsData := ChannelsList{
		EmailContacts:     emailContacts,
		PhoneContacts:     phoneContacts,
		SlackIntegrations: slackIntegrations,
	}
	data, err := json.Marshal(channelsData)
	if err != nil {
		log.Panic(err)
	}
	message.ClientPublish(pbClient, "messenger", data, attrs)
	log.Println(monitor.Name, ": Pubsub message sent!")
}

func addMessagesEntries(messages []interface{}) {
	conn.InsertMany(
		conn.MSCollection(),
		messages,
	)
}

func handleEvent(eventType string, monitor check.Monitor, eventH event.Event, wg *sync.WaitGroup) {
	log.Println(monitor.Name, ": Handling", eventType, "event")
	var emailContacts []event.Contact
	var phoneContacts []event.Contact
	var slackIntegrations []event.Integration
	var level string
	sms, email, slack := false, false, false
	messages := make([]interface{}, len(eventH.Contacts)+len(eventH.Integrations))
	monitorName := monitor.Name
	monitorURL := monitor.URL
	monitorPath := monitor.Path
	monitorCLS := monitor.CLS
	monitorID := monitor.ID.Hex()

	if strings.Compare(eventType, DOWNTIME) == 0 {
		level = CRITICAL
	}

	if strings.Compare(eventType, UPTIME) == 0 {
		level = DONE
	}

	for index, contact := range eventH.Contacts {
		var cntct event.Contact
		var cntctParent event.Contact
		var mtype string

		err := conn.FindOne(
			conn.CCollection(),
			&bson.M{"_id": contact},
		).Decode(&cntct)
		if err != nil {
			log.Fatal(err)
		}
		err = conn.FindOne(
			conn.CCollection(),
			&bson.M{"_id": cntct.ContactParent},
		).Decode(&cntctParent)
		if err != nil {
			log.Fatal(err)
		}
		cntct.Name = cntctParent.Name
		cntct.LastName = cntctParent.LastName
		if cntct.CLS == EMAILCONTACT {
			emailContacts = append(emailContacts, cntct)
			if email == false {
				email = true
			}
			mtype = EMAIL
		}
		if cntct.CLS == PHONECONTACT {
			phoneContacts = append(phoneContacts, cntct)
			if sms == false {
				sms = true
			}
			mtype = SMS
		}
		messages[index] = bson.M{
			"sent":         true,
			"message_type": mtype,
			"level":        level,
			"created_at":   time.Now(),
			"user":         monitor.User,
			"monitor_name": monitorName,
			"monitor_url":  monitorURL,
			"monitor_path": monitorPath,
			"monitor_cls":  monitorCLS,
			"monitor_id":   monitorID,
			"contact_name": cntct.Name,
		}
	}

	for index, integration := range eventH.Integrations {
		var intgr event.Integration
		mtype := SLACK
		err := conn.FindOne(
			conn.ICollection(),
			&bson.M{"_id": integration},
		).Decode(&intgr)
		if err != nil {
			log.Fatal(err)
		}
		if intgr.CLS == SLACKINTEGRATION {
			slackIntegrations = append(slackIntegrations, intgr)
			if slack == false {
				slack = true
			}
		}
		messages[len(eventH.Contacts)+index] = bson.M{
			"sent":             true,
			"message_type":     mtype,
			"level":            level,
			"created_at":       time.Now(),
			"user":             monitor.User,
			"monitor_name":     monitorName,
			"monitor_url":      monitorURL,
			"monitor_path":     monitorPath,
			"monitor_cls":      monitorCLS,
			"monitor_id":       monitorID,
			"integration_name": intgr.TeamName + " - " + intgr.IncWebhookChannel,
		}
	}
	sendMessage(eventType, monitor, sms, email, slack, emailContacts, phoneContacts, slackIntegrations)
	if len(messages) > 0 {
		addMessagesEntries(messages)
	}
	log.Println(monitor.Name, ":", eventType, "event handled!")
	wg.Done()
}

func checkEvent(eventC event.Event, monitor check.Monitor, wg *sync.WaitGroup) {
	var wgRules sync.WaitGroup
	for _, rule := range eventC.Rules {
		var metric event.Metric
		err := conn.FindOne(
			conn.MTCollection(),
			&bson.M{"_id": rule.Metric},
		).Decode(&metric)
		if err != nil {
			log.Fatal(err)
		}
		if metric.Name == METRICRESPONSE {
			value, err := strconv.Atoi(rule.Value)
			if err != nil {
				log.Fatal(err)
			}
			switch rule.Operator {
			case GTE:
				if monitor.Response >= value {
					wgRules.Add(1)
					handleEvent(DOWNTIME, monitor, eventC, &wgRules)
				}
			case LT:
				if monitor.Response < value {
					wgRules.Add(1)
					handleEvent(UPTIME, monitor, eventC, &wgRules)
				}
			}
		}
	}
	wgRules.Wait()
	wg.Done()
}

func checkMonitor(mntrColl *mongo.Collection, monitor check.Monitor, wg *sync.WaitGroup) {
	log.Println(monitor.Name, ": Processing!")
	if monitor.Response != 0 {
		log.Println(monitor.Name, ": Checking events")
		evtColl := conn.ECollection()
		evtCur := conn.Find(evtColl, &bson.M{"monitor": monitor.ID})
		var wgEvents sync.WaitGroup
		for evtCur.Next(context.TODO()) {
			var event event.Event
			err := evtCur.Decode(&event)
			if err != nil {
				log.Fatal(err)
			}
			wgEvents.Add(1)
			checkEvent(event, monitor, &wgEvents)
		}
		wgEvents.Wait()
		evtCur.Close(context.TODO())
		log.Println(monitor.Name, ": Events checked!")
	} else {
		log.Println(monitor.Name, ": No check yet")
	}
	filter := bson.D{{"controlled", false}}
	update := bson.D{
		{"$set", bson.D{
			{"controlled", true},
		}},
	}
	conn.Update(mntrColl, &filter, &update)
	log.Println(monitor.Name, ": Processed!")
	wg.Done()
}

func main() {
	pbClient = message.Start()
	mntrColl := conn.MCollection()
	for true {
		mntrCursor := conn.Find(mntrColl, &bson.M{"controlled": false})

		var wg sync.WaitGroup

		for mntrCursor.Next(context.TODO()) {
			var monitor check.Monitor
			err := mntrCursor.Decode(&monitor)
			if err != nil {
				log.Fatal(err)
			}
			wg.Add(1)
			go checkMonitor(mntrColl, monitor, &wg)
		}
		wg.Wait()
		mntrCursor.Close(context.TODO())
		time.Sleep(time.Second)
	}
}
