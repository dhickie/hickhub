package config

// Const values for the different types of device
const (
	TypeTv = "TV"
)

// Const values for the different sub types of device
const (
	SubTypeWebOsTv = "WebOsTV"
)

// Device represents a device which is controllable from the HickHub
type Device struct {
	Type         string      `json:"type"`
	SubType      string      `json:"sub_type"`
	ID           string      `json:"id"`
	Capabilities []string    `json:"capabilities"`
	Info         interface{} `json:"info"`
}

// WebOsTvDeviceInfo represents the extra information needed by a WebOS TV
type WebOsTvDeviceInfo struct {
	IPAddress string `json:"ip_address"`
	ClientKey string `json:"client_key"`
}
