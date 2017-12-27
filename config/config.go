package config

// Config represents the overall config of the application
type Config struct {
	API APIConfig `json:"api"`
	Tv  TvConfig  `json:"tv"`
}
