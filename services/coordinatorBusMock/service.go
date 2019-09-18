package coordinatorBusMock

import (
	"fmt"
	"github.com/kuritka/onho.io/common/msgbus"
	"github.com/rs/zerolog/log"
)

var logger   = log.Log



type CoordinatorBusMock struct {
	options Options
}

func NewService(options Options) *CoordinatorBusMock {
	return &CoordinatorBusMock {
		options:options,
	}
}

func (c *CoordinatorBusMock) Run() error {

	mb := msgbus.NewMsgBus(c.options.QueueConnectionString)
	defer mb.Close()

	receiver, _ := mb.Register("serviceBus")

	receiver.
		AddCommandHandler("CMD-A",func(message msgbus.Message) {fmt.Printf("cmdA command: %s \n", message.Message)} ).
		Listen()

	var a string
	fmt.Println("listening")
	fmt.Scanln(&a)

	return nil
}

