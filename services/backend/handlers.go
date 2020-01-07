package backend

import (
	"github.com/kuritka/onho.io/common/utils"
	"net/http"
)

func (s *Server) handleHome()http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write( []byte("Hello from backend"))
		utils.FailFastOnError(err)
	}
}




func (s *Server) handleHealthProbe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

