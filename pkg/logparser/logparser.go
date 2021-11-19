package logparser

import (
	"encoding/json"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/config"
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
)

var whitelist []string

// ParseOpenCanaryLog takes logs formatted by OpenCanary and converts them
// into a standardized log format.
func ParseOpenCanaryLog(jsonLog string) (StandardLog, error) {
	var standardLog StandardLog
	var opencanaryLog OpencanaryLog

	err := json.Unmarshal([]byte(jsonLog), &opencanaryLog)
	if err != nil {
		log.Logger.Warn().Msgf("Error unmarshalling JSON: %s", err)
		return StandardLog{}, err
	}

	// Drop log if source IP is whitelisted.
	if isWhitelisted(opencanaryLog.SrcHost) {
		log.Logger.Info().Msgf("IP: %s is whitelisted. Dropping log.", opencanaryLog.SrcHost)
		return StandardLog{}, nil
	}

	// Extract logdata/message as json.
	logdataMarshalled, err := json.Marshal(opencanaryLog.Logdata)
	if err != nil {
		log.Logger.Warn().Msgf("Error unmarshalling log type: %s", err)
		return StandardLog{}, err
	}

	// Parse to standardized log format.
	standardLog.DeviceID = config.Config.DeviceID
	standardLog.LogID = 0 // LogID is set by the backend.
	standardLog.DstHost = opencanaryLog.DstHost
	standardLog.DstPort = opencanaryLog.DstPort
	standardLog.SrcHost = opencanaryLog.SrcHost
	standardLog.SrcPort = opencanaryLog.SrcPort
	standardLog.LogTimeStamp = opencanaryLog.LocalTime
	standardLog.Message = string(logdataMarshalled)
	// Get severity level of log type.
	standardLog.Level, err = getSeverityLevel(opencanaryLog.LogType)
	if err != nil {
		log.Logger.Warn().Msgf("Error getting severity level: %s", err)
		return StandardLog{}, err
	}
	standardLog.LogType = OpencanaryLogTypes[opencanaryLog.LogType]
	// Include raw opencanary log for redundancy.
	standardLog.RawLog = jsonLog

	return standardLog, nil
}
