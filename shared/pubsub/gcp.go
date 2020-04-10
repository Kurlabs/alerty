package pubsub

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckURL is sent to Cloud Function
type CheckURL struct {
	ID             primitive.ObjectID `json:"id"`
	URL            string             `json:"url"`
	ActualResponse int                `json:"actual_response"`
}

// Start a PubSub Client
func Start() *pubsub.Client {
	ctx := context.Background()
	projectID := os.Getenv("project_id")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// Publish publishes a new message to the given topic
func ClientPublish(client *pubsub.Client, topicName string, message []byte, attrs map[string]string) {
	ctx := context.Background()
	topic := client.Topic(topicName)
	msg := &pubsub.Message{
		Data:       message,
		Attributes: attrs,
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Message published.")
}
