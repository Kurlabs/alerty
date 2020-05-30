package checking

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/Kurlabs/alerty/internal/check"
	"github.com/Kurlabs/alerty/internal/check/publishing"
	"github.com/Kurlabs/alerty/internal/check/storage/inmongo"
	models "github.com/Kurlabs/alerty/shared/mongo"
	message "github.com/Kurlabs/alerty/shared/pubsub"
	"github.com/jasonlvhit/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func buildAndSend(cur *mongo.Cursor, pbClient *pubsub.Client) error {
	domNum := 50
	counter := 0
	var checkURL []message.CheckURL

	for cur.Next(context.TODO()) {
		var monitor check.Monitor
		err := cur.Decode(&monitor)
		if err != nil {
			return err
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
	return nil
}

// Run Looks monitors and makes groups of 50 domains
// Send those domains to check-url cloudfuction via pubsub
func Run(cls string, frequency int) {
	// Instance pubsub pool connection
	pbClient := message.Start()
	// Instace Mongo Collection
	collection := models.MCollection()
	log.Println("frequency:", frequency)
	checkRepo := inmongo.NewMonitorsRepository(collection)
	cur, err := checkRepo.Find(&bson.M{"_cls": cls,
		"validated": true, "frequency": frequency, "status": true})
	if err != nil {
		log.Fatal(err)
	}
	err = buildAndSend(cur, pbClient)
	if err != nil {
		log.Fatal(err)
	}
}

// Cronjob Define all cronjobs and exec the run method
func Cronjob(cls string) {
	gocron.Every(1).Minute().Do(Run, cls, 1)
	gocron.Every(2).Minutes().Do(Run, cls, 2)
	gocron.Every(3).Minutes().Do(Run, cls, 3)
	gocron.Every(5).Minutes().Do(Run, cls, 5)
	gocron.Every(10).Minutes().Do(Run, cls, 10)
	gocron.Every(15).Minutes().Do(Run, cls, 15)
	gocron.Every(30).Minutes().Do(Run, cls, 30)
	<-gocron.Start()
}
