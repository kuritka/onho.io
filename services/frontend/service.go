package frontend

import (
	"fmt"
	"github.com/kuritka/onho.io/common/msgbus"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/kuritka/onho.io/common/utils"
)



type Frontend struct {
	options Options
}

	func NewService(options Options)  *Frontend {
	utils.FailOnEmptyString(options.ClientID, "missing clientID")
	utils.FailOnEmptyString(options.ClientSecret, "missing clientSecret")
	utils.FailOnEmptyString(options.CookieStoreKey, "missing cookieStoreKey")
	utils.FailOnEmptyString(options.QueueConnectionString, "queue connection string")
	utils.FailOnLessOrEqualToZero(options.Port, "invalid port number")

	return & Frontend{
		options ,

	}
}


func (f *Frontend) Run() error {

	msgBus :=  msgbus.NewMsgBus(f.options.QueueConnectionString)
	defer msgBus.Close()
	_, publisher :=  msgBus.Register("frontend")

	i := 3
	commandPublisher := func(data string){
		publisher.Command("recognised",data)
		fmt.Println( fmt.Sprintf("%v -> COMMAND SEND" , i))
		i++
	}

	oauth,_ := NewIDP(&f.options)

	server := NewServer(mux.NewRouter(), &f.options, oauth, commandPublisher)
	listenAddr :=  fmt.Sprintf(":%v",f.options.Port)
	log.Printf("listening on %s",listenAddr)
	return http.ListenAndServe(listenAddr, server)
}
