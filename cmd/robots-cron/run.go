package main

import (
	"strings"

	"github.com/Kurlabs/alerty/internal/crawl/crawling"
	"github.com/Kurlabs/alerty/shared/env"
	message "github.com/Kurlabs/alerty/shared/pubsub"

	// Internal calls
	models "github.com/Kurlabs/alerty/shared/mongo"
)

const (
	MONITORTYPE = "Monitor.WebsiteMonitor"
)

func main() {
	// Instance pubsub pool connection
	pbClient := message.Start()
	// Instace Mongo Collection
	mbCollectionCursor := models.MCollection()
	if strings.Compare(env.Config.Level, "debug") == 0 {
		crawling.RunRobot(MONITORTYPE, 1, pbClient, mbCollectionCursor)
		return
	}
	// production url
	crawling.RobotCronjob(MONITORTYPE, pbClient, mbCollectionCursor)
}
