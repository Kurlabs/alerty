package main

import (
	"strings"

	"github.com/Kurlabs/alerty/internal/crawl/crawling"
	"github.com/Kurlabs/alerty/shared/env"
)

const (
	MONITORTYPE = "Monitor.WebsiteMonitor"
)

func main() {
	if strings.Compare(env.Config.Level, "debug") == 0 {
		crawling.RunRobot(MONITORTYPE, 1)
		return
	}
	// production url
	crawling.RobotCronjob(MONITORTYPE)
}
