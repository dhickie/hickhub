package tv

import (
	"fmt"
	"net/http"

	"github.com/dhickie/openhub/messaging"
	"github.com/gorilla/mux"
)

// TurnOff sends a message to turn off the TV
func TurnOff(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeTurnOff, nil)
	publishMessage(w, msg)
}

// ChannelUp sends a message to go up one channel
func ChannelUp(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeChannelUp, nil)
	publishMessage(w, msg)
}

// ChannelDown sends a message to go down one channel
func ChannelDown(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeChannelDown, nil)
	publishMessage(w, msg)
}

// SetChannel sends a message to set the TV to the specified channel
func SetChannel(w http.ResponseWriter, r *http.Request) {
	channel := mux.Vars(r)["channel"]
	msg := messaging.NewMessage(messaging.MessageTypeSetChannel, channel)
	publishMessage(w, msg)
}

// VolumeUp sends a message to increase the volume by one
func VolumeUp(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeVolumeUp, nil)
	publishMessage(w, msg)
}

// VolumeDown sends a message to decrease the volume by one
func VolumeDown(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeVolumeDown, nil)
	publishMessage(w, msg)
}

func SetVolume(w http.ResponseWriter, r *http.Request) {
	volume := mux.Vars(r)["volume"]
	msg := messaging.NewMessage(messaging.MessageTypeSetVolume, volume)
	publishMessage(w, msg)
}

// Play sends a message to do play the TV
func Play(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypePlay, nil)
	publishMessage(w, msg)
}

// Pause sends a message to pause the TV
func Pause(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypePause, nil)
	publishMessage(w, msg)
}

// Stop sends a message to stop the TV
func Stop(w http.ResponseWriter, r *http.Request) {
	msg := messaging.NewMessage(messaging.MessageTypeStop, nil)
	publishMessage(w, msg)
}

func publishMessage(w http.ResponseWriter, msg messaging.Message) {
	err := messaging.Publish(messaging.TopicTv, msg)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
}
