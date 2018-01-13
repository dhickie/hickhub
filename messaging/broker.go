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
	return membroker.Subscribe(topic, func(m membroker.Message) {
		// Unmarshal the message
		msg := new(Message)
		err := json.Unmarshal(m.Data, msg)
		if err != nil {
			return
		}

		// Add the reply topic to the message
		msg.Reply = m.Reply

		// Call the callback
		callback(*msg)
	})
}

// Request publishes a message and then waits for the response until the specified timeout
func Request(topic string, msg Message, timeout int) (Message, error) {
	// Marshal the message
	mashalled, err := json.Marshal(msg)
	if err != nil {
		return Message{}, err
	}

	// Publish the request
	reply, err := membroker.Request(topic, mashalled, timeout)
	if err != nil {
		return Message{}, err
	}

	// Unmarshal the reply
	rply := new(Message)
	err = json.Unmarshal(reply.Data, rply)
	return *rply, err
}
