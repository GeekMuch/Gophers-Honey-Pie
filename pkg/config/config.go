package config

import (
	"io/ioutil"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"

	"gopkg.in/yaml.v3"
)

type SetService struct {
	DeviceID uint32 `json:"device_id"`
}
type Record struct {
	Hostname string `yaml:"Hostname"`
	DeviceID uint32 `yaml:"DeviceID"`
}

type Services struct {
	SSH    bool `yaml:"SSH"`
	FTP    bool `yaml:"FTP"`
	RDP    bool `yaml:"RDP"`
	SMB    bool `yaml:"SMB"`
	TELNET bool `yaml:"TELNET"`
}

type Config struct {
	Record   Record   `yaml:"Settings"`
	Services Services `yaml:"Services"`
}

func AddDeviceIDtoYAML(hostname string) {
	dID := api.GetDeviceID(hostname)

	log.Logger.Info().Msgf("[*]\tAdding Device ID to YAML")

	config := Config{
		Record: Record{Hostname: hostname, DeviceID: dID}}

	data, err := yaml.Marshal(&config)

	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}
	err2 := ioutil.WriteFile(api.ConfPath, []byte(data), 0755)

	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	log.Logger.Info().Msgf("[+]\tFirst time configuration [DONE]")

}

func ReadConfigFile() {

	yfile, err := ioutil.ReadFile(api.ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	settings := make(map[string]Record)
	err2 := yaml.Unmarshal(yfile, &settings)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}

	services := make(map[string]Services)
	err3 := yaml.Unmarshal(yfile, &services)
	if err3 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err3)
	}

	log.Logger.Info().Msgf("[*]Settings: \n\t\tHostname:\t%v \n\t\tDeviceID:\t%v",
		settings["Settings"].Hostname,
		settings["Settings"].DeviceID)
	log.Logger.Info().Msgf("[*]Services: \n\t\tSSH:\t%v \n\t\tFTP:\t%v \n\t\tRDP:\t%v \n\t\tSMB:\t%v \n\t\tTELNET:\t%v \n",
		services["Services"].SSH,
		services["Services"].FTP,
		services["Services"].RDP,
		services["Services"].SMB,
		services["Services"].TELNET)
}

func CheckIfDeviceIDExits(hostname string) {
	yfile, err := ioutil.ReadFile(api.ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	settings := make(map[string]Record)
	err2 := yaml.Unmarshal(yfile, &settings)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	if settings["Settings"].DeviceID == 0 {
		AddDeviceIDtoYAML(hostname)
	} else {
		ReadConfigFile()
	}
}
