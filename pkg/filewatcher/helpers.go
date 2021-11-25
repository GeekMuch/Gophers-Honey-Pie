package filewatcher

import (
	"fmt"
	"github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"io"
	"os"
	"runtime"
)

// getOffsetFromFile fetches the offset value stored in the given file.
func getOffsetFromFile(filepath string) (uint32, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		logger.Logger.Error().Msgf("Error closing offset file: %s", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Logger.Error().Msgf("Error closing offset file: %s", err)
		}
	}(file)

	var offset uint32

	for {
		_, err = fmt.Fscanf(file, "%d", &offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Logger.Error().Msgf("Scanning offset from file: %s", err)
			return offset, err
		}
	}

	return offset, nil
}

// saveOffsetToFile saves the given offset to the given file.
func saveOffsetToFile(filepath string, offset uint32) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY, 0644)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	if err != nil {
		logger.Logger.Error().Msgf("Error opening offset file: %s", err)
		return err
	}

	_, err = file.WriteAt([]byte(fmt.Sprint(offset)), 0)
	if err != nil {
		logger.Logger.Error().Msgf("Error writing to offset file: %s", err)
		return err
	}

	return nil
}

// resetOffsetFile is used to set the offset value in the given offset
// file to 0.
func resetOffsetFile(filepath string) error {
	err := saveOffsetToFile(filepath, 0)
	if err != nil {
		return err
	}
	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// enableFilePolling checks for OS type at runtime and enables/disables
// file polling based on OS. Helper function for StartNewFileWatcher.
// See https://github.com/hpcloud/tail/issues/54 for more info.
func enableFilePolling() bool {
	currentOS := runtime.GOOS
	switch currentOS {
	case "windows":
		logger.Logger.Debug().Msgf("Windows detected. Using polling")
		return true
	case "linux":
		logger.Logger.Debug().Msgf("Linux detected. Using inotify/fsnotify")
		return false
	}
	logger.Logger.Warn().Msgf("Error getting OS. Using default polling setting")
	return false
}
