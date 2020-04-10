package sendmail

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

// Attributes contains all the required data to send notifications
type Attributes struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data       []byte     `json:"data"`
	Attributes Attributes `json:attributes`
}

// SendEmail sends a new email with the given info
func SendEmail(ctx context.Context, m PubSubMessage) error {
	domain := "mg.alerty.online"

	// You can find the Private API Key in your Account Menu, under "Settings":
	// (https://app.mailgun.com/app/account/security)
	privateAPIKey := os.Getenv("mailgun_key")

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, privateAPIKey)

	sender := "alerty.message@alerty.online"
	subject := m.Attributes.Subject
	recipient := m.Attributes.Email
	body := string(m.Data)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, "", recipient)
	message.SetHtml(string(body))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Printf("Sending email notification to: %v", recipient)
	// Send the message	with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return nil
}
