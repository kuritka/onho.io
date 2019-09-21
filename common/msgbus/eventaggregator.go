package msgbus

import (
	"github.com/kuritka/onho.io/common/utils"
)

type eventAggregator struct {
	listeners map[string][]func(message Message)
}

func newEventAggregator() *eventAggregator {
	return &eventAggregator{
		listeners: make(map[string][]func(message Message)),
	}
}

func (ea *eventAggregator) AddListener(key string, callback func(Message)) {
	utils.FailOnNil(callback, "callback is not set")
	utils.FailOnEmptyString(key, "key is empty")
	ea.listeners[key] = append(ea.listeners[key], callback)
}

func (ea *eventAggregator) Publish(message Message){
	if ea.listeners[message.Name] != nil {
		for _, f := range ea.listeners[message.Name] {
			f(message)
		}
	}
}