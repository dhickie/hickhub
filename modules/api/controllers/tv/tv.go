package tv

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dhickie/openhub/messaging"
	"github.com/gorilla/mux"
)

// TurnOff sends a message to turn off the TV
func TurnOff(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeTurnOff, nil)
	publishMessage(w, r, msg)
}

// ChannelUp sends a message to go up one channel
func ChannelUp(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeChannelUp, nil)
	publishMessage(w, r, msg)
}

// ChannelDown sends a message to go down one channel
func ChannelDown(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeChannelDown, nil)
	publishMessage(w, r, msg)
}

// SetChannel sends a message to set the TV to the specified channel
func SetChannel(w http.ResponseWriter, r *http.Request) {
	channel, err := strconv.Atoi(mux.Vars(r)["channel"])
	if err != nil {
		badRequest(w, err)
		return
	}

	msg := messaging.NewMessage(messaging.MessageTypeSetChannel, channel)
	publishMessage(w, r, msg)
}

// VolumeUp sends a message to increase the volume by one
func VolumeUp(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeVolumeUp, nil)
	publishMessage(w, r, msg)
}

// VolumeDown sends a message to decrease the volume by one
func VolumeDown(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeVolumeDown, nil)
	publishMessage(w, r, msg)
}

// SetVolume sets the current volume of the TV to the provided value
func SetVolume(w http.ResponseWriter, r *http.Request) {
	volume, err := strconv.Atoi(mux.Vars(r)["volume"])
	if err != nil {
		badRequest(w, err)
		return
	}

	msg := messaging.NewMessage(messaging.MessageTypeSetVolume, volume)
	publishMessage(w, r, msg)
}

// Play sends a message to do play the TV
func Play(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypePlay, nil)
	publishMessage(w, r, msg)
}

// Pause sends a message to pause the TV
func Pause(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypePause, nil)
	publishMessage(w, r, msg)
}

// Stop sends a message to stop the TV
func Stop(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeStop, nil)
	publishMessage(w, r, msg)
}

func publishMessage(w http.ResponseWriter, r *http.Request, msg messaging.Message) {
	deviceID := mux.Vars(r)["deviceId"]
	msg.DeviceID = deviceID
	err := messaging.Publish(messaging.TopicTv, msg)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprint(err), 400)
}
