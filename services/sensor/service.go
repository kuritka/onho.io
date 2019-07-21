package sensor

import (
	"github.com/kuritka/onho.io/common/bus"
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/common/qutils"
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
	conn, ch := qutils.GetChannel(sm.options.ConnectionString)
	var provider = bus.NewProvider(conn, ch)
	provider.Register(sm.options.Name)
	defer provider.Close()
	return nil
}