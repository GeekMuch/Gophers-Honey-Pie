package config

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

func CheckForC2Server(C2 string) {

	C2Host := C2
	C2Port := Config.Port
	if !(len(C2Host) > 0) {
		log.Logger.Error().Msgf("[X]\tSite unreachable, no C2 set \n")
		log.Logger.Fatal()
		os.Exit(1)
	}
	log.Logger.Info().Msg("[*]\tChecking if C2 is Alive ")

	timeout := 10 * time.Second

	host := fmt.Sprintf("%s:%d", C2Host, C2Port)
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tC2 is unreachable, [ERROR] -  \n", err)
		log.Logger.Fatal()
	}
	log.Logger.Info().Msgf("[!]\tC2 is alive on -> %s", conn.LocalAddr().String())
}

/*
	Get local ip of this RPI
*/
func GetIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", nil
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
	//conn, err := net.Dial("udp", "8.8.8.8:80")
	//if err != nil {
	//	return nil, err
	//}
	//localAddr := conn.LocalAddr().(*net.UDPAddr)
	//
	//return localAddr.IP, nil
}

func AuthenticationToken() string {
	// Create a Bearer string by appending string access token
	// TODO: Change token to environment variable
	var bearer = "Bearer " + Config.DeviceKey
	return bearer
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
		ipstr, err := GetIP()
		if err != nil {
			log.Logger.Warn().Msgf("[X]\tError getting ip address: %s", err)
		}
		log.Logger.Info().Msgf("[!]\tIP is -> %s", ipstr)
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				log.Logger.Warn().Msgf("[X]\tError getting ip address: %s", err)
			}
		}(conn)
	}
}
func UpdateSystem() {
	log.Logger.Warn().Msg("[+]\tFetching updates!")
	// fmt.Println("[+] Fetching updates!")
	cmd := exec.Command("bash", "-c", "sudo apt update && sudo apt upgrade -y && sudo apt autoremove -y &> /dev/null")
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Logger.Error().Msgf("[X]\tCommand running failed [ERROR] - \n", err)
	}
	log.Logger.Info().Msgf("[*]\t[DONE] Updating")
}
