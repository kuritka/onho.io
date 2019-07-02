package coordinator

import (
	"github.com/kuritka/onho.io/common/utils"
)


type EventAggregator struct {
	listeners map[string][]func(EventData)
}

func NewEventAggregator() *EventAggregator {
	return &EventAggregator{
		listeners: make(map[string][]func(EventData)),
	}
}

func (ea *EventAggregator) AddListener(name string, callback func(EventData)) {
	utils.FailOnNil(callback, "callback is not set")
	utils.FailOnEmptyString(name, "name is empty")
	ea.listeners[name] = append(ea.listeners[name], callback)
}


func (ea *EventAggregator) PublishEvent(name string, data EventData ){
	utils.FailOnEmptyString(name, "name is empty")
	if ea.listeners[name] != nil {
		for _, r := range ea.listeners[name] {
			r(data)
		}
	}
}
