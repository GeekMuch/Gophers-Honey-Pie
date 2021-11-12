package api

import (
	"bytes"
	"encoding/json"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"net/http"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
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
	ipAddr := config.Config.IpStr

	// Encode the ip_addr to postbody
	postBody, _ := json.Marshal(map[string]string{
		"ip_str": ipAddr,
	})
	return postBody
}

func RegisterDevice() {
	responseBody := bytes.NewBuffer(createPostBody())
	// Create a new request using http
	req, err := http.NewRequest("POST", getAddDeviceURL(), responseBody)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -  \n", err)

	}
	// add authorization header to the req
	req.Header.Add("Authorization", config.AuthenticationToken())

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
	conf := model.PiConfResponse{
		Status:    "",
		DeviceId:  deviceId.Id,
		NICVendor: config.Config.NICVendor,
		Hostname:  config.Config.Hostname,
		Services:  config.Config.Services,
	}
	err = config.UpdateConfig(conf)
	if err != nil {
		log.Logger.Fatal().Msgf("Error updating conf with device id: %s", err)
		return
	}
}
