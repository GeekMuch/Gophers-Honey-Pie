package api

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

/*
	Returns URL for add device
*/
func getAddDeviceURL() string {
	c2_host := config.Config.HostName
	url := "http://" + c2_host + ":8000/api/devices/addDevice"
	log.Logger.Info().Msg(url)
	return url
}

/*
	Get local ip of this RPI
*/
func Get_ip() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Logger.Error().Msgf("[X]\tConnection is down! [ERROR] -  \n", err)
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func authenticationToken() string {
	// Create a Bearer string by appending string access token
	// TODO: Change token to environment variable
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"
	return bearer
}

func createPostBody() []byte {
	ipAddr := Get_ip().String()

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
	req.Header.Add("Authorization", authenticationToken())

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
	config.Config.IpStr = Get_ip().String()

}
