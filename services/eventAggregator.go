package services

import (
	"github.com/kuritka/onho.io/common/utils"
)

const (
		DataSourceDiscovered = "DataSourceDiscovered"
		MessageReceivedPrefix = "MessageReceived_"
)


type EventAggregator struct {
	listeners map[string][]func(interface{})
}

func NewEventAggregator() *EventAggregator {
	return &EventAggregator{
		listeners: make(map[string][]func(interface{})),
	}
}

func (ea *EventAggregator) AddListener(name string, callback func(interface{})) {
	utils.FailOnNil(callback, "callback is not set")
	utils.FailOnEmptyString(name, "name is empty")
	ea.listeners[name] = append(ea.listeners[name], callback)
}


func (ea *EventAggregator) PublishEvent(name string, data interface{} ){
	utils.FailOnEmptyString(name, "name is empty")
	if ea.listeners[name] != nil {
		for _, r := range ea.listeners[name] {
			r(data)
		}
	}
}
