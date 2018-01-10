package logging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dhickie/hickhub/messaging/payloads"

	"github.com/dhickie/hickhub/config"
	"github.com/dhickie/hickhub/log"
	"github.com/dhickie/hickhub/messaging"
)

var (
	logger      log.Logger
	currentFile string
)

// Launch launches the logging module which makes sure logs get written out to file
func Launch(appConfig config.Config) {
	// Create the logger based on the current date
	currentFile = getLogFileName()
	var err error
	logger, err = log.NewLogger(currentFile)
	if err != nil {
		panic(err)
	}

	// Subscribe the logger to the logging topic
	messaging.Subscribe(messaging.TopicLogging, subscriber)

	// Set up the routine to periodically flush the logs out to file
	go worker()
}

func subscriber(msg messaging.Message) {
	// Just log the log from the payload
	log := new(payloads.LogPayload)
	err := json.Unmarshal([]byte(msg.Payload), log)
	if err != nil {
		// Just swallow, not like we can log it...
	}

	logger.Log(*log)
}

func worker() {
	ticker := time.NewTicker(time.Minute)

	for {
		<-ticker.C
		// Make sure the filename is still appropriate
		correctFile := getLogFileName()
		if correctFile != currentFile {
			// We've moved on to a new file, move the logger on to it
			err := logger.NewFile(correctFile)
			if err != nil {
				continue
			}

			currentFile = correctFile
		}

		// Flush the logger to write out to file
		logger.Flush()
	}
}

func getLogFileName() string {
	return fmt.Sprintf("%v.log", time.Now().Format("20060102"))
}
