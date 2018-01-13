package models

// Constants for all the possible types of state to be read or modified
const (
	StatePower   = "power"
	StateVolume  = "volume"
	StateChannel = "channel"
)

// DeviceState represents the current state of a particular state component of the device
type DeviceState struct {
	Type  string      `json:"type"`
	State interface{} `json:"state"`
}

// VolumeState represents the state of a device with the volume state
type VolumeState struct {
	Volume  int  `json:"volume"`
	IsMuted bool `json:"is_muted"`
}
