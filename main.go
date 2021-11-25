package main

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots"

	//"github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots/opencanary"

	// "github.com/GeekMuch/Gophers-Honey-Pie/pkg/helper"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

/*
	Runs functions in order.
*/
func main() {
	log.InitLog(true)
	config.Initialize()
	if !config.CheckIfDeviceIDExits() {
		api.RegisterDevice()
		log.Logger.Info().Msgf("[+]\tFirst time configuration [DONE]")
	}
	if err := honeypots.Initialize(); err != nil {
		log.Logger.Fatal().Msgf("Error initializing honeypots: %s", err)
		return
	}
	for {
	go api.GetConfFromBackend()
	go api.Heartbeat()
	//opencanary.ReadFromToCanaryConfig()
	//opencanary.WriteToCanaryConfigFile()
	//opencanary.Start()

	}

	// api.GetDeviceIDFromAPI()
	// config.AddDeviceIDtoYAML()
	// config.StartSetupSequence()
}
