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

	tv := s.PathPrefix("/tv/").Subrouter()
	setupTvRoutes(tv)

	return r
}

func setupTvRoutes(r *mux.Router) {
	r.HandleFunc("/power/off/", tv.TurnOff)

	r.HandleFunc("/channel/up/", tv.ChannelUp)
	r.HandleFunc("/channel/down/", tv.ChannelDown)
	r.HandleFunc("/channel/set/{channel}/", tv.SetChannel)

	r.HandleFunc("/volume/up/", tv.VolumeUp)
	r.HandleFunc("/volume/down/", tv.VolumeDown)
	r.HandleFunc("/volume/set/{volume}/", tv.VolumeDown)

	r.HandleFunc("/media/play/", tv.Play)
	r.HandleFunc("/media/pause/", tv.Pause)
	r.HandleFunc("/media/stop/", tv.Stop)
}
