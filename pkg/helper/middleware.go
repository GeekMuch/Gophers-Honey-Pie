package helper

import (
	"net"
	"os"
	"time"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

func CheckForC2Server(C2 string) {

	C2_host := C2
	if !(len(C2_host) > 0) {
		log.Logger.Error().Msgf("[X]\tSite unreachable, no C2 set \n")
		log.Logger.Fatal()
		os.Exit(1)
	}
	log.Logger.Info().Msgf("[*]\tChecking if C2 with C2 is Alive -> %s", C2_host)

	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", C2_host+":8000", timeout)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tSite unreachable, [ERROR] -  \n", err)
		log.Logger.Fatal()
	}
	log.Logger.Info().Msgf("[*]\tC2 is alive on -> %s", conn.LocalAddr().String())
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

func AuthenticationToken() string {
	// Create a Bearer string by appending string access token
	// TODO: Change token to environment variable
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"
	return bearer
}
