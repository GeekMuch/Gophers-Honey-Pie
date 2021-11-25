package opencanary

import (
	"encoding/json"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/filewatcher"
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

func Initialize() error {
	err := readFromCanaryConfig()
	if err != nil {
		return err
	}
	_ = stopCanary()
	err = startCanary()
	if err != nil {
		log.Logger.Error().Msgf("Error starting opencanary: %s", err)
		return err
	}

	logChannel := filewatcher.NewLogChannel()
	go startChannel(logChannel)
	go func() {
		err = filewatcher.StartNewFileWatcher("/var/tmp/opencanary.log", "/home/zinja/offset.txt", logChannel)
		if err != nil {
			log.Logger.Error().Msgf("Filewatcher error: %s", err)
		}
	}()
	return nil
}

func stopCanary() error {
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

func startCanary() error {
	log.Logger.Info().Msg("[X]\tStarting OpenCanary!")
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
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Logger.Error().Msgf("Error closing file: %s", err)
		}
	}(file)
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
