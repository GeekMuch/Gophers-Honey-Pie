package opencanary

import (
	"errors"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/helpers"
)

// getOpenCanaryLogLevel retrieves the log severity level based on provided log
// type.
func getOpenCanaryLogLevel(logType int) (int, error) {
	// TODO: Probably a better way to check level.
	if helpers.IsInTypeArray(logType, CriticalTypes) {
		return CRITICAL, nil
	} else if helpers.IsInTypeArray(logType, ScanTypes) {
		return SCAN, nil
	} else if helpers.IsInTypeArray(logType, InformationalTypes) {
		return INFORMATIONAL, nil
	}
	return -1, errors.New("log type is not assigned to a severity level")
}
