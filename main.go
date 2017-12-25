package main

import (
	"github.com/dhickie/openhub/config"
	"github.com/dhickie/openhub/modules/api"
)

func main() {
	// Launch the API
	go api.Launch(config.Config{
		API: config.APIConfig{
			Port: 10001,
		},
	})

	select {}
}
