package config

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Mikkelhost/Gophers-Honey/pkg/model"

	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"

	"gopkg.in/yaml.v3"
)

var Config *model.PiConf
var ConfPath = "/boot/config.yml"

func Initialize() {
	readConfigFile()
	CheckForInternet()
	CheckForC2Server(Config.C2)
	ip, err := GetIP()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError getting ip: %s", err)
		return
	}
	if Config.IpStr != ip {
		Config.IpStr = ip
		//todo update database on backend

	}
}

func readConfigFile() {

	yFile, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in reading YAML  - ", err)
	}

	// settings := make(map[string]model.PiConf)
	conf := model.PiConf{}
	err2 := yaml.Unmarshal(yFile, &conf)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError in unmarshal YAML - ", err2)
	}

	log.Logger.Info().Msgf("[*]\tSettings: \n\t\tC2:\t\t%v \n\t\tIPStr:\t\t%v \n\t\tHostname:\t%v \n\t\tNIC Vendor:\t%v \n\t\tMAC:\t%v \n\t\tConfigured:\t%v \n\t\tPort:\t\t%v \n\t\tDeviceID:\t%v \n\t\tDeviceKey:\t%v",
		conf.C2,
		conf.IpStr,
		conf.Hostname,
		conf.NICVendor,
		conf.Mac,
		conf.Configured,
		conf.Port,
		conf.DeviceID,
		conf.DeviceKey)

	log.Logger.Info().Msgf("[*]\tUpdated Services in config file: \n\t\tFTP:\t%v \n\t\tSSH:\t%v \n\t\tTELNET:\t%v \n\t\tHTTP:\t%v \n\t\tHTTPS:\t%v \n\t\tSMB:\t%v \n",
		conf.Services.FTP,
		conf.Services.SSH,
		conf.Services.TELNET,
		conf.Services.HTTP,
		conf.Services.HTTPS,
		conf.Services.SMB)

	Config = &conf
}

func rebootPi() error{
	log.Logger.Info().Msg("[X]\tRebooting Gophers Pi in 5 seconds!")
	time.Sleep(5 * time.Second)
	cmd := exec.Command("reboot" )
	err := cmd.Run()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError rebooting after Hostname change, command: %s", err)
		return err
	}
	return nil
}
func interfaceDown(iface string) error{
	log.Logger.Info().Msgf("[*]\tPutting Network interface %s DOWN", iface)

	cmd := exec.Command("ifconfig", iface, "down" )
	err := cmd.Run()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError in putting down  %s", err)
		return err
	}
	return nil
}
func interfaceUp(iface string) error{
	log.Logger.Info().Msgf("[*]\tPutting Network interface %s UP", iface)
	cmd := exec.Command("ifconfig", "eth0", "up" )
	err := cmd.Run()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError in putting down, command  %s", err)
		return err
	}
	return nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getNICVendorList() error {
	log.Logger.Info().Msgf("[ ! ]\tDownloading vendor list, please wait.. \n")

	cmd := exec.Command("wget", "http://standards-oui.ieee.org/oui/oui.csv", "-O", "pkg/config/NICVendors/vendors.csv")
	err := cmd.Run()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError in getNICVendor list, command  %s", err)
		return err
	}
	return nil
}
func readNICVendorFile(NICVendor string) string {

	if err := getNICVendorList(); err != nil {
		log.Logger.Warn().Msgf("[X]\tError getting vendor list: %s", err)
	}

	log.Logger.Info().Msgf("[*]\tReading NIC vendor list")

	var tmpMAC string

	in, err := os.Open("pkg/config/NICVendors/vendors.csv")
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError opening Vendor CSV file  %s", err)
	}

	r := csv.NewReader(in)

	for {
		column, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Logger.Warn().Msgf("[X]\tError in reading Vendor CSV file %s",err)
		}
		if strings.Contains(column[2], NICVendor) {
			tmpMAC = column[1]
			break
		}
	}
	// Concat first 3 vendor bytes with last 3 bytes
	val, _ := randomHex(3)
	finalMac:= tmpMAC + val


	// Separate into formal MAC address
	for i := 2; i < len(finalMac); i += 3 {
		finalMac = finalMac[:i] + ":" + finalMac[i:]
	}
	log.Logger.Info().Msgf("[!]\tNEW MAC is ->  %s",finalMac)
	return finalMac
}

func ChangeNICVendor(macAddress string, iface string) error{
	log.Logger.Debug().Msg("[!]\tChanging NIC Vendor!")

	if err := interfaceDown(iface); err != nil {
		log.Logger.Warn().Msgf("[X]\tError putting down network interface: %s", err)
		return err
	}

	cmd := exec.Command("macchanger","--mac", macAddress, iface )
	err := cmd.Run()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError in changing the NIC Vendor command: %s", err)
		return err
	}

	if err := interfaceUp(iface); err != nil {
		log.Logger.Warn().Msgf("[X]\tError getting network interface up: %s", err)
		return err
	}

	return nil
}

func updateHostname(hostname string)error{
	log.Logger.Debug().Msgf("[+]\tExecuting update hostname: %s", hostname)
	hostnameString := []byte(hostname)
	if hostname == "" {
		return nil
	}
	var err = ioutil.WriteFile("/etc/hostname", hostnameString, 0644)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError writing to /etc/hostname - ", err)
	}
	input, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError reading /etc/hosts - ", err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "127.0.1.1") {
			lines[i] = "127.0.1.1\t" + hostname
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("/etc/hosts", []byte(output), 0644)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError writing /etc/hosts - ", err)
	}

	// BIG NO NO RISK baaaad code
	//cmd := exec.Command("echo", hostname,">","/etc/hostname" )
	//err := cmd.Run()
	//if err != nil{
	//	log.Logger.Warn().Msgf("[X]\tError in Hostname command change: %s", err)
	//	return err
	//}
	return nil

}

func CheckIfDeviceIDExits() bool {
	log.Logger.Info().Msgf("[!]\tChecking device id: %d", Config.DeviceID)
	if Config.DeviceID == 0 {
		log.Logger.Warn().Msg("[!]\tDevice ID not set")
		return false
	} else {
		log.Logger.Info().Msg("[+]\tDevice ID set")
		return true
	}
}

func WriteConfToYAML() {

	log.Logger.Info().Msgf("[*]\tAdding configuration to YAML")

	data, err := yaml.Marshal(&Config)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in YAML Marshal - ", err)
	}
	err2 := ioutil.WriteFile(ConfPath, data, 0755)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError writing to YAML - ", err2)
	}
	log.Logger.Info().Msgf("[!]\tDevice ID is: %v", Config.DeviceID)
}

func UpdateConfig(conf model.PiConfResponse) error{
	//todo revert to old conf if something fails.
	//Making backup config
	//config := Config
	var rebootFlag = false
	if Config.DeviceID != conf.DeviceId {
		Config.DeviceID = conf.DeviceId
	}

	if Config.Hostname != conf.Hostname && conf.Hostname != "" {
		log.Logger.Info().Msgf("[ ! ]\tChange in hostname initialized\n")
		Config.Hostname = conf.Hostname
		if err := updateHostname(conf.Hostname); err != nil {
			log.Logger.Warn().Msgf("[X]\tError Changing Hostname: %s", err)
		}
		rebootFlag = true

		//todo Set hostname in respective files with func
	}

	if Config.NICVendor != conf.NICVendor && conf.NICVendor != "" {
		log.Logger.Info().Msgf("[ ! ]\tChange in NICVendor initialized\n")
		Config.NICVendor = conf.NICVendor
		macAddress := readNICVendorFile(conf.NICVendor)
		if err := ChangeNICVendor(macAddress, "eth0"); err != nil {
			log.Logger.Warn().Msgf("[X]\tError Changing NIC Vendor: %s", err)
		}
		//Todo generate new mac address with a new func
	}

	ip, err := GetIP()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError getting IP: %s", err)
		return err
	}

	if Config.IpStr != ip{
		Config.IpStr = ip
		//todo Make and update ip in backend
	}

	WriteConfToYAML()
	if rebootFlag {
		err := rebootPi()
		if err != nil {
			return err
		}
	}
	return nil
}