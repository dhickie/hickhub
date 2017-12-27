package tv

import (
	"github.com/dhickie/go-lgtv/control"
	"github.com/dhickie/openhub/config"
	"github.com/dhickie/openhub/messaging"
)

var controller tvController

// Launch launches the TV module using the specified config
func Launch(appConfig config.Config) {
	tvConfig := appConfig.Tv

	// Create the TV objects using the provided config
	tvMap := make(map[string]*control.LgTv)
	for _, v := range tvConfig.Tvs {
		tv, err := control.NewTV(v.IPAddress)
		if err != nil {
			continue
		}

		tvMap[v.ID] = &tv
	}

	// Connect to each TV using client key from the TV config
	for _, v := range tvConfig.Tvs {
		tvMap[v.ID].Connect(v.ClientKey)
	}

	// Setup the controller
	controller = tvController{
		Tvs: tvMap,
	}

	// Subscribe the controller to messages on the TV topic
	messaging.Subscribe(messaging.TopicTv, controller.subscriber)
}
