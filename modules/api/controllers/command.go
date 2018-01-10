package controllers

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dhickie/hickhub/config"
	"github.com/dhickie/hickhub/messaging"
	"github.com/gorilla/mux"
)

// Errors relating to finding devices
var (
	ErrDeviceNotFound         = errors.New("Unable to find device")
	ErrUnknownDeviceType      = errors.New("Unknown device type")
	ErrDeviceNotCapable       = errors.New("Device is unable to perform the requested command")
	ErrUnableToReadBody       = errors.New("An error occured trying to read the body of the request")
	ErrUnableToPublishMessage = errors.New("An error occured trying to publish a message to the relevant device module")
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
	command := vars["cmd"]

	// Remove the trailing slash if there is one
	command = strings.TrimRight(command, "/")

	// Work out of this command is actually supported by the device in question
	device, err := c.findDevice(deviceID)
	if err != nil {
		http.Error(w, ErrDeviceNotFound.Error(), 400)
		return
	}

	canServe, err := c.deviceIsCapable(device, command)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	} else if !canServe {
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

	// Publish the message
	err = c.publishMessage(device, command, string(body))
	if err != nil {
		http.Error(w, ErrUnableToPublishMessage.Error(), 500)
	}
}

// publishMessage publishes the command message to the appropriate topic
func (c *CommandController) publishMessage(device *config.Device, command, detail string) error {
	// Publish to the topic based on the type of device
	var topic string
	switch device.Type {
	case config.TypeTv:
		topic = messaging.TopicTv
	default:
		return ErrUnknownDeviceType
	}

	msg, err := messaging.NewCommandMessage(device.ID, command, detail)
	if err != nil {
		return err
	}

	return messaging.Publish(topic, msg)
}

// deviceIsCapable determines if the device with the given ID can perform the given command
func (c *CommandController) deviceIsCapable(device *config.Device, command string) (bool, error) {
	for _, v := range device.Capabilities {
		if v == command {
			return true, nil
		}
	}

	return false, nil
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
