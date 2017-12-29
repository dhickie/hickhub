package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dhickie/openhub/log"
	"github.com/dhickie/openhub/messaging"
	"github.com/gorilla/mux"
)

// TvController is an empty struct for representing the TV API Controller
type TvController struct {
}

// TurnOff sends a message to turn off the TV
func (t *TvController) TurnOff(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeTurnOff, nil)
	t.publishMessage(w, r, msg)
}

// ChannelUp sends a message to go up one channel
func (t *TvController) ChannelUp(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeChannelUp, nil)
	t.publishMessage(w, r, msg)
}

// ChannelDown sends a message to go down one channel
func (t *TvController) ChannelDown(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeChannelDown, nil)
	t.publishMessage(w, r, msg)
}

// SetChannel sends a message to set the TV to the specified channel
func (t *TvController) SetChannel(w http.ResponseWriter, r *http.Request) {
	channel, err := strconv.Atoi(mux.Vars(r)["channel"])
	if err != nil {
		log.Error(fmt.Sprintf("Bad request - Failed to convert channel to integer: %v", err.Error()))
		t.badRequest(w, err)
		return
	}

	msg := messaging.NewMessage(messaging.MessageTypeSetChannel, channel)
	t.publishMessage(w, r, msg)
}

// VolumeUp sends a message to increase the volume by one
func (t *TvController) VolumeUp(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeVolumeUp, nil)
	t.publishMessage(w, r, msg)
}

// VolumeDown sends a message to decrease the volume by one
func (t *TvController) VolumeDown(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeVolumeDown, nil)
	t.publishMessage(w, r, msg)
}

// SetVolume sets the current volume of the TV to the provided value
func (t *TvController) SetVolume(w http.ResponseWriter, r *http.Request) {
	volume, err := strconv.Atoi(mux.Vars(r)["volume"])
	if err != nil {
		log.Error(fmt.Sprintf("Bad request - failed to convert volume to integer: %v", err.Error()))
		t.badRequest(w, err)
		return
	}

	msg := messaging.NewMessage(messaging.MessageTypeSetVolume, volume)
	t.publishMessage(w, r, msg)
}

// Play sends a message to do play the TV
func (t *TvController) Play(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypePlay, nil)
	t.publishMessage(w, r, msg)
}

// Pause sends a message to pause the TV
func (t *TvController) Pause(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypePause, nil)
	t.publishMessage(w, r, msg)
}

// Stop sends a message to stop the TV
func (t *TvController) Stop(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeStop, nil)
	t.publishMessage(w, r, msg)
}

func (t *TvController) publishMessage(w http.ResponseWriter, r *http.Request, msg messaging.Message) {
	deviceID := mux.Vars(r)["deviceId"]
	msg.DeviceID = deviceID
	err := messaging.Publish(messaging.TopicTv, msg)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured publishing the message to the TV topic: %v", err.Error()))
		http.Error(w, fmt.Sprint(err), 500)
	}
}

func (t *TvController) badRequest(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprint(err), 400)
}
