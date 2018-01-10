package tv

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dhickie/go-lgtv/control"
	"github.com/dhickie/hickhub/log"
	"github.com/dhickie/hickhub/messaging"
	"github.com/dhickie/hickhub/messaging/payloads"
	"github.com/dhickie/hickhub/models"
)

// tvController controlls all TVs under its remit when an appropriate message is received
type tvController struct {
	Tvs        map[string]*control.LgTv
	ClientKeys map[string]string
}

// subscriber is the callback called when the TV module receives a message
func (c *tvController) subscriber(msg messaging.Message) {
	// We know this is a command message, so unmarshal the payload as such
	cmd := new(payloads.CommandPayload)
	err := json.Unmarshal([]byte(msg.Payload), cmd)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured unmarshalling the command payload: %v", err))
		return
	}

	// Perform the provided command on the TV with the given device ID
	tv, ok := c.Tvs[cmd.DeviceID]
	if ok {
		switch cmd.Command {
		case models.CommandTurnOff:
			err = tv.TurnOff()
		case models.CommandVolumeUp:
			err = tv.VolumeUp()
		case models.CommandVolumeDown:
			err = tv.VolumeDown()
		case models.CommandSetVolume:
			val, err := strconv.Atoi(cmd.Detail)
			if err != nil {
				log.Error(fmt.Sprintf("An error occured getting the target volume: %v", err))
				return
			}
			err = tv.SetVolume(val)
		case models.CommandChannelUp:
			err = tv.ChannelUp()
		case models.CommandChannelDown:
			err = tv.ChannelDown()
		case models.CommandSetChannel:
			val, err := strconv.Atoi(cmd.Detail)
			if err != nil {
				log.Error(fmt.Sprintf("An error occured getting the target channel number: %v", err))
				return
			}
			err = tv.SetChannel(val)
		}

		if err != nil {
			log.Error(fmt.Sprintf("An error occured performing the requested TV operation: %v", err))
		}
	} else {
		log.Error(fmt.Sprintf("Received message for unknown device ID: %v", cmd.DeviceID))
	}
}
