package controllers

import "github.com/dhickie/hickhub/config"

// ControlController encapsulates all information related to controlling a device
type ControlController struct {
	devices []config.Device
}

// NewControlController returns a new ControlController instance
func NewControlController(config config.Config) *ControlController {
	return &ControlController{config.Devices}
}
