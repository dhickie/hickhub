package messaging

import (
	"encoding/json"
	"time"

	"github.com/dhickie/hickhub/messaging/payloads"
)

// Represents the different type of messages that can be sent between modules
const (
	MessageTypeCommand = "command"
	MessageTypeLog     = "log"
)

// Message is a message sent between different modules
type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"` // Payload is the message content encoded as JSON
}

// NewCommandMessage returns a new command message with the provided command details
func NewCommandMessage(deviceID, command, detail string) (Message, error) {
	payload := payloads.CommandPayload{
		DeviceID: deviceID,
		Command:  command,
		Detail:   detail,
	}

	return NewMessage(MessageTypeCommand, payload)
}

// NewLogMessage returns a new log message with the provided log details
func NewLogMessage(logType, message string, timeStamp time.Time) (Message, error) {
	payload := payloads.LogPayload{
		Type:      logType,
		Message:   message,
		Timestamp: timeStamp,
	}

	return NewMessage(MessageTypeLog, payload)
}

// NewMessage returns a new message of the specified type with the specified payload
func NewMessage(msgType string, payload interface{}) (Message, error) {
	msg := new(Message)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return *msg, err
	}

	return Message{
		Type:    msgType,
		Payload: string(payloadJSON),
	}, nil
}
