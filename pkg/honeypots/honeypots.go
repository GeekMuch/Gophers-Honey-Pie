package honeypots

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/honeypots/opencanary"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

func Initialize() error {
	err := opencanary.Initialize()
	if err != nil {
		log.Logger.Warn().Msgf("[X]\tError initializing opencanary: %s", err)
		return err
	}
	return nil
}