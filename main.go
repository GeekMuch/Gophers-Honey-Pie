package main

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots"
	"sync"

	//"github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots/opencanary"

	// "github.com/GeekMuch/Gophers-Honey-Pie/pkg/helper"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

/*
	Runs functions in order.
*/
func main() {
	var wg sync.WaitGroup

	log.InitLog(true)
	config.Initialize()
	if !config.CheckIfDeviceIDExits() {
		api.RegisterDevice()
		log.Logger.Info().Msgf("[+]\tFirst time configuration [DONE]")
	}

	log.Logger.Info().Msgf("[+]\t Initializing honeypots")
	if err := honeypots.Initialize(); err != nil {
		log.Logger.Fatal().Msgf("Error initializing honeypots: %s", err)
		return
	}

	wg.Add(1)
	go func() {
		api.GetConfFromBackend()
		log.Logger.Info().Msgf("GetConfFromBackend goroutine done")
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		api.Heartbeat()
		log.Logger.Info().Msgf("Heartbeat goroutine done")
		wg.Done()
	}()
	//opencanary.ReadFromToCanaryConfig()
	//opencanary.WriteToCanaryConfigFile()
	//opencanary.Start()
	//Test
	wg.Wait()

	// api.GetDeviceIDFromAPI()
	// config.AddDeviceIDtoYAML()
	// config.StartSetupSequence()
}
