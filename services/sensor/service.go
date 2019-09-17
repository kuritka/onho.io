package sensor

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/common/msgbus"
	"github.com/kuritka/onho.io/common/utils"
)

var logger = log.Log

type Sensor struct {
	options Options
}

func NewService(options Options) *Sensor {

	utils.FailOnEmptyString(options.Name, "empty name")

	utils.FailOnEmptyString(options.ConnectionString, "empty connection string")

	return &Sensor{
		options,
	}
}

func (sm *Sensor) Run() error {
	mb := msgbus.NewMsgBus(sm.options.ConnectionString)
	defer mb.Close()

	listener, publisher := mb.Register("myService")

	listener.
		AddEventHandler("foo", func(message msgbus.Message) {	fmt.Printf("foo event: %s \n", message.Message)}).
		AddEventHandler("blah", func(message msgbus.Message) {fmt.Printf("blah event: %s \n", message.Message)}).
		AddCommandHandler("cmdA", func(message msgbus.Message) {fmt.Printf("cmdA command: %s \n", message.Message)}).
		AddCommandHandler("cmdB", func(message msgbus.Message) {fmt.Printf("cmdB command: %s \n", message.Message)}).
		Listen()


	publisher.Command( "cmdA", "A1")
	publisher.Command( "cmdB", "B1")
	publisher.Command( "cmdB", "B2")
	publisher.Command( "cmdB", "B3")
	//
	//publisher.Event("foo", "foo: FOO EVENT")
	//publisher.Event("foo", "foo: FOO EVENT2")
	//publisher.Event("foo", "foo: FOO EVENT3")
	//publisher.Event("blah", "foo: BLAH EVENT1")
	//publisher.Event("foo", "foo: FOO EVENT4")

	var a string
	fmt.Println("listening")
	fmt.Scanln(&a)

	return nil
}
