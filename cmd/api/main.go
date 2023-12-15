package main

import (
	"log"

	"github.com/ashiqsabith123/user-details-svc/pkg/config"
	"github.com/ashiqsabith123/user-details-svc/pkg/di"
	"github.com/ashiqsabith123/user-details-svc/pkg/helper"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(helper.Red("Error while loading config", err))
	}
	di.IntializeService(config)
}
