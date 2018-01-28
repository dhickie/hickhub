package config

// MessagingConfig represents messaging configuration for receiving messages from the internet
type MessagingConfig struct {
	MessagingServer string `json:"messaging_server"`
	APIServer       string `json:"api_server"`
	AuthKey         string `json:"auth_key"`
}
