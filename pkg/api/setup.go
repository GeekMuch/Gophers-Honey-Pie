package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"time"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"gopkg.in/yaml.v3"
)

type SendStruct struct {
	IpStr      string `json:"ip_str,omitempty"`
	Configured bool   `json:"configured"`
}

// Create struct to recive JSON format
type responseStruct struct {
	DeviceID   uint32 `json:"device_id"`
	Configured bool   `json:"configured"`
}

type Settiongs struct {
	Hostname  string `yaml:"Hostname"`
	Port      int    `yaml:"port"`
	DeviceID  uint32 `yaml:"DeviceID"`
	DeviceKey string `yaml:"DeviceKey"`
}

type Services struct {
	SSH    bool `yaml:"SSH" json:"SSH"`
	FTP    bool `yaml:"FTP" json:"FTP"`
	RDP    bool `yaml:"RDP" json:"RDP"`
	SMB    bool `yaml:"SMB" json:"SMB"`
	TELNET bool `yaml:"TELNET" json:"TELNET"`
}

var ConfPath string = "boot/config.yaml"

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
func GetURLForC2Server(hostname string) string {
	c2_host := hostname
	url := "http://" + c2_host + ":8000/api/devices/addDevice"
	return url
}

/*
	Check if C2 is Alive
*/
func CheckForC2Server(hostname string) {

	c2_host := hostname

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
func CheckForInternet(hostname string) {

	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		log.Logger.Error().Msgf("[X]\tConnection is down!")
	} else {
		log.Logger.Info().Msgf("[+]\tConnection is up!")
		log.Logger.Info().Msgf("[*]\tIP is -> %s", Get_ip())
		defer conn.Close()
	}
}

func GetHostnameYAML() string {

	yfile, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	settings := make(map[string]Settiongs)
	err2 := yaml.Unmarshal(yfile, &settings)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	log.Logger.Info().Msgf("[+] Hostname -> %s", settings["Settings"].Hostname)
	return settings["Settings"].Hostname
}

func GetDeviceID(hostname string) uint32 {
	ipAddr := Get_ip().String()

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"

	// Encode the data
	postBody, _ := json.Marshal(map[string]string{
		"ip_str": ipAddr,
	})

	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, err := http.NewRequest("POST", GetURLForC2Server(hostname), responseBody)
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

	var respStruct responseStruct

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&respStruct); err != nil {
		log.Logger.Error().Msgf("[X]\tError in decode.\n[ERROR] -  \n", err)
	}
	log.Logger.Info().Msgf("[+]\tNew DeviceID Added-> %d", respStruct.DeviceID)
	defer resp.Body.Close()

	return respStruct.DeviceID
}
