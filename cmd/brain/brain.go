package main

import (
	"github.com/Kurlabs/alerty/cmd/brain/router"
)

func main() {
	e := router.New()
	e.Start("0.0.0.0:3000")
}
