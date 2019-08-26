package main

import (
	"flag"
	"github.com/alfcope/checkouttest/config"
	"github.com/alfcope/checkouttest/pkg/logging"
	"github.com/alfcope/checkouttest/server"
)

func main() {
	configPath := flag.String("config", "", "path to configuration")
	flag.Parse()

	if *configPath == "" {
		*configPath = "./config"
	}

	configuration, err := config.LoadConfiguration(*configPath, "configuration")
	if err != nil {
		logging.Logger.Error("Shutting down. Error loading configuration: ", err.Error())
		return
	}

	api, err := server.NewCheckoutApi(configuration)
	if err != nil {
		logging.Logger.Error("Shutting down. Error initialing api: ", err.Error())
		return
	}

	api.RunServer(configuration.Server.Port)
}
