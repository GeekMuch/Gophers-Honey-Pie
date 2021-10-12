package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	model "github.com/Mikkelhost/Gophers-Honey/pkg/model"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

func getDeviceConfURL() string {
	c2_host := config.Config.c2
	url := "http://" + c2_host + ":8000/api/devices/getDeviceConf"
	log.Logger.Info().Msg(url)
	return url
}

func authToken() string {
	// Create a Bearer string by appending string access token
	// TODO: Change token to environment variable
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"
	return bearer
}

func getConfFromBackend() {
	// Create a Bearer string by appending string access token
	var bearer = authToken()

	sendStruct := &model.DeviceAuth{
		DeviceId:  config.Config.DeviceID,
		DeviceKey: "XxPFUhQ8R7kKhpgubt7v"}

	postBody, _ := json.Marshal(sendStruct)

	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, err := http.NewRequest("GET", "http://"+c2+":8000/api/devices/getDeviceConf", responseBody)
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
