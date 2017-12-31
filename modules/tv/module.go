package tv

import (
	"fmt"

	"github.com/dhickie/go-lgtv/control"
	"github.com/dhickie/hickhub/config"
	"github.com/dhickie/hickhub/log"
	"github.com/dhickie/hickhub/messaging"
)

var controller tvController

// Launch launches the TV module using the specified config
func Launch(appConfig config.Config) {
	log.Info("Launching TV module")
	tvConfig := appConfig.Tv

	// Create the TV objects using the provided config
	log.Info("Creating TV objects")
	tvMap := make(map[string]*control.LgTv)
	for _, v := range tvConfig.Tvs {
		tv, err := control.NewTV(v.IPAddress)
		if err != nil {
			log.Error(fmt.Sprintf("Error creating TV at %v: %v", v.IPAddress, err.Error()))
			continue
		}

		tvMap[v.ID] = &tv
	}

	// Connect to each TV using client key from the TV config
	log.Info("Connecting to TVs")
	for _, v := range tvConfig.Tvs {
		_, err := tvMap[v.ID].Connect(v.ClientKey)
		if err != nil {
			log.Error(fmt.Sprintf("An error occured connecting to TV at %v: %v", v.IPAddress, err.Error()))
		}
	}

	// Setup the controller
	controller = tvController{
		Tvs: tvMap,
	}

	// Subscribe the controller to messages on the TV topic
	log.Info("Subscribing to TV messaging topic")
	messaging.Subscribe(messaging.TopicTv, controller.subscriber)
}
