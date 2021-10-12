package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	model "github.com/Mikkelhost/Gophers-Honey/pkg/model"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"

	"gopkg.in/yaml.v3"
)

var Config *model.PiConf
var ConfPath string = "boot/config.yaml"

type Services struct {
	SSH    bool `yaml:"SSH" json:"SSH"`
	FTP    bool `yaml:"FTP" json:"FTP"`
	RDP    bool `yaml:"RDP" json:"RDP"`
	SMB    bool `yaml:"SMB" json:"SMB"`
	TELNET bool `yaml:"TELNET" json:"TELNET"`
}

func getConfFromBackend(C2 string, deviceID uint32) {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"

	sendStruct := &model.DeviceAuth{
		DeviceId:  deviceID,
		DeviceKey: "XxPFUhQ8R7kKhpgubt7v"}

	postBody, _ := json.Marshal(sendStruct)

	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, err := http.NewRequest("GET", "http://"+C2+":8000/api/devices/getDeviceConf", responseBody)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)

	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
	}

	var respStruct model.PiConfResponse

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&respStruct); err != nil {
		log.Logger.Error().Msgf("[X]\tError in decode.\n[ERROR] -  \n", err)
	}
	log.Logger.Info().Msgf("[+]\t Added list of serices -> %v", respStruct)
	defer resp.Body.Close()
}

func AddDeviceIDtoYAML() {
	dID := Config.DeviceID

	log.Logger.Info().Msgf("[*]\tAdding Device ID to YAML")

	Configuration := model.PiConf{
		DeviceID: dID}

	data, err := yaml.Marshal(&Configuration)

	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}
	err2 := ioutil.WriteFile(ConfPath, []byte(data), 0755)

	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	log.Logger.Info().Msgf("[+]\tFirst time configuration [DONE]")

}

func ReadConfigFile() {

	yfile, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	// settings := make(map[string]model.PiConf)
	settings := model.PiConf{}
	err2 := yaml.Unmarshal(yfile, &settings)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}

	log.Logger.Info().Msgf("[*]Settings: \n\t\tC2:\t%v \n\t\tPort:\t%v \n\t\tDeviceID:\t%v \n\t\tDeviceKey:\t%v",
		settings.C2,
		settings.Port,
		settings.DeviceID,
		settings.DeviceKey)
	log.Logger.Info().Msgf("[*]Services: \n\t\tSSH:\t%v \n\t\tFTP:\t%v \n\t\tRDP:\t%v \n\t\tSMB:\t%v \n\t\tTELNET:\t%v \n",
		settings.Services.SSH,
		settings.Services.FTP,
		settings.Services.RDP,
		settings.Services.SMB,
		settings.Services.TELNET)
}

func CheckIfDeviceIDExits() {
	C2 := Config.C2
	yfile, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	settings := model.PiConf{}
	err2 := yaml.Unmarshal(yfile, &settings)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	if settings.DeviceID == 0 {
		AddDeviceIDtoYAML()
	} else {
		ReadConfigFile()
		getConfFromBackend(C2, settings.DeviceID)
	}
}
