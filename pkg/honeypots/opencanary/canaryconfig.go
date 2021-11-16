package opencanary

import (
	"encoding/json"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"io/ioutil"
	"os/exec"

	//"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"os"
	// model "github.com/Mikkelhost/Gophers-Honey/pkg/model"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

var CanaryConfPath = "/etc/opencanaryd/opencanary.conf" //"boot/opencanary.conf"
var conf *canaryConf

func Initialize() error{
	err := readFromCanaryConfig()
	if err != nil {
		return err
	}
	stopCanary()
	err = startCanary()
	if err != nil {
		log.Logger.Error().Msgf("Error starting opencanary: %s", err)
		return err
	}
	return nil
}

func stopSMB() error{
	log.Logger.Info().Msg("[X]\tStopping Samba!")
	// fmt.Println("[+] Fetching updates!")
	cmd := exec.Command("systemctl", "stop","smdb")
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Logger.Info().Msgf("[*]\t[DONE] Samba Stoped")
	return nil
}

func startSMB() error{
	log.Logger.Info().Msg("[!]\tStarting Samba!")
	// fmt.Println("[+] Fetching updates!")
	cmd := exec.Command("systemctl", "start", "smdb")
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Logger.Info().Msgf("[*]\t[DONE] Samba Started")
	return nil
}


func stopCanary() error{
	log.Logger.Info().Msg("[X]\tStopping OpenCanary!")
	// fmt.Println("[+] Fetching updates!")
	cmd := exec.Command("opencanaryd", "--stop")
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Logger.Info().Msgf("[*]\t[DONE] OpenCanary Stoped")
	return nil
}

func startCanary() error{
	log.Logger.Info().Msg("[!]\tStarting OpenCanary!")
	// fmt.Println("[+] Fetching updates!")
	cmd := exec.Command("opencanaryd", "--start")
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Logger.Info().Msgf("[*]\t[DONE] OpenCanary Started")
	return nil
}

// Sets the opencanary conf pointer
func readFromCanaryConfig() error {
	file, err := os.Open(CanaryConfPath)
	defer file.Close()
	if err != nil {
		log.Logger.Warn().Msgf("Error opening file: %s", err)
	}
	jFile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError in reading JSON  - ", err)
	}
	//log.Logger.Debug().Msgf("ConfFile: %s", jFile)
	err = json.Unmarshal(jFile, &conf)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
	}
	//log.Logger.Debug().Msgf("JSON : %v", *conf)
	//log.Logger.Warn().Msgf("Hello %v", responseModel)

	return err
}

func writeToCanaryConfigFile(responseModel model.PiConfResponse) error {
	log.Logger.Info().Msgf("[*]\tAdding configuration to JSON %v", responseModel)
	log.Logger.Info().Msgf("[*]\tcanary conf %v", *conf)
	confBak := conf
	conf.SshEnabled = responseModel.Services.SSH
	conf.FtpEnabled = responseModel.Services.FTP
	conf.TelnetEnabled = responseModel.Services.TELNET
	conf.HttpEnabled = responseModel.Services.HTTP
	conf.SmbEnabled = responseModel.Services.SMB
	data, err := json.MarshalIndent(&conf, "", "    ")
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError in JSON Marshal - ", err)
		conf = confBak
		return err
	}
	err = ioutil.WriteFile(CanaryConfPath, []byte(data), 0755)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError writing to JSON - ", err)
		conf = confBak
		return err
	}

	return nil
}

func UpdateCanary(conf model.PiConfResponse) error {
	if conf.Services.SMB == true {
		if err := startSMB(); err != nil {
			log.Logger.Warn().Msgf("Error starting SMB: %s", err)
		}
	}
	if conf.Services.SMB == false {
		if err := stopSMB(); err != nil {
			log.Logger.Warn().Msgf("Error stopping SMB: %s", err)
		}
	}

	if err := stopCanary(); err != nil {
		log.Logger.Warn().Msgf("Error stopping opencanary: %s", err)
	}
	if err := writeToCanaryConfigFile(conf); err != nil {
		log.Logger.Warn().Msgf("Error writing to canary conf: %s", err)
	}
	if err := startCanary(); err != nil {
		log.Logger.Warn().Msgf("Error starting opencanary: %s", err)
	}
	return nil
}