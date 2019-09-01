package msgbus

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)

type eventAggregator struct {
	listeners map[string][]func(amqp.Delivery)
}

func newEventAggregator() *eventAggregator {
	return &eventAggregator{
		listeners: make(map[string][]func(amqp.Delivery)),
	}
}

func (ea *eventAggregator) AddListener(key string, callback func(amqp.Delivery)) {
	utils.FailOnNil(callback, "callback is not set")
	utils.FailOnEmptyString(key, "key is empty")
	ea.listeners[key] = append(ea.listeners[key], callback)
}

func (ea *eventAggregator) Publish(name string, data amqp.Delivery ){
	utils.FailOnEmptyString(name, "name is empty")
	if ea.listeners[name] != nil {
		for _, f := range ea.listeners[name] {
			f(data)
		}
	}
}