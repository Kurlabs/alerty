package main

import (
	"strings"

	"github.com/Kurlabs/alerty/internal/check/checking"

	// "github.com/Kurlabs/alerty/shared"

	// Internal calls

	"github.com/Kurlabs/alerty/shared/env"
)

const (
	MONITORTYPE = "Monitor.WebsiteMonitor"
)

func main() {
	if strings.Compare(env.Config.Level, "debug") == 0 {
		checking.Run(MONITORTYPE, 1)
		return
	}
	// production url
	checking.Cronjob(MONITORTYPE)
}
