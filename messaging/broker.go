package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/dhickie/go-membroker"
	"github.com/dhickie/openhub/log/models"
)

// Publish publishes the provided message to the specified topic
func Publish(topic string, msg Message) error {
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	fmt.Println(string(json))

	membroker.Publish(topic, []byte(json))
	return nil
}

// Subscribe subscribes the provided callback function to the provided topic
func Subscribe(topic string, callback func(Message)) int {
	return membroker.Subscribe(topic, func(b []byte) {
		// Unmarshal the byte array as a raw json map to know what type we're dealing with
		var msg Message
		var raw map[string]*json.RawMessage
		err := json.Unmarshal(b, &msg)
		if err != nil {
			return
		}
		err = json.Unmarshal(b, &raw)
		if err != nil {
			return
		}

		// Work out what the payload is based on the type of message
		var payload interface{}
		if msg.Type == MessageTypeSetVolume || msg.Type == MessageTypeSetChannel {
			payload = new(int)
		} else if msg.Type == MessageTypeLog {
			payload = new(models.Log)
		}

		// If there's a payload, unmarshal it
		if payload != nil {
			err = json.Unmarshal(*raw["payload"], payload)
			if err != nil {
				return
			}
		}

		// Add the payload to the message and invoke the callback
		switch payload.(type) {
		case *int:
			msg.Payload = *payload.(*int)
		case *models.Log:
			msg.Payload = *payload.(*models.Log)
		}
		callback(msg)
	})
}
