package sendsms

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kurkop/twilio-go"
)

// Attributes is the message data structure
type Attributes struct {
	Number string `json:"number"`
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data       []byte     `json:"data"`
	Attributes Attributes `json:attributes`
}

// SendSMS sends a new SMS message using twillio client
func SendSMS(ctx context.Context, m PubSubMessage) error {
	sid := os.Getenv("sid")
	token := os.Getenv("token")

	from := "+13852658542"
	to := m.Attributes.Number
	message := string(m.Data)

	log.Printf("to: %s!", to)
	log.Printf("msg: %s!", message)

	client := twilio.NewClient(sid, token, nil)

	// Send a message
	msg, err := client.Messages.SendMessage(from, to, message, nil)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	fmt.Printf("Response: %s\n", msg)

	return nil
}
