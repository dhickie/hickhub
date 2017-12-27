package main

import (
	"github.com/dhickie/openhub/config"
	"github.com/dhickie/openhub/modules/api"
	"github.com/dhickie/openhub/modules/tv"
)

func main() {
	config := config.Config{
		API: config.APIConfig{
			Port: 10001,
		},
		Tv: config.TvConfig{
			Tvs: []config.TvInfo{
				config.TvInfo{
					ID:        "tv1",
					IPAddress: "192.168.1.130",
					ClientKey: "",
				},
			},
		},
	}

	// Launch the Modules
	go api.Launch(config)
	go tv.Launch(config)

	select {}
}
