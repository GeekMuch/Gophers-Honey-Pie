package main

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/daemon"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"sync"
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
		daemon.UpdateDevice()
		log.Logger.Info().Msgf("UpdateDevice goroutine done")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		daemon.Heartbeat()
		log.Logger.Info().Msgf("Heartbeat goroutine done")
		wg.Done()
	}()

	wg.Wait()
}
