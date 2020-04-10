package sendslack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Attributes is the message data structure
type Attributes struct {
	Text            string `json:"text"`
	AttachmentTitle string `json:"attachment_title"`
	AttachmentText  string `json:"attachment_text"`
	WebhookURL      string `json:"webhook_url"`
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data       []byte     `json:"data"`
	Attributes Attributes `json:attributes`
}

// Attributes is the message data structure
type SlackAttachment struct {
	AttachmentTitle string `json:"title"`
	AttachmentText  string `json:"text"`
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type SlackMessage struct {
	Text        string            `json:"text"`
	WebhookURL  string            `json:"webhook_url"`
	Attachments []SlackAttachment `json:"attachments"`
}

// SendSlackMessage sends a new message to a slack channel
func SendSlackMessage(ctx context.Context, m PubSubMessage) error {
	text := m.Attributes.Text
	log.Printf("text: %s!", text)

	// Instance message (https://api.slack.com/incoming-webhooks)
	var msg = SlackMessage{
		Text:        m.Attributes.Text,
		WebhookURL:  m.Attributes.WebhookURL,
		Attachments: []SlackAttachment{{AttachmentTitle: m.Attributes.AttachmentTitle, AttachmentText: m.Attributes.AttachmentText}},
	}

	// Struct to Bytes
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Printf(string(b))

	req, err := http.NewRequest("POST", m.Attributes.WebhookURL, bytes.NewBuffer(b))

	// Set the header in the request.
	req.Header.Set("Content-Type", "application/json")

	// Execute the request.
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("Error connection to Brain API", err)
	}
	defer resp.Body.Close()

	return nil
}

// Test data
// text Notification
// attachment_title Domain: https://alerty.online
// attachment_text it is sent from CloudFunctions
// webhook_url <slack hook>
