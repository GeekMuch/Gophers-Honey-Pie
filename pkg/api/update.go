package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/helper"
	model "github.com/Mikkelhost/Gophers-Honey/pkg/model"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

func getDeviceConfURL() string {
	C2Host := config.Config.C2
	url := "http://" + C2Host + ":8000/api/devices/getDeviceConf"
	log.Logger.Info().Msg(url)
	return url
}

func GetConfFromBackend() {
	// Create a Bearer string by appending string access token
	var bearer = helper.AuthenticationToken()

	sendStruct := &model.DeviceAuth{
		DeviceId:  config.Config.DeviceID,
		DeviceKey: "XxPFUhQ8R7kKhpgubt7v"}

	postBody, _ := json.Marshal(sendStruct)

	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, err := http.NewRequest("GET", "http://"+config.Config.C2+":8000/api/devices/getDeviceConf", responseBody)
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
	//log.Logger.Info().Msgf("[+]\t Updated configs -> %v", respStruct)
	defer resp.Body.Close()

	config.Config.Services = respStruct.Services

	log.Logger.Info().Msgf("[*]Updated Services in config file: \n\t\tSSH:\t%v \n\t\tFTP:\t%v \n\t\tRDP:\t%v \n\t\tSMB:\t%v \n\t\tTELNET:\t%v \n",
		config.Config.Services.SSH,
		config.Config.Services.FTP,
		config.Config.Services.RDP,
		config.Config.Services.SMB,
		config.Config.Services.TELNET)
}

func Heartbeat () {
	for {
		var bearer = helper.AuthenticationToken()

		sendStruct := &model.Heartbeat{
			DeviceID: config.Config.DeviceID,
			TimeStamp: time.Now()}

		postBody, _ := json.Marshal(sendStruct)

		responseBody := bytes.NewBuffer(postBody)

		req, err := http.NewRequest("POST", "http://"+config.Config.C2+":8000/api/devices/heartbeat", responseBody)
		if err != nil {
			log.Logger.Error().Msgf("[X]\tError in http request.\n[ERROR] -  \n", err)
		}

		req.Header.Add("Authorization", bearer)
		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
		}

		resp.Body.Close()
		time.Sleep(time.Second*30)
	}

}
