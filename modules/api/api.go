package api

import (
	"fmt"
	"net/http"

	"github.com/dhickie/hickhub/config"
	"github.com/dhickie/hickhub/log"
	"github.com/dhickie/hickhub/modules/api/controllers"
	"github.com/gorilla/mux"
)

type middleware struct {
	h http.Handler
}

func (m middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log details of the incoming request
	log.Info(fmt.Sprintf("Incoming HTTP request: Path - %v, IP Address - %v", r.URL.Path, r.RemoteAddr))
	m.h.ServeHTTP(w, r)
}

// Launch configures and then launches the API module
func Launch(appConfig config.Config) {
	log.Info("Launching API module")
	config := appConfig.API

	log.Info("Setting up routes")
	r := setupRoutes(appConfig)
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), middleware{r})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to launch http listener: %v", err.Error()))
	}
}

func setupRoutes(appConfig config.Config) *mux.Router {
	disc := controllers.NewDiscoveryController(appConfig)
	cmd := controllers.NewCommandController(appConfig)

	r := mux.NewRouter()
	s := r.PathPrefix("/api/").Subrouter()

	s.HandleFunc("/devices", disc.GetDevices).Methods("GET")
	s.HandleFunc(`/device/{id}/command/{cmd:[a-zA-Z0-9=\-\/]+}`, cmd.ControlDevice).Methods("POST")

	return r
}
