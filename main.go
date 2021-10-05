package main

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

/*
	Runs fucntions in order.
*/
func main() {
	log.InitLog(true)
	config.ReadConfigFile()
	if !config.CheckIfDeviceIDExits() {
		api.GetDeviceIDFromAPI()
		config.AddDeviceIDtoYAML()
	}

	// api.GetDeviceIDFromAPI()
	// config.AddDeviceIDtoYAML()
	// config.StartSetupSequence()
}
