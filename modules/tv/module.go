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
	devices := appConfig.Devices

	// Create the TV objects using the provided config
	log.Info("Creating TV objects")
	tvMap := make(map[string]*control.LgTv)
	keyMap := make(map[string]string)

	for _, v := range devices {
		if v.Type == config.TypeTv && v.SubType == config.SubTypeWebOsTv {
			info := v.Info.(*config.WebOsTvDeviceInfo)
			tv, err := control.NewTV(info.IPAddress)
			if err != nil {
				log.Error(fmt.Sprintf("Error creating TV at %v: %v", info.IPAddress, err.Error()))
				continue
			}

			// Try to connect to the TV using the client key from the config
			_, err = tv.Connect(info.ClientKey)
			if err != nil {
				// No biggie if we can't connect - it might not be turned on
				log.Warn(fmt.Sprintf("Unable to connect to TV at %v: %v", info.IPAddress, err.Error()))
			}

			tvMap[v.ID] = &tv
			keyMap[v.ID] = info.ClientKey
		}
	}

	// Setup the controller
	controller = tvController{
		Tvs:        tvMap,
		ClientKeys: keyMap,
	}

	// Subscribe the controller to messages on the TV topic
	log.Info("Subscribing to TV messaging topic")
	messaging.Subscribe(messaging.TopicTv, controller.subscriber)
}
