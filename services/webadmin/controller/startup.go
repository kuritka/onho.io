package controller

import (
	"net/http"
)

type startup struct {
	options Options
	websocket *websocketController
}

func NewStartup(options Options) *startup {
	websocket := newWebsocketController(options)
	return &startup{
		options:options,
		websocket: websocket,
	}
}

func (wa *startup) Init() {
	wa.registerRoutes()
	wa.registerFileServers()
}

func (wa *startup) registerRoutes() {
	http.HandleFunc("/ws", wa.websocket.handleMessage)
}

func (wa *startup)  registerFileServers() {
	http.Handle("/public/", http.FileServer(http.Dir("services/webadmin/assets")))
	http.Handle("/public/lib/",
		http.StripPrefix("/public/lib/", http.FileServer(http.Dir("services/webadmin/node_modules"))))

}