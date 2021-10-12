package main

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/helper"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

/*
	Runs functions in order.
*/
func main() {
	log.InitLog(true)
	config.ReadConfigFile()
	helper.CheckForInternet()
	helper.UpdateSystem()
	helper.CheckForC2Server(config.Config.C2)
	config.Config.IpStr = helper.GetIP().String()
	if !config.CheckIfDeviceIDExits() {
		api.GetDeviceIDFromAPI()
		config.WriteConfToYAML()
	}
	go api.GetConfFromBackend()
	go api.Heartbeat()
	for {}


	// api.GetDeviceIDFromAPI()
	// config.AddDeviceIDtoYAML()
	// config.StartSetupSequence()
}
