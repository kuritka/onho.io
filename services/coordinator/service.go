package coordinator

import (
	"github.com/kuritka/onho.io/common/log"
)

type Coordinator struct {
	options Options
}

var logger   = log.Log


func NewService(options Options) *Coordinator {
	return &Coordinator {
		options:options,
	}
}


func (c *Coordinator) Run() error {

	NewQueueListener(c.options).ListenForNewSource()

	return nil
}


