package api

import (
	"fmt"
	"net/http"

	"github.com/dhickie/openhub/config"
	"github.com/dhickie/openhub/log"
	"github.com/dhickie/openhub/modules/api/controllers/tv"
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
	r := setupRoutes()
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), middleware{r})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to launch http listener: %v", err.Error()))
	}
}

func setupRoutes() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/").Subrouter()

	tv := s.PathPrefix("/tv/{deviceId}/").Subrouter()
	setupTvRoutes(tv)

	return r
}

func setupTvRoutes(r *mux.Router) {
	r.HandleFunc("/power/off", tv.TurnOff).Methods("POST")

	r.HandleFunc("/channel/up", tv.ChannelUp).Methods("POST")
	r.HandleFunc("/channel/down", tv.ChannelDown).Methods("POST")
	r.HandleFunc("/channel/set/{channel}", tv.SetChannel).Methods("POST")

	r.HandleFunc("/volume/up", tv.VolumeUp).Methods("POST")
	r.HandleFunc("/volume/down", tv.VolumeDown).Methods("POST")
	r.HandleFunc("/volume/set/{volume}", tv.SetVolume).Methods("GET")

	r.HandleFunc("/media/play", tv.Play).Methods("POST")
	r.HandleFunc("/media/pause", tv.Pause).Methods("POST")
	r.HandleFunc("/media/stop", tv.Stop).Methods("POST")
}
