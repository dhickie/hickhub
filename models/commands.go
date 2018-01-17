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
)

// SetChannelDetail is sent with commands to set the current channel
type SetChannelDetail struct {
	ChannelNumber int    `json:"channel_number"`
	ChannelName   string `json:"channel_name"`
}
