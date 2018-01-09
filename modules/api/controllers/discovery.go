package controllers

import "github.com/dhickie/hickhub/config"

// DiscoveryController encapsulates all API endpoints relating to device discovery
type DiscoveryController struct {
	devices []config.Device
}

// NewDiscoveryController creates a new discovery controller from the given HickHub config
func NewDiscoveryController(config config.Config) *DiscoveryController {
	return &DiscoveryController{config.Devices}
}

// GetDevices gets all the devices which are controlled by this HickHub
func (c *DiscoveryController) GetDevices() []config.Device {
	return c.devices
}
