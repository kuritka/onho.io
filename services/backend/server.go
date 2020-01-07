package backend

import (

	"net/http"

	"github.com/gorilla/mux"

)


type Server struct {
	router           *mux.Router

}

func NewServer( mux *mux.Router ) *Server {
	server := Server{mux}
	server.routes()
	return &server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//https://neoteric.eu/blog/how-to-serve-static-files-with-golang/
	s.router.ServeHTTP(w, r)
}
