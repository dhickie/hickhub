package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dhickie/hickhub/config"
	"github.com/dhickie/hickhub/messaging"
	"github.com/dhickie/hickhub/messaging/payloads"
	"github.com/gorilla/mux"
)

// Errors relating to finding devices
var (
	ErrDeviceNotFound               = errors.New("Unable to find device")
	ErrUnknownDeviceType            = errors.New("Unknown device type")
	ErrDeviceNotCapable             = errors.New("Device is unable to perform the requested command")
	ErrUnableToReadBody             = errors.New("An error occured trying to read the body of the request")
	ErrUnableToPublishMessage       = errors.New("An error occured trying to publish a message to the relevant device module")
	ErrUnableToReadResult           = errors.New("An error occured trying to unmarshal the response from the device's module")
	ErrDeviceUnableToPerformCommand = errors.New("An error occured when the device tried to perform the requested command")
	ErrFailedToMarshalNewState      = errors.New("An error occured trying to marshal the new state of the device")
)

// CommandController encapsulates all information related to controlling a device
type CommandController struct {
	devices []config.Device
}

// NewCommandController returns a new ControlController instance
func NewCommandController(config config.Config) *CommandController {
	return &CommandController{config.Devices}
}

// ControlDevice issues a command to control a particular device
func (c *CommandController) ControlDevice(w http.ResponseWriter, r *http.Request) {
	// Get the ID of the device this is aimed for
	vars := mux.Vars(r)
	deviceID := vars["id"]
	state := vars["state"]
	command := vars["cmd"]

	// Remove the trailing slash if there is one
	command = strings.TrimRight(command, "/")

	// Work out of this command is actually supported by the device in question
	device, err := c.findDevice(deviceID)
	if err != nil {
		http.Error(w, ErrDeviceNotFound.Error(), 400)
		return
	}

	canServe := c.deviceIsCapable(device, state, command)
	if !canServe {
		http.Error(w, ErrDeviceNotCapable.Error(), 400)
		return
	}

	// Get the body of the request (this is the detail of the command)
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err == io.EOF {
		body = []byte("")
	} else if err != nil {
		http.Error(w, ErrUnableToReadBody.Error(), 500)
		return
	}

	// Publish the message and request a reply containing the devices new state
	reply, err := c.publishMessage(device, state, command, string(body))
	if err != nil {
		http.Error(w, ErrUnableToPublishMessage.Error(), 500)
		return
	}

	// Unmarshal the payload as a CommandResultMessage
	result := new(payloads.CommandResultPayload)
	if err = json.Unmarshal([]byte(reply.Payload), result); err != nil {
		http.Error(w, ErrUnableToReadResult.Error(), 500)
		return
	}

	// Check that the result was a success
	if !result.Success {
		http.Error(w, fmt.Sprintf("An error occured when the device tried to perform the command: %v", result.Error), 500)
		return
	}

	// Respond with the device's new state
	stateJSON, err := json.Marshal(result.NewState.State)
	if err != nil {
		http.Error(w, ErrFailedToMarshalNewState.Error(), 500)
	}
	w.Write(stateJSON)
}

// publishMessage publishes the command message to the appropriate topic
func (c *CommandController) publishMessage(device *config.Device, state, command, detail string) (messaging.Message, error) {
	// Publish to the topic based on the type of device
	var topic string
	switch device.Type {
	case config.TypeTv:
		topic = messaging.TopicTv
	default:
		return messaging.Message{}, ErrUnknownDeviceType
	}

	msg, err := messaging.NewCommandMessage(device.ID, state, command, detail)
	if err != nil {
		return messaging.Message{}, err
	}

	return messaging.Request(topic, msg, 1000)
}

// deviceIsCapable determines if the device with the given ID can perform the given command
func (c *CommandController) deviceIsCapable(device *config.Device, state, command string) bool {
	if val, ok := device.Capabilities[state]; ok {
		// The device has this state, see if we can perform the command
		for _, v := range val {
			if v == command {
				return true
			}
		}
	}

	return false
}

// findDevice finds the device with the given ID in the controller's device list
func (c *CommandController) findDevice(deviceID string) (*config.Device, error) {
	for _, v := range c.devices {
		if v.ID == deviceID {
			return &v, nil
		}
	}

	return nil, ErrDeviceNotFound
}
