package config

import (
	"io/ioutil"

	// model "github.com/Mikkelhost/Gophers-Honey/pkg/model"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"

	"gopkg.in/yaml.v3"
)
type opencanaryConfiguration struct {
	C2         string  `json:"ftp.enabled"`
	IpStr      string  `json:"ip_str"`
	Hostname   string  `json:"hostname"`
	Services   Service `json:"services"`
}
type Service struct {
	SSH    bool `bson:"ssh" yaml:"ssh" json:"ssh"`
	FTP    bool `bson:"ftp" yaml:"ftp"s json:"ftp"`
	TELNET bool `bson:"telnet" yaml:"telnet" json:"telnet"`
	HTTP   bool `bson:"http" yaml:"http" json:"http"`
	HTTPS  bool `bson:"https" yaml:"https" json:"https"`
	SMB    bool `bson:"smb" yaml:"smb" json:"smb"`
}

// var Config *model.PiConf
// var CanaryConfPath string = "/etc/opencanaryd/opencanary.conf"
var CanaryConfPath string = "boot/opencanary.conf"


func ReadFromToCanaryConfig() {
	yfile, err := ioutil.ReadFile(CanaryConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in reading YAML  - ", err)
	}

	// settings := make(map[string]model.PiConf)
	conf := opencanaryConfiguration{}
	err2 := yaml.Unmarshal(yfile, &conf)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError in unmarshal YAML - ", err2)
	}
	log.Logger.Info().Msgf("\n test: %v", &conf)
		
}

// func WriteToCanaryConfig() {
	
// }