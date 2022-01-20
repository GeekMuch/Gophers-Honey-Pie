package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"io"
	"net/http"
	"time"
)

/*
	Returns URL for add device
*/
func getAddDeviceURL() string {
	C2Host := config.Config.C2
	C2Protocol := config.Config.C2Protocol
	C2Port := config.Config.Port
	url := fmt.Sprintf("%s://%s:%d/api/devices/addDevice", C2Protocol, C2Host, C2Port)
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
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] - %s \n", err)

	}
	// add authorization header to the req
	req.Header.Add("Authorization", config.AuthenticationToken())

	// Send req using http Client
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] - %s \n", err)
	}
	// responseConfig := config.Config
	decoder := json.NewDecoder(resp.Body)
	var deviceId struct {
		Id uint32 `json:"device_id"`
	}
	if err = decoder.Decode(&deviceId); err != nil {
		log.Logger.Error().Msgf("[X]\tError in decode.\n[ERROR] -  \n", err)
	}
	log.Logger.Info().Msgf("[+]\tNew DeviceID Added-> %v", deviceId.Id)

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Logger.Error().Msgf("[X]\tError closing body: %s", err)
		}
	}(resp.Body)

	conf := model.PiConfResponse{
		Status:    "",
		DeviceId:  deviceId.Id,
		NICVendor: config.Config.NICVendor,
		Hostname:  config.Config.Hostname,
		Services:  config.Config.Services,
	}
	err = config.UpdateConfig(conf)
	if err != nil {
		log.Logger.Fatal().Msgf("[X]\tError updating conf with device id: %s", err)
		return
	}
}
