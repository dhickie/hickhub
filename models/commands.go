package models

// Constants for all the possible types of command
const (
	CommandOn          = "on"
	CommandOff         = "off"
	CommandUp          = "up"
	CommandDown        = "down"
	CommandSet         = "set"
	CommandAdjust      = "adjust"
	CommandSetMute     = "setmute"
	CommandPlay        = "play"
	CommandPause       = "pause"
	CommandRewind      = "rewind"
	CommandFastForward = "fastforward"
	CommandLaunch      = "launch"
)

// SetChannelDetail is sent with commands to set the current channel
type SetChannelDetail struct {
	ExactChannelNumber     int    `json:"exact_channel_number"`
	ExactChannelName       string `json:"exact_channel_name"`
	FuzzyChannelIdentifier string `json:"fuzzy_channel_identifier"`
}
