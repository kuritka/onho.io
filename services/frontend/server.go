package frontend

import (
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

const (
	defaultLayout = "./services/frontend/templates/layout.html"
	templateDir   = "./services/frontend/templates/"
)

type (
	Idp int
)

const (
	GitHubProvider    Idp = iota
	//GoogleHubProvider
)

type Server struct {
	router           *mux.Router
	oauthCfg         *oauth2.Config
	store            *sessions.CookieStore
	templates        map[string]*template.Template
	upgrader         websocket.Upgrader
	login            *string
	commandPublisher func(data string)
}

func NewServer( mux *mux.Router, cookieStore *sessions.CookieStore, oauthCfg *oauth2.Config ,f func(data string) ) *Server {
	templates := map[string]*template.Template{}
	templates["home.html"] = template.Must(template.ParseFiles(templateDir+"home.html", defaultLayout))
	upgrader := websocket.Upgrader{ ReadBufferSize:  1024, WriteBufferSize: 1024,}
	server := Server{mux, oauthCfg,cookieStore, templates,upgrader, nil, f}
	server.routes()
	return &server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//https://neoteric.eu/blog/how-to-serve-static-files-with-golang/
	s.router.ServeHTTP(w, r)
}
