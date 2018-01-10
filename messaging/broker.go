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

// Subscribe subscribes the provided callback function to the provided topic
func Subscribe(topic string, callback func(Message)) int {
	return membroker.Subscribe(topic, func(b []byte) {
		// Unmarshal the message and call the callback
		msg := new(Message)
		err := json.Unmarshal(b, msg)
		if err != nil {
			return
		}

		callback(*msg)
	})
}
