package daemon

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"time"
)

func UpdateDevice() {
	for {
		respStruct := api.GetConfFromBackend()

		log.Logger.Info().Msg("Updating config")
		err := config.UpdateConfig(respStruct)
		if err != nil {
			log.Logger.Warn().Msgf("[X]\tError updating config", err)
		}
		if config.Config.Services != respStruct.Services {
			config.Config.Services = respStruct.Services
			//Todo enable correct OpenCanary setting with new func
			//Todo start and stop opencanary
			if err = honeypots.UpdateServices(respStruct); err != nil {
				log.Logger.Warn().Msgf("[X]\tError updating opencanary: %s", err)
			}
		}
		log.Logger.Info().Msgf("[*]\tIP address: %s", config.Config.IpStr)
		log.Logger.Info().Msgf("[*]\tMAC address: %s", config.Config.Mac)
		log.Logger.Info().Msgf("[*]\tUpdated Services in config file from backend: "+
			"\n\tHostname: \t%v "+
			"\n\tNICVendor:\t%v "+
			"\n\tDeviceID:\t%v "+
			"\n\tStatus:\t%v "+
			"\n\t\tSSH:\t%v \n\t\tFTP:\t%v \n\t\tTELNET:\t%v \n\t\tHTTP:\t%v	 \n\t\tSMB:\t%v \n",
			respStruct.Hostname,
			respStruct.NICVendor,
			respStruct.DeviceId,
			respStruct.Status,
			respStruct.Services.SSH,
			respStruct.Services.FTP,
			respStruct.Services.TELNET,
			respStruct.Services.HTTP,
			respStruct.Services.SMB)

		log.Logger.Info().Msg("End of UpdateDevice loop before sleep")
		time.Sleep(time.Second * 10)
		log.Logger.Info().Msg("End of UpdateDevice loop after sleep")
	}
}

func Heartbeat() {
	for {
		api.SendHeartbeat()
	}
}
