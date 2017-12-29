package tv

import (
	"fmt"

	"github.com/dhickie/go-lgtv/control"
	"github.com/dhickie/openhub/log"
	"github.com/dhickie/openhub/messaging"
)

// tvController controlls all TVs under its remit when an appropriate message is received
type tvController struct {
	Tvs map[string]*control.LgTv
}

// subscriber is the callback called when the TV module receives a message
func (c *tvController) subscriber(msg messaging.Message) {
	// If we don't know about the device with this ID, then don't do anything
	log.Info(fmt.Sprintf("TV Controller received message: Type - %v, Device: %v", msg.Type, msg.DeviceID))
	tv, ok := c.Tvs[msg.DeviceID]
	if ok {
		// Work out which command we need to do based on the type of message
		switch msg.Type {
		case messaging.MessageTypeTurnOff:
			tv.TurnOff()
			break
		case messaging.MessageTypeVolumeUp:
			tv.VolumeUp()
			break
		case messaging.MessageTypeVolumeDown:
			tv.VolumeDown()
			break
		case messaging.MessageTypeSetVolume:
			tv.SetVolume(msg.Payload.(int))
			break
		case messaging.MessageTypeChannelUp:
			tv.ChannelUp()
			break
		case messaging.MessageTypeChannelDown:
			tv.ChannelDown()
			break
		case messaging.MessageTypeSetChannel:
			tv.SetChannel(msg.Payload.(int))
			break
		}
	} else {
		log.Warn(fmt.Sprintf("Received message for unknown device ID: %v", msg.DeviceID))
	}
}
