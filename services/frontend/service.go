package frontend

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"

	"github.com/kuritka/onho.io/common/utils"
)

func BuildContainer(options Options) *dig.Container {
	container := dig.New()
	err := container.Provide(mux.NewRouter)
	utils.FailFastOnError(err)
	err = container.Provide(NewIDP)
	utils.FailFastOnError(err)
	err = container.Provide(NewServer)
	utils.FailFastOnError(err)
	err = container.Provide(options)
	utils.FailFastOnError(err)
	err = container.Provide(func() *Options{ return &options})
	utils.FailFastOnError(err)
	return container
}





type Frontend struct {
	options Options
}

func NewService(options Options)  *Frontend {
	utils.FailOnEmptyString(options.ClientID, "missing clientID")
	utils.FailOnEmptyString(options.ClientSecret, "missing clientSecret")
	utils.FailOnEmptyString(options.CookieStoreKey, "missing cookieStoreKey")
	utils.FailOnLessOrEqualToZero(options.Port, "invalid port number")

	return & Frontend{
		options ,

	}
}


func (f *Frontend) Run() error {
	oauth,_ := NewIDP(&f.options)
	server := NewServer(mux.NewRouter(), &f.options, oauth)
	listenAddr :=  fmt.Sprintf(":%v",f.options.Port)
	log.Printf("listening on %s",listenAddr)
	return http.ListenAndServe(listenAddr, server)
}
