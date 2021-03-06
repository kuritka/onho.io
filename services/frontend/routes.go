package frontend

import (
	"encoding/gob"
	"net/http"

	"golang.org/x/oauth2"
)


const (
	sessionStoreKey = "fugu"
	staticDir = "./services/frontend/"
)

func init() {
	gob.Register(&oauth2.Token{})

}

func (s *Server) routes() {
	// register handlers here
	s.router.PathPrefix("/static/").Handler(notFoundHook{http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir+"static")))})

	s.router.Handle("/", s.handleHome(sessionStoreKey)).Methods("GET")
	s.router.Handle("/api/auth/start", s.handleStart(sessionStoreKey)).Methods("GET")
	s.router.HandleFunc("/api/auth/signing-callback",s.handleAuthCallback(sessionStoreKey)).Methods("GET")
	s.router.HandleFunc("/api/auth/destroy-session", s.handleDestroySession(sessionStoreKey)).Methods("GET")
	s.router.HandleFunc("/api/websockets/stream", s.serveWebSockets(s.commandPublisher))
	s.router.Handle("/health", s.handleHealthProbe()).Methods("GET")
}


type hookedResponseWriter struct {
	http.ResponseWriter
	ignore bool
}

func (hrw *hookedResponseWriter) writeHeader(status int) {
	hrw.ResponseWriter.WriteHeader(status)
	if status == 404 {
		hrw.ignore = true
	}
}

func (hrw *hookedResponseWriter) write(p []byte) (int, error) {
	if hrw.ignore {
		return len(p), nil
	}
	return hrw.ResponseWriter.Write(p)
}

type notFoundHook struct {
	h http.Handler
}

func (nfh notFoundHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nfh.h.ServeHTTP(&hookedResponseWriter{ResponseWriter: w}, r)
}
