package config

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
)

func StartSetupSequence() {
	hostname := api.GetHostnameYAML()
	api.CheckForInternet(hostname)
	api.CheckForC2Server(hostname)
	// UpdateSystem()
	CheckIfDeviceIDExits(hostname)

}
