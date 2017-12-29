package models

import "time"

// Constants for different types of log
const (
	LogTypeInfo  = "info"
	LogTypeWarn  = "warn"
	LogTypeError = "error"
)

// Log represents a log item to be stored
type Log struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}
