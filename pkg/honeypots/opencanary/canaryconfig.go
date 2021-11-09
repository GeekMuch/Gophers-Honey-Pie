package opencanary

import (
	"encoding/json"
	"io/ioutil"

	//"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"os"
	// model "github.com/Mikkelhost/Gophers-Honey/pkg/model"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)



var CanaryConfPath = "boot/opencanary.conf" // PATH: /etc/opencanaryd/opencanary.conf
var conf *canaryConf

func ReadFromToCanaryConfig() {
	file, err := os.Open(CanaryConfPath)
	defer file.Close()
	if err != nil {
		log.Logger.Warn().Msgf("Error opening file: %s", err)
	}
	
	jFile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError in reading JSON  - ", err)
	}
	log.Logger.Debug().Msgf("ConfFile: %s", jFile)
	err = json.Unmarshal(jFile, &conf)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
	}
	log.Logger.Debug().Msgf("JSON : %v", *conf)
}


func WriteToCanaryConfigFile() {

	log.Logger.Info().Msgf("[*]\tAdding configuration to JSON")
	conf.FtpEnabled = !conf.FtpEnabled
	conf.SshPort = 1337
	data, err := json.MarshalIndent(&conf, "", "    ")
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in JSON Marshal - ", err)
	}

	err2 := ioutil.WriteFile("boot/canaryconf2.conf", []byte(data), 0755)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError writing to JSON - ", err2)
	}
	log.Logger.Info().Msgf("[!]  FTP: %v SSH port : %v", conf.FtpEnabled, conf.SshPort)
}

