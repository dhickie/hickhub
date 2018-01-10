package payloads

import "time"

// CommandPayload is the payload sent with command messages
type CommandPayload struct {
	DeviceID string `json:"device_id"`
	Command  string `json:"command"`
	Detail   string `json:"detail"` // Detail provides extra detail about the command, again encoded as JSON
}

// LogPayload is the payload sent with log messages
type LogPayload struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}
