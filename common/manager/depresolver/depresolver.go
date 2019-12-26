package depresolver

import (
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/kuritka/onho.io/common/msgbus"
	"github.com/kuritka/onho.io/common/utils"
	"golang.org/x/oauth2"
	"sync"
)

type IDependencyResolver interface {
	// returns db connected to app
	MustResolveDatabase() *sql.DB
	// returns rabbit-mq
	MustResolveRabbitMQ() msgbus.IMsgBus
	//MustResolveGrpc() *grpcpool.Pool
	MustResolveEnvironment() string

	MustResolveGithubOAuth() *oauth2.Config
}


type FromEnvDependenciesResolver struct {
	database struct {
		initPlatformOnce sync.Once

	}

	queues struct {
		initPlatformMQOnce sync.Once
	}

	oauth struct {
		initPlatformMQOnce sync.Once
	}

	Options Dependencies
}

type Dependencies struct {
	Environment string
	Port int
	Db *sql.DB
	MsgBus msgbus.IMsgBus
	Auth *oauth2.Config
	CookieStore *sessions.CookieStore
}

func NewFromEnvDependencyResolver() *FromEnvDependenciesResolver{
	return &FromEnvDependenciesResolver{}
}


func (r *FromEnvDependenciesResolver) MustResolveDatabase() *FromEnvDependenciesResolver{
	utils.NotImplemented()
	return nil
}

func (r *FromEnvDependenciesResolver) MustResolveRabbitMQ() *FromEnvDependenciesResolver{
	r.queues.initPlatformMQOnce.Do(func (){
		connectionString := utils.MustGetStringFlagFromEnv("ONHO_RABBIT_CONNECTION_STRING")
		r.Options.MsgBus = msgbus.NewMsgBus(connectionString)
	})
	return r
}



func (r *FromEnvDependenciesResolver) MustResolveEnvironment() *FromEnvDependenciesResolver {
	r.Options.Environment = utils.MustGetStringFlagFromEnv("ONHO_ENVIRONMENT")
	return r
}

func (r *FromEnvDependenciesResolver) MustResolvePort() *FromEnvDependenciesResolver {
	r.Options.Port = utils.MustGetIntFlagFromEnv("ONHO_CLIENT_PORT")
	return r
}


func (r *FromEnvDependenciesResolver) MustResolveCookieStore() *FromEnvDependenciesResolver {
	cookieStoreKey := utils.MustGetStringFlagFromEnv("ONHO_OAUTH_COOKIE_KEY")
	r.Options.CookieStore = sessions.NewCookieStore([]byte(cookieStoreKey))
	return r
}



const (
	githubAuthorizeUrl = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
	redirectUrl        = ""
)

func (r *FromEnvDependenciesResolver) MustResolveGithubOAuth() *FromEnvDependenciesResolver {
	r.oauth.initPlatformMQOnce.Do(func() {
		r.Options.Auth = &oauth2.Config{
			ClientID:     utils.MustGetStringFlagFromEnv("ONHO_OAUTH_CLIENTID"),
			ClientSecret: utils.MustGetStringFlagFromEnv("ONHO_OAUTH_CLIENT_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  githubAuthorizeUrl,
				TokenURL: githubTokenUrl,
			},
			RedirectURL: redirectUrl,
			Scopes:      []string{"repo"},
		}
	})
	return r
}