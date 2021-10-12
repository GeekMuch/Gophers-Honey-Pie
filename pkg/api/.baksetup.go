package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"time"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"gopkg.in/yaml.v3"
)

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

/*
	Returns URL for add device
*/
func GetURLForC2Server(c2 string) string {
	c2_host := c2
	url := "http://" + c2_host + ":8000/api/devices/addDevice"
	return url
}

/*
	Check if C2 is Alive
*/
func CheckForC2Server(c2 string) {

	c2_host := config.Config.c2
	log.Logger.Info().Msgf("[*]\tChecking if C2 with c2 is Alive -> %s", c2_host)

	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", c2_host+":8000", timeout)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tSite unreachable, [ERROR] -  \n", err)
		log.Logger.Fatal()
	}
	log.Logger.Info().Msgf("[*]\tC2 is alive on -> %s", conn.LocalAddr().String())
}

/*
	Runs update command to update RPI
*/
func UpdateSystem() {
	log.Logger.Warn().Msg("[*]\tFetching updates!")
	// fmt.Println("[+] Fetching updates!")
	cmd := exec.Command("bash", "-c", "sudo apt update && sudo apt upgrade -y && sudo apt autoremove -y &> /dev/null")
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout
	err := cmd.Run()
	log.Logger.Info().Msgf("[+]\t[DONE] Updating")
	if err != nil {
		log.Logger.Error().Msgf("[X]\tCommand running failed [ERROR] - \n", err)
	}
}

/*
	Checks if RPI has internet
*/
func CheckForInternet() {

	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		log.Logger.Error().Msgf("[X]\tConnection is down!")
	} else {
		log.Logger.Info().Msgf("[+]\tConnection is up!")
		log.Logger.Info().Msgf("[*]\tIP is -> %s", Get_ip())
		defer conn.Close()
	}
}

func GetConfigYAML() {

	yfile, err := ioutil.ReadFile(config.ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	err2 := yaml.Unmarshal(yfile, config.Config)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	log.Logger.Info().Msgf("[+] c2 -> %s", &config.Config.c2)
}

func GetDeviceID(c2 string) {
	ipAddr := Get_ip().String()

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"

	// Encode the data
	postBody, _ := json.Marshal(map[string]string{
		"ip_str": ipAddr,
	})

	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, err := http.NewRequest("POST", GetURLForC2Server(config.Config.c2), responseBody)
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

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&config.Config); err != nil {
		log.Logger.Error().Msgf("[X]\tError in decode.\n[ERROR] -  \n", err)
	}
	log.Logger.Info().Msgf("[+]\tNew DeviceID Added-> %d", config.Config.DeviceID)
	defer resp.Body.Close()

}
