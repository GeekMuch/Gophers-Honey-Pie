package logparser

import "errors"

// getOpenCanaryLogLevel retrieves the log severity level based on provided log
// type.
func getOpenCanaryLogLevel(logType int) (int, error) {
	// TODO: Probably a better way to check level.
	if isInTypeArray(logType, CriticalTypes) {
		return CRITICAL, nil
	} else if isInTypeArray(logType, ScanTypes) {
		return SCAN, nil
	} else if isInTypeArray(logType, InformationalTypes) {
		return INFORMATIONAL, nil
	}
	return -1, errors.New("log type is not assigned to a severity level")
}

// isInTypeArray checks if an element is present in an array.
func isInTypeArray(element int, array []int) bool {
	for _, x := range array {
		if x == element {
			return true
		}
	}
	return false
}
