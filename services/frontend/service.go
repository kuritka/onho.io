package frontend

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
	//utils.FailIfFileNotExist("./onho.crt","missing ./onho.crt")
	//utils.FailIfFileNotExist("./onho.crt","missing ./onho.key")

	msgBus :=  f.dependencies.MsgBus
	defer msgBus.Close()
	_, publisher :=  msgBus.Register("frontend")

	i := 3
	commandPublisher := func(data string){
		publisher.Command("recognised",data)
		fmt.Println( fmt.Sprintf("%v -> COMMAND SEND" , i))
		i++
	}

	server := NewServer(mux.NewRouter(), f.dependencies.CookieStore, f.dependencies.Auth, commandPublisher)
	listenAddr :=  fmt.Sprintf(":%v",f.dependencies.Port)
	log.Printf("listening HTTP on %s",listenAddr)
	return http.ListenAndServe(listenAddr, server)
	//return http.ListenAndServeTLS(listenAddr,"./onho.crt","./onho.key", server)
}

