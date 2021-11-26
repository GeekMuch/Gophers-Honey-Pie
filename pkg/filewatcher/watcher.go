package filewatcher

import (
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"github.com/hpcloud/tail"
)

// StartNewFileWatcher reads the provided logfile as it is updated. Passes
// read lines to the provided channel. File truncation and replacement is
// handled. Should be run as go routine.
func StartNewFileWatcher(logFilepath, offsetFilepath string, logChannel *LogChannel) error {
	enablePolling := enableFilePolling()

	tailFile, err := tail.TailFile(logFilepath, tail.Config{
		Follow: true,
		ReOpen: true, // Config.ReOpen = true is analogous to linux command "tail -F".
		Poll:   enablePolling,
	})

	if err != nil {
		log.Logger.Error().Msgf("Tailfile error: %s", err)
		return err
	}

	var index uint32 = 0
	var offset uint32

	// Read offset value from file if it exists. Else set offset to 0.
	if fileExists(offsetFilepath) {
		offset, err = getOffsetFromFile(offsetFilepath)
		if err != nil {
			log.Logger.Error().Msgf("Error getting offset from offset file: %s", err)
			return err
		}
	} else {
		offset = 0
		err = saveOffsetToFile(offsetFilepath, offset)
		if err != nil {
			log.Logger.Error().Msgf("Error saving offset to offset file: %s", err)
			return err
		}
	}

	for line := range tailFile.Lines {
		if index >= offset {
			log.Logger.Info().Msgf("New line in log: %s", line.Text)
			// Send log line to log channel.
			logChannel.Logs <- line.Text
			offset++
			err = saveOffsetToFile(offsetFilepath, offset)
			if err != nil {
				log.Logger.Error().Msgf("Error saving offset to offset file: %s", err)
				return err
			}
		}
		index++
	}

	// Wait for more changes
	err = tailFile.Wait()
	if err != nil {
		log.Logger.Error().Msgf("Tailfile wait error: %s", err)
		return err
	}

	return nil
}
