package main

import (
	"github.com/dhickie/openhub/config"
	"github.com/dhickie/openhub/modules/api"
	"github.com/dhickie/openhub/modules/logging"
	"github.com/dhickie/openhub/modules/tv"
)

func main() {
	// Read in the current configuration
	config, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// The logging module is a special case - launch it
	// synchonously before moving on to the other modules
	logging.Launch(config)

	// Launch the Modules
	go api.Launch(config)
	go tv.Launch(config)

	select {}
}
