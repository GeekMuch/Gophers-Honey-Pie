package api

import (
	"bytes"
	"encoding/json"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	opencanaryconfig "github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots/opencanary"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"net/http"
	"time"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

func GetConfFromBackend() {
	getConf:
	for {
		// Create a Bearer string by appending string access token
		log.Logger.Info().Msg("Start of GetConfFromBackend loop")
		var bearer = config.AuthenticationToken()

		sendStruct := &model.DeviceAuth{
			DeviceId:  config.Config.DeviceID}

		postBody, _ := json.Marshal(sendStruct)

		responseBody := bytes.NewBuffer(postBody)

		log.Logger.Info().Msg("Creating http request for getDeviceConf")
		// Create a new request using http
		req, err := http.NewRequest("GET", "http://"+config.Config.C2+":8000/api/devices/getDeviceConf", responseBody)
		if err != nil {
			log.Logger.Info().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
			time.Sleep(time.Second * 5)
			goto getConf
		}
		// add authorization header to the req
		req.Header.Add("Authorization", bearer)

		// Send req using http Client

		log.Logger.Info().Msg("Sending http request")
		client := &http.Client{
			Timeout: time.Second*3,
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
			time.Sleep(time.Second * 5)
			goto getConf
		}

		var respStruct model.PiConfResponse
		log.Logger.Info().Msg("Decoding response")
		decoder := json.NewDecoder(resp.Body)

		if err := decoder.Decode(&respStruct); err != nil {
			log.Logger.Error().Msgf("[X]\tError in decode.\n[ERROR] -  \n", err)
			time.Sleep(time.Second * 5)
			goto getConf
		}
		err = resp.Body.Close()
		if err != nil {
			log.Logger.Error().Msgf("[X]\tError in response body close.\n[ERROR] -  \n", err)

			return
		}
		log.Logger.Info().Msg("Updating config")
		err = config.UpdateConfig(respStruct)
		if err != nil {
			log.Logger.Warn().Msgf("[X]\tError updating config", err)
		}

		if config.Config.Services != respStruct.Services {
			config.Config.Services = respStruct.Services
			//Todo enable correct OpenCanary setting with new func
			//Todo start and stop opencanary
			if err := opencanaryconfig.UpdateCanary(respStruct); err != nil {
				log.Logger.Warn().Msgf("[X]\tError updating opencanary: %s", err)
			}
		}
		log.Logger.Info().Msgf("[*]\tIP address: %s", config.Config.IpStr)
		log.Logger.Info().Msgf("[*]\tUpdated Services in config file from backend: " +
			"\n\tHostname: \t%v " +
			"\n\tNICVendor:\t%v " +
			"\n\tDeviceID:\t%v " +
			"\n\tStatus:\t%v " +
			"\n\t\tSSH:\t%v \n\t\tFTP:\t%v \n\t\tTELNET:\t%v \n\t\tHTTP:\t%v \n\t\tSMB:\t%v \n",
			respStruct.Hostname,
			respStruct.NICVendor,
			respStruct.DeviceId,
			respStruct.Status,
			respStruct.Services.SSH,
			respStruct.Services.FTP,
			respStruct.Services.TELNET,
			respStruct.Services.HTTP,
			respStruct.Services.SMB)
		log.Logger.Info().Msg("End of GetConfFromBackend loop before sleep")
		time.Sleep(time.Second * 10)
		log.Logger.Info().Msg("End of GetConfFromBackend loop after sleep")
	}
}

func Heartbeat() {
	loop:
	for {
		var bearer = config.AuthenticationToken()

		sendStruct := &model.Heartbeat{
		DeviceID: config.Config.DeviceID}

		postBody, _ := json.Marshal(sendStruct)

		responseBody := bytes.NewBuffer(postBody)

		req, err := http.NewRequest("POST", "http://"+config.Config.C2+":8000/api/devices/heartbeat", responseBody)
		if err != nil {
			log.Logger.Error().Msgf("[X]\tError in http request.\n[ERROR] -  \n", err)
			time.Sleep(time.Second * 5)
			goto loop
		}

		req.Header.Add("Authorization", bearer)
		// Send req using http Client
		client := &http.Client{
			Timeout: time.Second*3,
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
			time.Sleep(time.Second * 5)
			goto loop
		}

		err = resp.Body.Close()
		if err != nil {
			log.Logger.Error().Msgf("[X]\tError in heartbeat response body.\n[ERROR] -  \n", err)

			return
		}
		time.Sleep(time.Second * 30)
		log.Logger.Info().Msgf("[*]\tHeartbeat ->  DeviceID: %v \n", sendStruct.DeviceID)
	}

}
