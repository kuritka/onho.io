package msgbus

import (
	"fmt"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)

type (
	IMsgBusListener interface {
		AddCommandHandler(name string, f func()) *msgBusListenerImpl
		AddEventHandler(name string, f func()) *msgBusListenerImpl
		Listen()
	}

	msgBusListenerImpl struct {
		serviceName     string
		eventHandlers   map[string]func(deliveries <-chan amqp.Delivery)
		commandHandlers map[string]func(deliveries <-chan amqp.Delivery)
		discoveryQueue  string
		eventQueue      string
		commandQueue    string
		qm              *queueManagerImpl
	}
)

func newMsgBusListener(serviceName string, msgBusImpl *BusImpl, discoveryQueue string, serviceEvent string, serviceCommand string) *msgBusListenerImpl {
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	return &msgBusListenerImpl{
		serviceName,
		map[string]func(input <-chan amqp.Delivery){},
		map[string]func(input <-chan amqp.Delivery){},
		discoveryQueue,
		serviceEvent,
		serviceCommand,
		qm,
	}
}

func (l *msgBusListenerImpl) AddCommandHandler(name string, f func(input <-chan amqp.Delivery)) *msgBusListenerImpl {
	l.commandHandlers[name] = f
	return l
}

func (l *msgBusListenerImpl) AddEventHandler(name string, f func(input <-chan amqp.Delivery)) *msgBusListenerImpl {
	l.eventHandlers[name] = f
	return l
}

func (l *msgBusListenerImpl) Listen() {

	for eventHandler := range l.eventHandlers {
		eventChannel, err := l.qm.createQueueIfNotExists(l.eventQueue, true).bindToQueue("", serviceEventExchange).consumeFromChannel()
		utils.FailOnError(err, "event exchange")
		go l.eventHandlers[eventHandler](eventChannel)
	}

	for cmdHandler := range l.commandHandlers {
		commandChannel, err := l.qm.createQueueIfNotExists(l.commandQueue, true).bindToQueue(cmdHandler, serviceCommandExchange).consumeFromChannel()
		utils.FailOnError(err, "command exchange")
		go l.commandHandlers[cmdHandler](commandChannel)
	}

	discos, err := l.qm.createQueueIfNotExists(l.discoveryQueue, true).bindToQueue("", serviceDiscoveryExchange).consumeFromChannel()
	utils.FailOnError(err, "discovery exchange")

	go listenForDiscoverRequests(discos)

}

func listenForDiscoverRequests(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		fmt.Println(string(d.Body))
	}
}
