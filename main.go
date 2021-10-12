package main

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/helper"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

/*
	Runs fucntions in order.
*/
func main() {
	log.InitLog(true)
	config.ReadConfigFile()
	helper.CheckForC2Server(config.Config.C2)
	if !config.CheckIfDeviceIDExits() {
		api.GetDeviceIDFromAPI()
		config.AddDeviceIDtoYAML()
	}
	api.GetConfFromBackend()
	config.AddDeviceIDtoYAML()

	// api.GetDeviceIDFromAPI()
	// config.AddDeviceIDtoYAML()
	// config.StartSetupSequence()
}
