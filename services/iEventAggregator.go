package services

type IEventAggregator interface {
	AddListener(name string, callback func(interface{}))
	PublishEvent(name string, data interface{})
}
