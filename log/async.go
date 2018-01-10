package log

import (
	"fmt"
	"time"

	"github.com/dhickie/hickhub/messaging"
)

// Constants for different types of log
const (
	LogTypeInfo  = "INFO"
	LogTypeWarn  = "WARN"
	LogTypeError = "ERROR"
)

// Info publishes an Info logging message to be processed by the logging module
func Info(message string) {
	doLog(LogTypeInfo, message)
}

// Warn publishes a Warn logging message to be processed by the logging module
func Warn(message string) {
	doLog(LogTypeWarn, message)
}

// Error publishes an Error logging message to be processed by the logging module
func Error(message string) {
	doLog(LogTypeError, message)
}

func doLog(logType, message string) {
	timestamp := time.Now()

	// Write the log out to console
	fmt.Printf("%v - %v - %v\r\n", timestamp, logType, message)

	msg, err := messaging.NewLogMessage(logType, message, timestamp)
	if err != nil {
		// Just swallow, not like we can log it...
		return
	}

	messaging.Publish(messaging.TopicLogging, msg)
}
