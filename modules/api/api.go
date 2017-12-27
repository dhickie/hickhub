package api

import (
	"fmt"
	"net/http"

	"github.com/dhickie/openhub/config"
	"github.com/dhickie/openhub/modules/api/controllers/tv"
	"github.com/gorilla/mux"
)

// Launch configures and then launches the API module
func Launch(appConfig config.Config) {
	config := appConfig.API

	r := setupRoutes()
	http.ListenAndServe(fmt.Sprintf(":%v", config.Port), r)
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
