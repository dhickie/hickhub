package payloads

import (
	"time"

	"github.com/dhickie/hickhub/models"
)

// CommandPayload is the payload sent with command messages
type CommandPayload struct {
	DeviceID string `json:"device_id"`
	State    string `json:"state"` // The state (ie. "volume", "power" etc.) of the device that is being changed
	Command  string `json:"command"`
	Detail   string `json:"detail"` // Detail provides extra detail about the command, again encoded as JSON
}

// CommandResultPayload is the payload sent with the result of command messages
type CommandResultPayload struct {
	Success  bool               `json:"success"`
	Error    string             `json:"error"`
	NewState models.DeviceState `json:"new_state"`
}

// LogPayload is the payload sent with log messages
type LogPayload struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}
