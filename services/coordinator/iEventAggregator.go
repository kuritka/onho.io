package coordinator

type IEventAggregator interface {
	AddListener(name string, callback func(EventData))
	PublishEvent(name string, data EventData)
}

