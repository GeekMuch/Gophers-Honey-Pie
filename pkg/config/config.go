package config

import (
	"io/ioutil"

	model "github.com/Mikkelhost/Gophers-Honey/pkg/model"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"

	"gopkg.in/yaml.v3"
)

var Config *model.PiConf
var ConfPath string = "/boot/config.yml"

func ReadConfigFile() {

	yfile, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in reading YAML  - ", err)
	}

	// settings := make(map[string]model.PiConf)
	conf := model.PiConf{}
	err2 := yaml.Unmarshal(yfile, &conf)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError in unmarshal YAML - ", err2)
	}

	log.Logger.Info().Msgf("[*] Settings: \n\t\tC2:\t\t%v \n\t\tIPStr:\t\t%v \n\t\tHostname:\t%v \n\t\tMAC:\t%v \n\t\tConfigured:\t%v \n\t\tPort:\t\t%v \n\t\tDeviceID:\t%v \n\t\tDeviceKey:\t%v",
		conf.C2,
		conf.IpStr,
		conf.Hostname,
		conf.Mac,
		conf.Configured,
		conf.Port,
		conf.DeviceID,
		conf.DeviceKey)

	log.Logger.Info().Msgf("[*] Services: \n\t\tSSH:\t%v \n\t\tFTP:\t%v \n\t\tTELNET:\t%v \n\t\tHTTP:\t%v \n\t\tHTTPS:\t%v \n\t\tSMB:\t%v \n",
		conf.Services.SSH,
		conf.Services.FTP,
		conf.Services.TELNET,
		conf.Services.HTTP,
		conf.Services.HTTPS,
		conf.Services.SMB)

	Config = &conf
	// log.Logger.Debug().Msgf("Config: %v", *Config)
}

func CheckIfDeviceIDExits() bool {
	if Config.DeviceID == 0 {
		log.Logger.Warn().Msg("[!] Device ID not set")
		return false
	} else {
		log.Logger.Info().Msg("[+] Device ID set")
		return true
	}
}

func WriteConfToYAML() {

	log.Logger.Info().Msgf("[*]\tAdding configuration to YAML")

	data, err := yaml.Marshal(&Config)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in YAML Marshal - ", err)
	}

	err2 := ioutil.WriteFile(ConfPath, []byte(data), 0755)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError writing to YAML - ", err2)
	}
	log.Logger.Info().Msgf("[!] Device ID is: %v", Config.DeviceID)
}
