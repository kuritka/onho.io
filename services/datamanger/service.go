package datamanager

import (
	"github.com/kuritka/onho.io/common/log"
)

type DataManager struct {
	options Options
}

var logger   = log.Log


func NewService(options Options ) *DataManager {
	return &DataManager {
		options:options,
	}
}


func (c *DataManager) Run() error {

	NewDbController(c.options).ReadPersistentQueue()

	return nil
}


