package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/helper"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

/*
	Returns URL for add device
*/
func getAddDeviceURL() string {
	C2Host := config.Config.C2
	url := "http://" + C2Host + ":8000/api/devices/addDevice"
	log.Logger.Info().Msg(url)
	return url
}

func createPostBody() []byte {
	ipAddr := helper.GetIP().String()

	// Encode the ip_addr to postbody
	postBody, _ := json.Marshal(map[string]string{
		"ip_str": ipAddr,
	})
	return postBody
}

func GetDeviceIDFromAPI() {

	responseBody := bytes.NewBuffer(createPostBody())
	// Create a new request using http
	req, err := http.NewRequest("POST", getAddDeviceURL(), responseBody)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)

	}
	// add authorization header to the req
	req.Header.Add("Authorization", helper.AuthenticationToken())

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)
	}
	// responseConfig := config.Config
	decoder := json.NewDecoder(resp.Body)
	var deviceId struct {
		Id uint32 `json:"device_id"`
	}
	if err := decoder.Decode(&deviceId); err != nil {
		log.Logger.Error().Msgf("[X]\tError in decode.\n[ERROR] -  \n", err)
	}
	log.Logger.Info().Msgf("[+]\tNew DeviceID Added-> %v", deviceId.Id)
	defer resp.Body.Close()
	config.Config.DeviceID = deviceId.Id
	config.Config.IpStr = helper.GetIP().String()
}
