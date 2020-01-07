package backend

import (
	"encoding/gob"
	"golang.org/x/oauth2"
)


func init() {
	gob.Register(&oauth2.Token{})

}

func (s *Server) routes() {
	s.router.Handle("/", s.handleHome()).Methods("GET")
	s.router.Handle("/health", s.handleHealthProbe()).Methods("GET")
}

