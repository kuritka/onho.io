package sensor

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/common/msgbus"
	"github.com/kuritka/onho.io/common/utils"
)

var logger   = log.Log

type Sensor struct {
	options Options
}


func NewService(options Options) *Sensor {

	utils.FailOnEmptyString(options.Name,"empty name")

	utils.FailOnEmptyString(options.ConnectionString,"empty connection string")

	return &Sensor{
		options,
	}
}


func (sm *Sensor) Run() error {
	mb := msgbus.NewMsgBus(sm.options.ConnectionString)
	defer mb.Close()

	listener, _ := mb.Register("myService")

	listener.Listen()
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