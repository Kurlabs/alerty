package publishing

import (
	"encoding/json"
	"log"

	"github.com/Kurlabs/alerty/shared/env"
	message "github.com/Kurlabs/alerty/shared/pubsub"

	"cloud.google.com/go/pubsub"
	// Internal calls
)

// SendMessage recieve a list of checkURL and send it to cloudfunction
func SendMessage(checkURL []message.CheckURL, pbClient *pubsub.Client, pbTopic string) {
	output, err := json.Marshal(checkURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("CheckURL: ", string(output))
	log.Println("BrainURL: ", env.Config.BrainURL)
	msgAttrs := map[string]string{
		"BrainURL": env.Config.BrainURL,
	}
	message.ClientPublish(pbClient, pbTopic, output, msgAttrs)
}
