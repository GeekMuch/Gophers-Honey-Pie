package opencanary

import (
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

func Start() {
	log.Logger.Debug().Msgf("Global piconf: %v", config.Config)

	for {

	}
}