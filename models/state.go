package models

// Constants for all the possible types of state to be read or modified
const (
	StatePower    = "power"
	StateVolume   = "volume"
	StateChannel  = "channel"
	StatePlayback = "playback"
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

// PowerState represents the current power state of a device
type PowerState struct {
	PowerOn bool `json:"power_on"`
}

// ChannelState represents the current channel state of a device
type ChannelState struct {
	ChannelNumber int    `json:"channel_number"`
	ChannelName   string `json:"channel_name"`
}
