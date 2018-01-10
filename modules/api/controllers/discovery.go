package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dhickie/hickhub/config"
)

// ErrMarshalError is sent when an error occurs marshaling the JSON response
var ErrMarshalError = errors.New("An error occured marshaling the response payload")

// DiscoveryController encapsulates all API endpoints relating to device discovery
type DiscoveryController struct {
	devices []config.Device
}

// NewDiscoveryController creates a new discovery controller from the given HickHub config
func NewDiscoveryController(config config.Config) *DiscoveryController {
	return &DiscoveryController{config.Devices}
}

// GetDevices gets all the devices which are controlled by this HickHub
func (c *DiscoveryController) GetDevices(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(c.devices)
	if err != nil {
		http.Error(w, ErrMarshalError.Error(), 500)
		return
	}

	w.Write(response)
}
