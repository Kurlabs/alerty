package checking

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/jasonlvhit/gocron"
	"github.com/Kurlabs/alerty/internal/check"
	"github.com/Kurlabs/alerty/internal/check/publishing"
	models "github.com/Kurlabs/alerty/shared/mongo"
	message "github.com/Kurlabs/alerty/shared/pubsub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Run Looks monitors and makes groups of 50 domains
// Send those domains to check-url cloudfuction via pubsub
func Run(cls string, frequency int, pbClient *pubsub.Client, collection *mongo.Collection) {
	log.Println("frequency:", frequency)
	cur := models.Find(collection, &bson.M{"_cls": cls,
		"validated": true, "frequency": frequency, "status": true})
	domNum := 50
	counter := 0
	var checkURL []message.CheckURL

	for cur.Next(context.TODO()) {
		var monitor check.Monitor
		err := cur.Decode(&monitor)
		if err != nil {
			log.Fatal(err)
		}
		checkURL = append(checkURL, message.CheckURL{ID: monitor.ID, URL: monitor.URL, ActualResponse: monitor.Response})

		counter++

		if counter >= domNum {
			publishing.SendMessage(checkURL, pbClient, "check-url")
			counter = 0
			checkURL = []message.CheckURL{}
		}
	}

	if len(checkURL) > 0 {
		publishing.SendMessage(checkURL, pbClient, "check-url")
	}

	log.Println("CheckURL:", checkURL)
}

// Cronjob Define all cronjobs and exec the run method
func Cronjob(cls string, pbClient *pubsub.Client, mbCollectionCursor *mongo.Collection) {
	gocron.Every(1).Minute().Do(Run, cls, 1, pbClient, mbCollectionCursor)
	gocron.Every(2).Minutes().Do(Run, cls, 2, pbClient, mbCollectionCursor)
	gocron.Every(3).Minutes().Do(Run, cls, 3, pbClient, mbCollectionCursor)
	gocron.Every(5).Minutes().Do(Run, cls, 5, pbClient, mbCollectionCursor)
	gocron.Every(10).Minutes().Do(Run, cls, 10, pbClient, mbCollectionCursor)
	gocron.Every(15).Minutes().Do(Run, cls, 15, pbClient, mbCollectionCursor)
	gocron.Every(30).Minutes().Do(Run, cls, 30, pbClient, mbCollectionCursor)
	<-gocron.Start()
}
