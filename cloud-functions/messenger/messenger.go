package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"cloud.google.com/go/pubsub"
)

// SMSTEMPLATE allows to define the SMS notification message
const (
	SMSTEMPLATE   = "Hey {{name}}, Your {{object_type}} ({{object}}) is {{status}}."
	SLACKTEMPLATE = "Your {{object_type}} ({{object}}) is {{status}}."
)

// Attributes contains all the required data to send notifications
type Attributes struct {
	ObjectType string `json:"object_type"`
	Object     string `json:"object"`
	Email      bool   `json:"email,string"`
	Sms        bool   `json:"sms,string"`
	Slack      bool   `json:"slack,string"`
	Event      string `json:"event"`
}

// Contact represents a contact to send alert
type Contact struct {
	Name     string
	LastName string
	Email    string
	Area     int
	Number   int
}

// Integrations like Slack
type Integrations struct {
	WebhookURL string `json:"inc_webhook_url"`
}

// ChannelsList contains Contacts array field to handle json parsing
type ChannelsList struct {
	EmailContacts     []Contact
	PhoneContacts     []Contact
	SlackIntegrations []Integrations
}

// PubSubMessage contains Pub Sub message information
type PubSubMessage struct {
	Data       []byte     `json:"data"`
	Attributes Attributes `json:attributes`
}

// parseEmailTemplate builds the email template with the required data
func parseEmailTemplate(attrs Attributes, status, name string) string {
	tmpltData := struct {
		Name, ObjectType, Object, Status string
	}{
		Name:       name,
		ObjectType: attrs.ObjectType,
		Object:     attrs.Object,
		Status:     status,
	}
	var buffer bytes.Buffer
	tmpl, err := template.ParseFiles("email_alert.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(&buffer, tmpltData); err != nil {
		log.Fatal(err)
	}
	result := buffer.String()
	return result
}

// SendSlackNotifications sends Slack Message to a channel
func sendSlackNotifications(attrs Attributes, status string, slackInt []Integrations) {
	log.Println("Sending SMS notifications")
	for _, channel := range slackInt {
		replacer := strings.NewReplacer(
			"{{object_type}}", attrs.ObjectType,
			"{{object}}", attrs.Object,
			"{{status}}", status,
		)
		log.Println("Building SLACK message!")
		webhookURL := channel.WebhookURL
		attachmentTitle := "[ALERTY]: " + attrs.ObjectType + " - " + attrs.Object
		attachmentText := replacer.Replace(SLACKTEMPLATE)

		msgAttrs := map[string]string{
			"text":             "Incident",
			"attachment_title": attachmentTitle,
			"attachment_text":  attachmentText,
			"webhook_url":      webhookURL,
		}
		log.Println("Sending Slack message to", webhookURL)
		publish("send-slack-message", "", msgAttrs)
	}
}

// SendSMSNotifications sends SMS notification to user when an event occurrs
func sendSMSNotifications(attrs Attributes, status string, contacts []Contact) {
	log.Println("Sending SMS notifications")
	for _, contact := range contacts {
		replacer := strings.NewReplacer(
			"{{name}}", contact.Name,
			"{{object_type}}", attrs.ObjectType,
			"{{object}}", attrs.Object,
			"{{status}}", status,
		)
		log.Println("Building SMS message!")
		msg := replacer.Replace(SMSTEMPLATE)
		code := strconv.Itoa(contact.Area)
		number := strconv.Itoa(contact.Number)
		phone := strings.Join([]string{"+", code, number}, "")
		msgAttrs := map[string]string{
			"number": phone,
		}
		log.Println("Sending SMS to", phone)
		publish("send_sms", msg, msgAttrs)
	}
}

// SendEmailNotifications sends Email notification to user when an event occurrs
func sendEmailNotifications(attrs Attributes, status string, contacts []Contact) {
	log.Println("Sending Email notifications")
	for _, contact := range contacts {
		email := contact.Email
		body := parseEmailTemplate(attrs, status, contact.Name)
		msgAttrs := map[string]string{
			"subject": "[ALERTY]: " + attrs.ObjectType + " - " + attrs.Object,
		}
		msgAttrs["email"] = email
		log.Println("Sending Email to", email)
		publish("send_email", body, msgAttrs)
	}
}

// Publish publishes a new message to the given topic
func publish(topicName, message string, attrs map[string]string) {
	ctx := context.Background()
	projectID := os.Getenv("project_id")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	topic := client.Topic(topicName)
	msg := &pubsub.Message{
		Data:       []byte(message),
		Attributes: attrs,
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Message published.")
}

// Messenger handle the notificacion process
func Messenger(ctx context.Context, m PubSubMessage) error {
	var channels ChannelsList
	err := json.Unmarshal(m.Data, &channels)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Notification process started!")
	var status string
	switch m.Attributes.Event {
	case "uptime":
		status = "live"
	case "downtime":
		status = "down"
	}
	if m.Attributes.Sms {
		sendSMSNotifications(m.Attributes, status, channels.PhoneContacts)
	}
	if m.Attributes.Email {
		sendEmailNotifications(m.Attributes, status, channels.EmailContacts)
	}
	if m.Attributes.Slack {
		sendSlackNotifications(m.Attributes, status, channels.SlackIntegrations)
	}
	log.Println("Notification process finished!")
	return nil
}

// func main() {
// 	attrs := Attributes{
// 		Sms:        true,
// 		Email:      true,
// 		Event:      "downtime",
// 		Object:     "www.alerty.online",
// 		ObjectType: "domain",
// 	}
// 	var emailContacts []Contact
// 	emailContacts = append(
// 		emailContacts,
// 		Contact{Name: "Jhon", LastName: "Ramirez", Email: "jhoniscoding@gmail.com"},
// 	)
// 	var phoneContacts []Contact
// 	phoneContacts = append(
// 		phoneContacts,
// 		Contact{Name: "Jhon", LastName: "Ram", Area: 57, Number: "3195206895"},
// 	)
// 	contactsData := ChannelsList{
// 		EmailContacts: emailContacts,
// 		PhoneContacts: phoneContacts,
// 	}
// 	data, err := json.Marshal(contactsData)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	psmsg := PubSubMessage{Attributes: attrs, Data: []byte(data)}
// 	ctx := context.Background()
// 	Messenger(ctx, psmsg)
// }
