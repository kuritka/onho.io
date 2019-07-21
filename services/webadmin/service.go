package webadmin

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services"
	"github.com/kuritka/onho.io/services/webadmin/controller"
	"net/http"
	"strconv"
)

type WebAdmin struct {
	options    controller.Options
	aggregator services.IEventAggregator
}

var logger   = log.Log


func NewService(options controller.Options, aggregator services.IEventAggregator) *WebAdmin {
	utils.FailOnNil(aggregator,"event aggregator")
	return &WebAdmin {
		options:options,
		aggregator: aggregator,
	}
}


func (c *WebAdmin) Run() error {

	fmt.Printf("listening on port :"+strconv.Itoa(c.options.Port))
	controller.NewStartup(c.options).Init()

	http.ListenAndServe(":"+strconv.Itoa(c.options.Port), nil)

	return nil
}



