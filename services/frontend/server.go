package frontend

import (
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

const (
	defaultLayout = "services/frontend/templates/layout.html"
	templateDir   = "services/frontend/templates/"
)

const (
	githubAuthorizeUrl = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
	redirectUrl        = ""
)

type (
	Idp int
)

const (
	GitHubProvider    Idp = iota
	//GoogleHubProvider
)

type Server struct {
	router    *mux.Router
	oauthCfg  *oauth2.Config
	store     *sessions.CookieStore
	templates map[string]*template.Template
	upgrader  websocket.Upgrader
	login 	  *string
}

func NewServer( mux *mux.Router, config *Options, oauthCfg *oauth2.Config ) *Server {
	cookieStore := sessions.NewCookieStore([]byte(config.CookieStoreKey))
	templates := map[string]*template.Template{}
	templates["home.html"] = template.Must(template.ParseFiles(templateDir+"home.html", defaultLayout))
	upgrader := websocket.Upgrader{ ReadBufferSize:  1024, WriteBufferSize: 1024,}
	server := Server{mux, oauthCfg,cookieStore, templates,upgrader, nil}
	server.routes()
	return &server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//https://neoteric.eu/blog/how-to-serve-static-files-with-golang/
	s.router.ServeHTTP(w, r)
}

func NewIDP(config *Options) (*oauth2.Config, error){

	switch config.Idp  {
		case GitHubProvider:
			return &oauth2.Config{
				ClientID:     config.ClientID,
				ClientSecret: config.ClientSecret,
				Endpoint: oauth2.Endpoint{
					AuthURL:  githubAuthorizeUrl,
					TokenURL: githubTokenUrl,
				},
				RedirectURL: redirectUrl,
				Scopes:      []string{"repo"},
			},nil
	}
	log.Fatal().Msgf("not implemented %v",config.Idp)
	return nil, nil
}

