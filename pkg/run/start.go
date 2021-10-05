package config

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
)

func StartSetupSequence() {
	api.GetHostnameYAML()
	api.CheckForInternet()
	api.CheckForC2Server()
	// UpdateSystem()
	CheckIfDeviceIDExits()

}
