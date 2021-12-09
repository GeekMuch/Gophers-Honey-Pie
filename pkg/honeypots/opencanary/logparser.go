package opencanary

import (
	"encoding/json"
	"time"

	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/api"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/filewatcher"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
)

// startChannelListener is used to listen on the given log channel and
// parse the results to the OpenCanary parser.
func startChannelListener(logChannel *filewatcher.LogChannel) {
	for {
		select {
		case msg := <-logChannel.Logs:
			canaryLog, err := ParseOpenCanaryLog(msg)
			if err != nil {
				log.Logger.Error().Msgf("Error parsing log: %s", err)
				logChannel.Logs <- msg // TODO: possible stuck logs?
			}
			err = api.SendLog(canaryLog)
			if err != nil {
				for err != nil {
					err = api.SendLog(canaryLog)
				}
			}
		}
	}
}

// ParseOpenCanaryLog takes logs formatted by OpenCanary and converts them
// into a standardized log format.
func ParseOpenCanaryLog(jsonLog string) (model.Log, error) {
	var standardLog model.Log
	var opencanaryLog OpenCanaryLog

	err := json.Unmarshal([]byte(jsonLog), &opencanaryLog)
	if err != nil {
		log.Logger.Warn().Msgf("Error unmarshalling log JSON: %s", err)
		return model.Log{}, err
	}

	// Extract logdata/message as json.
	logdataMarshalled, err := json.Marshal(opencanaryLog.Logdata)
	if err != nil {
		log.Logger.Warn().Msgf("Error unmarshalling log type: %s", err)
		return model.Log{}, err
	}

	// Parse log local time to RFC3339 format.
	parsedLogTime, err := time.Parse("2006-01-02 15:04:05.9", opencanaryLog.UTCTime)
	if err != nil {
		log.Logger.Error().Msgf("Error parsing log local time: %s", err)
		return model.Log{}, err
	}
	log.Logger.Debug().Msgf(parsedLogTime.String())

	// Parse to standardized log format.
	standardLog.DeviceID = config.Config.DeviceID
	standardLog.LogID = 0 // LogID is set by the backend.
	standardLog.DstHost = opencanaryLog.DstHost
	if dstPort, _ := opencanaryLog.DstPort.Int64(); dstPort < 0 {
		standardLog.DstPort = 0
	} else {
		dstPort, _ = opencanaryLog.DstPort.Int64()
		standardLog.DstPort = uint16(dstPort)
	}
	standardLog.SrcHost = opencanaryLog.SrcHost
	if srcPort, _ := opencanaryLog.SrcPort.Int64(); srcPort < 0 {
		standardLog.SrcPort = 0
	} else {
		srcPort, _ = opencanaryLog.SrcPort.Int64()
		standardLog.SrcPort = uint16(srcPort)
	}
	standardLog.LogTimeStamp = parsedLogTime
	standardLog.Message = string(logdataMarshalled)
	// Get severity level of log type.
	standardLog.Level, err = getOpenCanaryLogLevel(opencanaryLog.LogType)
	if err != nil {
		log.Logger.Warn().Msgf("Error getting severity level: %s", err)
		return model.Log{}, err
	}
	standardLog.LogType = OpenCanaryLogTypes[opencanaryLog.LogType]
	// Include raw opencanary log for storage.
	standardLog.RawLog = jsonLog

	return standardLog, nil
}
