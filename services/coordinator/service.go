package coordinator

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/common/msgbus"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services"
)

type Coordinator struct {
	options Options
	aggregator services.IEventAggregator
}

var logger   = log.Log

func NewService(options Options, aggregator services.IEventAggregator) *Coordinator {
	utils.FailOnNil(aggregator,"event aggregator")
	return &Coordinator {
		options:options,
		aggregator: aggregator,
	}
}

func (c *Coordinator) Run() error {
	NewDatabaseConsumer(c.aggregator, c.options.QueueConnectionString)
	newWebAppConsumer(c.options,c.aggregator)//.ListenForDiscoveryRequest()

	msgBus :=  msgbus.NewMsgBus(c.options.QueueConnectionString)
	defer msgBus.Close()
	listener, _ :=  msgBus.Register("coordinator")

	listener.
		AddCommandHandler("cmd-tick",).
		Listen()

	var a string
	fmt.Println("listening")
	fmt.Scanln(&a)
	return nil
	//q:= NewQueueListener(c.options,c.aggregator)
	//go q.ListenForNewSource()
	//return nil
}


