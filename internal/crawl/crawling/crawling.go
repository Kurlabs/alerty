package crawling

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/Kurlabs/alerty/internal/check/publishing"
	"github.com/Kurlabs/alerty/internal/check/storage/inmongo"
	"github.com/Kurlabs/alerty/internal/crawl"
	models "github.com/Kurlabs/alerty/shared/mongo"
	message "github.com/Kurlabs/alerty/shared/pubsub"
	"github.com/jasonlvhit/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	collection = models.MCollection()
	pbClient   = message.Start()
)

func buildAndSend(cur *mongo.Cursor, pbClient *pubsub.Client) error {
	domNum := 5
	counter := 0
	var checkURL []message.CheckURL

	for cur.Next(context.TODO()) {
		var monitor crawl.Monitor
		err := cur.Decode(&monitor)
		if err != nil {
			return err
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
	return nil
}

// RunRobot Looks monitors and makes groups of 5 domains
// Send those domains to check-url cloudfuction via pubsub
func RunRobot(cls string, robotFrequency int) {
	// Instace Mongo Collection
	checkRepo := inmongo.NewMonitorsRepository(collection)
	cur, err := checkRepo.Find(&bson.M{"_cls": cls, "robot": true,
		"validated": true, "robot_frequency": robotFrequency, "status": true})
	if err != nil {
		log.Fatal(err)
	}
	err = buildAndSend(cur, pbClient)
	if err != nil {
		log.Fatal(err)
	}
}

// RobotCronjob Define all cronjobs and exec the run method
func RobotCronjob(cls string) {
	gocron.Every(30).Minute().Do(RunRobot, cls, 30)
	gocron.Every(90).Minutes().Do(RunRobot, cls, 90)
	gocron.Every(720).Minutes().Do(RunRobot, cls, 720)
	gocron.Every(1440).Minutes().Do(RunRobot, cls, 1440)
	gocron.Every(10080).Minutes().Do(RunRobot, cls, 10080)
	<-gocron.Start()
}
