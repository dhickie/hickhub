package log

import (
	"time"

	"github.com/dhickie/openhub/log/models"
	"github.com/dhickie/openhub/messaging"
)

// Info publishes an Info logging message to be processed by the logging module
func Info(message string) {
	log := models.Log{
		Type:      models.LogTypeInfo,
		Timestamp: time.Now(),
		Message:   message,
	}
	doPublish(log)
}

// Warn publishes a Warn logging message to be processed by the logging module
func Warn(message string) {
	log := models.Log{
		Type:      models.LogTypeWarn,
		Timestamp: time.Now(),
		Message:   message,
	}
	doPublish(log)
}

// Error publishes an Error logging message to be processed by the logging module
func Error(message string) {
	log := models.Log{
		Type:      models.LogTypeError,
		Timestamp: time.Now(),
		Message:   message,
	}
	doPublish(log)
}

func doPublish(log models.Log) {
	msg := messaging.NewMessage(messaging.MessageTypeLog, log)
	messaging.Publish(messaging.TopicLogging, msg)
}
