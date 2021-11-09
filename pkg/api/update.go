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

func GetConfFromBackend() {
	for {
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
			log.Logger.Info().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)

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
		resp.Body.Close()

		if config.Config.Services != respStruct.Services {
			config.Config.Services = respStruct.Services
			

		}
		config.Config.IpStr = helper.GetIP().String()
		config.WriteConfToYAML()

		log.Logger.Info().Msgf("[*] Updated Services in config file: \n\t\tSSH:\t%v \n\t\tFTP:\t%v \n\t\tRDP:\t%v \n\t\tSMB:\t%v \n\t\tTELNET:\t%v \n",
			config.Config.Services.SSH,
			config.Config.Services.FTP,
			config.Config.Services.TELNET,
			config.Config.Services.HTTP,
			config.Config.Services.HTTPS,
			config.Config.Services.SMB)


		time.Sleep(time.Second * 10)
	}
}

func Heartbeat() {
	for {
		var bearer = helper.AuthenticationToken()

		sendStruct := &model.Heartbeat{
			DeviceID: config.Config.DeviceID}

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
		time.Sleep(time.Second * 30)
		log.Logger.Info().Msgf("[*]\tHeartbeat ->  DeviceID: %v \n", sendStruct.DeviceID)
	}

}
