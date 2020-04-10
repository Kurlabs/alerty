package crawling

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/jasonlvhit/gocron"
	"github.com/Kurlabs/alerty/internal/crawl"
	models "github.com/Kurlabs/alerty/shared/mongo"
	message "github.com/Kurlabs/alerty/shared/pubsub"

	"github.com/Kurlabs/alerty/internal/check/publishing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RunRobot Looks monitors and makes groups of 5 domains
// Send those domains to check-url cloudfuction via pubsub
func RunRobot(cls string, robotFrequency int, pbClient *pubsub.Client, collection *mongo.Collection) {
	log.Println("robot frequency:", robotFrequency)
	cur := models.Find(collection, &bson.M{"_cls": cls, "robot": true,
		"validated": true, "robot_frequency": robotFrequency, "status": true})
	domNum := 5
	counter := 0
	var checkURL []message.CheckURL

	for cur.Next(context.TODO()) {
		var monitor crawl.Monitor
		err := cur.Decode(&monitor)
		if err != nil {
			log.Fatal(err)
		}
		checkURL = append(checkURL, message.CheckURL{ID: monitor.ID, URL: monitor.URL, ActualResponse: monitor.Response})

		counter++

		if counter >= domNum {
			publishing.SendMessage(checkURL, pbClient, "robot-run")
			counter = 0
			checkURL = []message.CheckURL{}
		}
	}

	if len(checkURL) > 0 {
		publishing.SendMessage(checkURL, pbClient, "robot-run")
	}

	log.Println("CheckURL:", checkURL)
}

// RobotCronjob Define all cronjobs and exec the run method
func RobotCronjob(cls string, pbClient *pubsub.Client, mbCollectionCursor *mongo.Collection) {
	gocron.Every(30).Minutes().Do(RunRobot, cls, 30, pbClient, mbCollectionCursor)
	gocron.Every(90).Minutes().Do(RunRobot, cls, 90, pbClient, mbCollectionCursor)
	gocron.Every(720).Minutes().Do(RunRobot, cls, 720, pbClient, mbCollectionCursor)
	gocron.Every(1440).Minutes().Do(RunRobot, cls, 1440, pbClient, mbCollectionCursor)
	gocron.Every(10080).Minutes().Do(RunRobot, cls, 10080, pbClient, mbCollectionCursor)
	<-gocron.Start()
}
