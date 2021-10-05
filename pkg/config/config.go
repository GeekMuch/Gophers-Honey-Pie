package config

import (
	"io/ioutil"

	model "github.com/Mikkelhost/Gophers-Honey/pkg/model"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"

	"gopkg.in/yaml.v3"
)

var Config *model.PiConf
var ConfPath string = "boot/config.yaml"

func ReadConfigFile() {

	yfile, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	// settings := make(map[string]model.PiConf)
	conf := model.PiConf{}
	err2 := yaml.Unmarshal(yfile, &conf)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}

	log.Logger.Info().Msgf("[*]Settings: \n\t\tHostname:\t%v \n\t\tPort:\t%v \n\t\tDeviceID:\t%v \n\t\tDeviceKey:\t%v \n\t\tIPStr:\t%v \n\t\tConfigured:\t%v",
		conf.HostName,
		conf.Port,
		conf.DeviceID,
		conf.DeviceKey,
		conf.IpStr,
		conf.Configured)

	log.Logger.Info().Msgf("[*]Services: \n\t\tSSH:\t%v \n\t\tFTP:\t%v \n\t\tRDP:\t%v \n\t\tSMB:\t%v \n\t\tTELNET:\t%v \n",
		conf.Services.SSH,
		conf.Services.FTP,
		conf.Services.RDP,
		conf.Services.SMB,
		conf.Services.TELNET)

	Config = &conf
	// log.Logger.Debug().Msgf("Config: %v", *Config)
}

func CheckIfDeviceIDExits() bool {
	if Config.DeviceID == 0 {
		log.Logger.Warn().Msg("Device ID not set")
		return false
	} else {
		log.Logger.Info().Msg("Device ID set")
		return true
	}
}

func AddDeviceIDtoYAML() {

	log.Logger.Info().Msgf("[*]\tAdding Device ID to YAML")

	data, err := yaml.Marshal(&Config)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	err2 := ioutil.WriteFile(ConfPath, []byte(data), 0755)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	log.Logger.Info().Msgf("[+]\tFirst time configuration [DONE]")
	log.Logger.Info().Msgf("Device ID is: %v", Config.DeviceID)
}
