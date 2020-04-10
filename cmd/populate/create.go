package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	models "github.com/Kurlabs/alerty/shared/mongo"

	// Internal calls

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const monitorsCount = 1

var userID, _ = primitive.ObjectIDFromHex("5d541ef3ddf0ee3f514181a5")
var random *rand.Rand
var urlArray = [...]string{
	"alerty.online", "alerty.tech", "www.kurlabs.com", "gitlab.com", "falsexyz.com", "http://facebook.com", "https://google.com",
	"https://app.asana.com", "https://mail.google.com", "web.whatsapp.com", "tour.golang.org", "sentry.io", "cloudflare.com",
}
var ipArray = [...]string{
	"35.186.190.139", "35.196.25.76", "54.68.233.232", "52.37.153.208", "159.122.100.42", "158.176.86.249", "192.155.215.197",
	"50.97.60.166", "169.55.87.188", "168.1.68.251", "184.173.213.195", "169.61.108.35", "159.8.77.42", "169.38.84.49",
}
var pathArray = [...]string{
	"/", "/path", "/another/path", "/document/d/1pOfV2vCacnm6eIp719dwAVchq80CV5hi_keC73ezkus/", "/mail/u/1/#inbox",
	"/organizations/datagran-ot/issues/?project=1453739", "/kurlabs-alerty/alerty", "/0/696687137857060/1136012473018722",
}
var protocolArray = [...]string{"tcp", "udp"}
var boolArray = [...]bool{true, false}
var freqArray = [...]int{1, 2, 3, 5, 10, 15, 30}

func initRandom() {
	s := rand.NewSource(time.Now().Unix())
	random = rand.New(s)
}

func getRandomIndex(n int) int {
	return random.Intn(n)
}

func buildWebsiteMonitorData(name string) bson.M {
	url := urlArray[getRandomIndex(len(urlArray))]
	ssl := boolArray[getRandomIndex(2)]
	status := boolArray[getRandomIndex(2)]
	validated := true
	timeout := random.Intn(100) + 10
	frequency := freqArray[getRandomIndex(len(freqArray))]
	controlled := boolArray[getRandomIndex(2)]
	monitor := bson.M{
		"_cls":       "Monitor.WebsiteMonitor",
		"user":       userID,
		"created_at": time.Now(),
		"timeout":    timeout,
		"name":       name,
		"status":     status,
		"validated":  validated,
		"ssl":        ssl,
		"url":        url,
		"frequency":  frequency,
		"controlled": controlled,
	}
	return monitor
}

func buildSocketMonitorData(name string) bson.M {
	url := urlArray[getRandomIndex(len(urlArray))]
	ip := ipArray[getRandomIndex(len(ipArray))]
	path := pathArray[getRandomIndex(len(pathArray))]
	protocol := protocolArray[getRandomIndex(len(protocolArray))]
	ssl := boolArray[getRandomIndex(2)]
	status := boolArray[getRandomIndex(2)]
	validated := true
	timeout := random.Intn(100) + 10
	frequency := freqArray[getRandomIndex(len(freqArray))]
	controlled := boolArray[getRandomIndex(2)]
	port := random.Intn(10000) + 22
	monitor := bson.M{
		"_cls":       "Monitor.SocketMonitor",
		"user":       userID,
		"created_at": time.Now(),
		"timeout":    timeout,
		"name":       name,
		"status":     status,
		"validated":  validated,
		"ssl":        ssl,
		"url":        url,
		"ip":         ip,
		"port":       port,
		"path":       path,
		"protocol":   protocol,
		"frequency":  frequency,
		"controlled": controlled,
	}
	return monitor
}

func buildEvents(monitor interface{}) bson.M {
	var metricID, _ = primitive.ObjectIDFromHex("5d6f317dfad1a42c5ee98805")
	var contactOne, _ = primitive.ObjectIDFromHex("5d6ad8b8d82e32c6c4a9311b")
	var contactTwo, _ = primitive.ObjectIDFromHex("5d6ad90ad8d0d5f699a93117")
	event := bson.M{
		"rules": []interface{}{
			bson.M{
				"metric":   metricID,
				"operator": ">=",
				"value":    "500",
			},
			bson.M{
				"metric":   metricID,
				"operator": "<",
				"value":    "500",
			},
		},
		"contacts":   []primitive.ObjectID{contactOne, contactTwo},
		"created_at": time.Now(),
		// "monitor" : monitor.ID
	}
	return event
}

func main() {
	initRandom()
	DBName := "alerty"
	monitorsCollection := "monitors"
	mongoHost := "localhost"
	if mh := os.Getenv("MONGO_HOST"); mh != "" {
		mongoHost = mh
	}
	client := models.Connect(DBName, mongoHost, "27017")
	collection := models.GetCollection(client, DBName, monitorsCollection)
	monitors := make([]interface{}, monitorsCount)
	for i := 1; i <= monitorsCount; i++ {
		option := boolArray[getRandomIndex(2)]
		monitorName := fmt.Sprintf("Monitor #%d", i)
		var monitor interface{}
		if option {
			monitor = buildWebsiteMonitorData(monitorName)
		} else {
			monitor = buildSocketMonitorData(monitorName)
		}
		monitors[i-1] = monitor
		event := buildEvents(monitor)
		evtColl := models.GetCollection(client, DBName, "events")
		models.Insert(evtColl, event)
	}
	models.InsertMany(collection, monitors)
	models.Close(client)
}
