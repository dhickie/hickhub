package messaging

import (
	"encoding/json"

	"github.com/dhickie/go-membroker"
)

// Publish publishes the provided message to the specified topic
func Publish(topic string, msg Message) error {
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	membroker.Publish(topic, []byte(json))
	return nil
}
