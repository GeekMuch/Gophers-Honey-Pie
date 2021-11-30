package honeypots

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots/opencanary"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
)

func Initialize() error {
	err := opencanary.Initialize()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError initializing opencanary: %s", err)
		return err
	}
	return nil
}

func UpdateServices(conf model.PiConfResponse) error {
	err := opencanary.UpdateCanary(conf)
	if err != nil {
		return err
	}
	return nil
}
