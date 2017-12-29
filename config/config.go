package config

import "io/ioutil"
import "encoding/json"

// Config represents the overall config of the application
type Config struct {
	API APIConfig `json:"api"`
	Tv  TvConfig  `json:"tv"`
}

// ReadConfig reads the config from the configuration JSON file
func ReadConfig() (Config, error) {
	config := new(Config)
	contents, err := ioutil.ReadFile("config.json")
	if err != nil {
		return *config, err
	}

	err = json.Unmarshal(contents, config)
	return *config, err
}
