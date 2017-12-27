package messaging

// Represents the different type of messages that can be sent between modules
const (
	MessageTypeTurnOff     = "turn_off"
	MessageTypeChannelUp   = "channel_up"
	MessageTypeChannelDown = "channel_down"
	MessageTypeSetChannel  = "set_channel"
	MessageTypeVolumeUp    = "volume_up"
	MessageTypeVolumeDown  = "volume_down"
	MessageTypeSetVolume   = "set_volume"
	MessageTypePlay        = "play"
	MessageTypePause       = "pause"
	MessageTypeStop        = "stop"
)

// Message is a message sent between different modules
type Message struct {
	Type     string      `json:"type"`
	DeviceID string      `json:"device_id"`
	Payload  interface{} `json:"payload"` // The payload will vary depending on the type of message being sent
}

// NewMessage returns a new message of the specified type with the specified payload
func NewMessage(msgType string, payload interface{}) Message {
	return Message{
		Type:    msgType,
		Payload: payload,
	}
}
