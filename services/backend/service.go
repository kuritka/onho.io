package backend

import (
	"fmt"
	"github.com/kuritka/onho.io/common/manager/depresolver"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/kuritka/onho.io/common/utils"
)



type Frontend struct {
	dependencies depresolver.Dependencies
}

func NewService(dependencies depresolver.Dependencies)  *Frontend {
	utils.FailOnNil(dependencies,"nil dependencies")
	return & Frontend{
		 dependencies,
	}
}




func (f *Frontend) Run() error {
	utils.FailIfFileNotExist("./onho.crt","missing ./onho.crt")
	utils.FailIfFileNotExist("./onho.crt","missing ./onho.key")

	//TODO: ADD websockets server and charts etc..

	server := NewServer(mux.NewRouter())
	listenAddr :=  fmt.Sprintf(":%v",f.dependencies.Port)
	log.Printf("listening HTTP on %s",listenAddr)

	return http.ListenAndServe(listenAddr, server)
	//return http.ListenAndServeTLS(listenAddr,"./onho.crt","./onho.key", server)
}

