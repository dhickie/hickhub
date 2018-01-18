package tv

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dhickie/go-lgtv/control"
	"github.com/dhickie/hickhub/log"
	"github.com/dhickie/hickhub/messaging"
	"github.com/dhickie/hickhub/messaging/payloads"
	"github.com/dhickie/hickhub/models"
	"github.com/dhickie/hickhub/utils"
)

// ErrCommandUnsupported is returned when a requested command is not supported by the TV
var ErrCommandUnsupported = errors.New("Requested action is unsupported")

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
	success := false
	errStr := ""
	deviceState := models.DeviceState{}

	if tv, ok := c.Tvs[cmd.DeviceID]; ok {
		clientKey := c.ClientKeys[cmd.DeviceID]

		// If this isn't a command to turn the TV on, then make sure we're connected first
		if cmd.State != models.StatePower || cmd.Command != models.CommandOn {
			newKey, err := tv.Connect(clientKey, 1000)
			if err != nil {
				errStr = fmt.Sprintf("An error occured trying to connect to the TV: %v", err)
			} else if newKey != clientKey {
				c.ClientKeys[cmd.DeviceID] = newKey
			}
		}

		// Figure out what state the command was for, if we didn't error trying to connect
		if errStr == "" {
			switch cmd.State {
			case models.StateVolume:
				err = handleVolumeCommand(tv, cmd.Command, cmd.Detail)
				if err == nil {
					var volState models.VolumeState
					volState, err = getVolumeState(tv)

					deviceState.Type = models.StateVolume
					deviceState.State = volState
				}
			case models.StateChannel:
				err = handleChannelCommand(tv, cmd.Command, cmd.Detail)
				if err == nil {
					var chanState models.ChannelState
					chanState, err = getChannelState(tv)

					deviceState.Type = models.StateChannel
					deviceState.State = chanState
				}
			case models.StatePower:
				err = handlePowerCommand(tv, clientKey, cmd.Command, cmd.Detail)
				if err == nil {
					powerOn := false
					if cmd.Command == models.CommandOn {
						powerOn = true
					}
					deviceState.Type = models.StatePower
					deviceState.State = models.PowerState{PowerOn: powerOn}
				}
			case models.StatePlayback:
				err = handlePlaybackCommand(tv, cmd.Command, cmd.Detail)
				// There's no device state for playback
			case models.StateInput:
				name, err := handleInputCommand(tv, cmd.Command, cmd.Detail)
				if err == nil {
					deviceState.Type = models.StateInput
					deviceState.State = models.InputState{InputName: name}
				}
			}

			if err != nil {
				errStr = fmt.Sprintf("An error occured performing the requested TV operation: %v", err)
				log.Error(errStr)
			} else {
				success = true
			}
		}
	} else {
		errStr = fmt.Sprintf("Received message for unknown device ID: %v", cmd.DeviceID)
		log.Error(errStr)
	}

	// Build the result payload and send it back
	resultMessage, err := messaging.NewCommandResultMessage(success, errStr, deviceState)
	if err != nil {
		// Log the error, we can't publish the result back :(
		log.Error(fmt.Sprintf("An error occured trying to create the result message: %v", err))
		return
	}

	messaging.Publish(msg.Reply, resultMessage)
}

func handleVolumeCommand(tv *control.LgTv, command string, detail string) error {
	switch command {
	case models.CommandSetMute:
		isMute, err := strconv.ParseBool(detail)
		if err != nil {
			return err
		}
		return tv.SetMute(isMute)
	case models.CommandUp:
		return tv.VolumeUp()
	case models.CommandDown:
		return tv.VolumeDown()
	case models.CommandSet:
		var val int
		err := json.Unmarshal([]byte(detail), &val)
		if err != nil {
			return err
		}
		return tv.SetVolume(val)
	case models.CommandAdjust:
		// Get how much we want to adjust by
		var val int
		err := json.Unmarshal([]byte(detail), &val)
		if err != nil {
			return err
		}

		// Get the current volume
		currentVol, err := tv.GetVolume()
		if err != nil {
			return err
		}

		// Set it to the new value
		newVol := currentVol + val
		if newVol < 0 {
			newVol = 0
		} else if newVol > 100 {
			newVol = 100
		}
		return tv.SetVolume(newVol)
	}

	return ErrCommandUnsupported
}

func getVolumeState(tv *control.LgTv) (models.VolumeState, error) {
	// Get the volume
	vol, err := tv.GetVolume()
	if err != nil {
		return models.VolumeState{}, err
	}

	// Get the mute status
	isMute, err := tv.GetMute()
	if err != nil {
		return models.VolumeState{}, err
	}

	return models.VolumeState{
		Volume:  vol,
		IsMuted: isMute,
	}, nil
}

func handleChannelCommand(tv *control.LgTv, command string, detail string) error {
	switch command {
	case models.CommandUp:
		return tv.ChannelUp()
	case models.CommandDown:
		return tv.ChannelDown()
	case models.CommandSet:
		val := new(models.SetChannelDetail)
		err := json.Unmarshal([]byte(detail), val)
		if err != nil {
			return err
		}

		// Get the list of possible channels
		channels, err := tv.ListChannels()
		if err != nil {
			return err
		}

		// Also get the list of possible apps
		apps, err := tv.ListInstalledApps()
		if err != nil {
			return err
		}

		// Match the request to closest channel we can find
		targetChannel, chanErr := utils.MatchChannel(*val, channels)
		// Also try to match to an app, like netflix
		targetApp, appErr := utils.MatchApp(val.FuzzyChannelIdentifier, apps)

		if chanErr != nil && appErr != nil {
			return chanErr
		} else if appErr != nil {
			return tv.SetChannel(targetChannel.ChannelNumber)
		}

		_, err = tv.LaunchApp(targetApp.ID)
		return err
	case models.CommandAdjust:
		// Work out how many channels to change by
		var val int
		err := json.Unmarshal([]byte(detail), &val)
		if err != nil {
			return err
		}

		// Get the current channel
		channel, err := tv.GetCurrentChannel()
		if err != nil {
			return err
		}

		// Adjust the channel
		newChannel := channel.ChannelNumber + val
		if newChannel < 0 {
			newChannel = 0
		}
		return tv.SetChannel(newChannel)
	}

	return ErrCommandUnsupported
}

func getChannelState(tv *control.LgTv) (models.ChannelState, error) {
	// Get the current channel
	channel, err := tv.GetCurrentChannel()
	if err != nil {
		return models.ChannelState{}, err
	}

	return models.ChannelState{
		ChannelName:   channel.ChannelName,
		ChannelNumber: channel.ChannelNumber,
	}, nil
}

func handlePowerCommand(tv *control.LgTv, clientKey, command, detail string) error {
	switch command {
	case models.CommandOff:
		return tv.TurnOff()
	case models.CommandOn:
		err := tv.TurnOn()
		if err != nil {
			return err
		}

		// Try to connect once a second, for 4 seconds
		ticker := time.NewTicker(1 * time.Second)
		for i := 0; i < 4; i++ {
			<-ticker.C
			_, err = tv.Connect(clientKey, 500)
			if err == nil {
				return nil
			}
		}

		return err
	}

	return ErrCommandUnsupported
}

func handlePlaybackCommand(tv *control.LgTv, command string, detail string) error {
	switch command {
	case models.CommandPlay:
		return tv.Play()
	case models.CommandPause:
		return tv.Pause()
	case models.CommandFastForward:
		return tv.FastForward()
	case models.CommandRewind:
		return tv.Rewind()
	}

	return ErrCommandUnsupported
}

func handleInputCommand(tv *control.LgTv, command string, detail string) (string, error) {
	switch command {
	case models.CommandSet:
		// Get the current list of possible inputs
		inputs, err := tv.ListExternalInputs()
		if err != nil {
			return "", err
		}

		// Find the match closest to the specified target
		var target string
		err = json.Unmarshal([]byte(detail), &target)
		if err != nil {
			return "", err
		}
		match, err := utils.MatchInput(target, inputs)
		if err != nil {
			return "", err
		}

		// Set the input to the closest match
		err = tv.SwitchInput(match.ID)
		return match.ID, err
	}

	return "", ErrCommandUnsupported
}
