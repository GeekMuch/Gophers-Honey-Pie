package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
)

func GetConfFromBackend() model.PiConfResponse {
	// Create a Bearer string by appending string access token
	log.Logger.Info().Msg("Start of GetConfFromBackend loop")
	var bearer = config.AuthenticationToken()

	sendStruct := &model.DeviceAuth{
		DeviceId: config.Config.DeviceID,
	}

	postBody, _ := json.Marshal(sendStruct)

	responseBody := bytes.NewBuffer(postBody)

	log.Logger.Info().Msg("Creating http request for getDeviceConf")
	// Create a new request using http
	req, err := http.NewRequest("GET", config.Config.C2Protocol+"://"+config.Config.C2+":8000/api/devices/getDeviceConf", responseBody)
	if err != nil {
		log.Logger.Info().Msgf("[X]\tError on request.\n[ERROR] -  \n", err)
		time.Sleep(time.Second * 5)
	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client

	log.Logger.Info().Msg("Sending http request")
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
		time.Sleep(time.Second * 5)
		return GetConfFromBackend()
	}

	var respStruct model.PiConfResponse
	log.Logger.Info().Msg("Decoding response")
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&respStruct); err != nil {
		log.Logger.Error().Msgf("[X]\tError in decode.\n[ERROR] -  \n", err)
		time.Sleep(time.Second * 5)
	}

	err = resp.Body.Close()
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in response body close.\n[ERROR] -  \n", err)
	}

	return respStruct
}

func SendHeartbeat() error {
	var bearer = config.AuthenticationToken()

	sendStruct := &model.Heartbeat{
		DeviceID: config.Config.DeviceID,
		IpStr:    config.Config.IpStr,
	}

	postBody, _ := json.Marshal(sendStruct)

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", config.Config.C2Protocol+"://"+config.Config.C2+":8000/api/devices/heartbeat", responseBody)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in http request.\n[ERROR] -  \n", err)
		return err
	}

	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in heartbeat response body.\n[ERROR] -  \n", err)
		return err
	}

	time.Sleep(time.Second * 30)
	log.Logger.Info().Msgf("[*]\tHeartbeat ->  DeviceID: %v \n", sendStruct.DeviceID)
	return nil
}
