package filewatcher

import (
	log "github.com/GeekMuch/Gophers-Honey-Pie/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/hpcloud/tail"
	"runtime"
)

// TODO need to test and choose one of the two filewatching methods

var (
	opencanaryLogfile = "/var/tmp/opencanary/opencanary.log"
	kernelLogfile     = "/var/log/kern.log"
)

// enableFilePolling checks for OS type at runtime and enables/disables
// file polling based on OS. Helper function for StartNewTailWatcher.
// See https://github.com/hpcloud/tail/issues/54 for more info.
func enableFilePolling() bool {
	os := runtime.GOOS
	switch os {
	case "windows":
		log.Logger.Debug().Msgf("Windows detected. Using polling.")
		return true
	case "linux":
		log.Logger.Debug().Msgf("Linux detected. Using inotify/fsnotify.")
		return false
	}
	log.Logger.Warn().Msgf("Error getting OS. Using default polling setting.")
	return false
}

// StartNewTailWatcher reads logfiles as they are updated. Passes each
// line read to the parser. File truncation and replacement is handled.
// Must be run as go routine.
// TODO might need to take logfile path as input instead.
func StartNewTailWatcher() error {
	enablePolling := enableFilePolling()

	tailFile, err := tail.TailFile(opencanaryLogfile, tail.Config{
		Follow: true,
		ReOpen: true, // Config.ReOpen = true is analogous to linux command "tail -F".
		Poll:   enablePolling,
	})

	if err != nil {
		log.Logger.Error().Msgf("Tailfile error: %s", err)
		return err
	}

	for line := range tailFile.Lines {
		log.Logger.Info().Msgf("New line in log: %s", line)
		// TODO send line to parser.
	}

	err = tailFile.Wait()
	if err != nil {
		log.Logger.Error().Msgf("Tailfile wait error: %s", err)
		return err
	}

	return nil
}

// StartNewFSWatcher generates events on file operations on logfiles.
// TODO might need to take logfile path as input instead.
func StartNewFSWatcher() error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Logger.Warn().Msgf("Error creating new watcher: %s", err)
		return err
	}

	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Logger.Warn().Msgf("Error closing watcher: %s", err)
		}
	}(watcher)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Logger.Info().Msgf("Watcher event: %#v", event)
				// TODO get message and send to parser.
			case err := <-watcher.Errors:
				log.Logger.Warn().Msgf("Watcher error: %s", err)
			}
		}
	}()

	err = watcher.Add(opencanaryLogfile)
	if err != nil {
		log.Logger.Warn().Msgf("Error adding file to watcher: %s", err)
		return err
	}

	<-done

	return nil
}
