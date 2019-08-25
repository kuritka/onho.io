package sensor

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/common/msgbus"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
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
		AddEventHandler("foo", func(input <-chan amqp.Delivery) {
			for i := range input {
				fmt.Printf("foo event: %s \n", i.Body)
			}
		}).
		AddEventHandler("blah", func(input <-chan amqp.Delivery) {
			for i := range input {
				fmt.Printf("blah event: %s \n", i.Body)
			}
		}).
		AddCommandHandler("cmdA", func(input <-chan amqp.Delivery) {
			for i := range input {
				fmt.Printf("cmdA command: %s \n", i.Body)
			}
		}).
		AddCommandHandler("cmdB", func(input <-chan amqp.Delivery) {
			for i := range input {
				fmt.Printf("cmdB command: %s \n", i.Body)
			}
		}).
		Listen()

	publisher.Command("cmdA", "blaah")
	publisher.Command("cmdA", "blaah2")
	publisher.Command("cmdB", "blaah3B")
	publisher.Command("cmdX", "cmdX")
	publisher.Command("cmdA", "blaah4")

	//conn, ch := qutils.GetChannel(sm.options.ConnectionString)
	//var provider = bus.NewProvider(conn, ch)
	//provider.Register(sm.options.Name)
	//defer provider.Close()
	//provider.Listen()

	var a string
	fmt.Println("listening")
	fmt.Scanln(&a)

	return nil
}
