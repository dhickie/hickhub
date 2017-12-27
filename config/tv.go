package config

// TvConfig contains the configuration needed for the TV module
type TvConfig struct {
	Tvs []TvInfo `json:"tvs"`
}

// TvInfo contains the information needed to connect to a TV
type TvInfo struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip_address"`
	ClientKey string `json:"client_key"`
}
