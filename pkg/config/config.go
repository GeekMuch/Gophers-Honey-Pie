package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	model "github.com/Mikkelhost/Gophers-Honey/pkg/model"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"

	"gopkg.in/yaml.v3"
)

// type DeviceAuth struct {
// 	DeviceId  uint32 `json:"device_id,omitempty"`
// 	DeviceKey string `json:"device_key,omitempty"`
// }

type SetService struct {
	DeviceID uint32 `json:"device_id"`
}

type Settings struct {
	Hostname  string `yaml:"Hostname"`
	Port      int    `yaml:"port"`
	DeviceID  uint32 `yaml:"DeviceID"`
	DeviceKey string `yaml:"DeviceKey"`
}

type Services struct {
	SSH    bool `yaml:"SSH" json:"SSH"`
	FTP    bool `yaml:"FTP" json:"FTP"`
	RDP    bool `yaml:"RDP" json:"RDP"`
	SMB    bool `yaml:"SMB" json:"SMB"`
	TELNET bool `yaml:"TELNET" json:"TELNET"`
}

type Configuration struct {
	Settings Settings `yaml:"Settings"`
	Services Services `yaml:"Services" json:"services"`
}

func getConfFromBackend(hostname string, deviceID uint32) {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"

	sendStruct := &model.DeviceAuth{
		DeviceId:  deviceID,
		DeviceKey: "XxPFUhQ8R7kKhpgubt7v"}

	postBody, _ := json.Marshal(sendStruct)

	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, err := http.NewRequest("GET", "http://"+hostname+":8000/api/devices/getDeviceConf", responseBody)
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

func AddDeviceIDtoYAML(hostname string) {
	dID := api.GetDeviceID(hostname)

	log.Logger.Info().Msgf("[*]\tAdding Device ID to YAML")

	Configuration := model.Device{
		DeviceID: dID}

	data, err := yaml.Marshal(&Configuration)

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

	settings := make(map[string]Settings)
	err2 := yaml.Unmarshal(yfile, &settings)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}

	services := make(map[string]Services)
	err3 := yaml.Unmarshal(yfile, &services)
	if err3 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err3)
	}

	log.Logger.Info().Msgf("[*]Settings: \n\t\tHostname:\t%v \n\t\tPort:\t%v \n\t\tDeviceID:\t%v \n\t\tDeviceKey:\t%v",
		settings["Settings"].Hostname,
		settings["Settings"].Port,
		settings["Settings"].DeviceID,
		settings["Settings"].DeviceKey)
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

	settings := make(map[string]Settings)
	err2 := yaml.Unmarshal(yfile, &settings)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	if settings["Settings"].DeviceID == 0 {
		AddDeviceIDtoYAML(hostname)
	} else {
		ReadConfigFile()
		getConfFromBackend(hostname, settings["Settings"].DeviceID)
	}
}
