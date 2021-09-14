package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLog(debug bool) {
	consoleWrite := zerolog.ConsoleWriter{Out: os.Stderr}
	Logger = zerolog.New(consoleWrite).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	if debug {
		Logger = zerolog.New(consoleWrite).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	}
}
