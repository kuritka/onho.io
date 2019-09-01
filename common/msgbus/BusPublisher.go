package msgbus

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)

type (
	IMsgBusPublisher interface {
		//Command can be send to multiple targets but only one can process
		Command(name string, message string)
		//Command can be send to multiple targets but all of them can process
		Event(name string, message string)
	}

	msgBusPublisherImpl struct {
		serviceName string
		mb          *BusImpl
		qm          *queueManagerImpl
	}
)

func newMessageBusPublisher(serviceName string, msgBusImpl *BusImpl) *msgBusPublisherImpl {
	utils.FailOnEmptyString(serviceName, "service name")
	utils.FailOnNil(msgBusImpl, "MessageBusImpl")
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	return &msgBusPublisherImpl{serviceName, msgBusImpl, qm}
}

func (p *msgBusPublisherImpl) Command(cmdName string, data string) {
	//p.qm.publishMessage(exchange.string(serviceCommandExchange), cmdName, amqp.Publishing{Body: []byte(data)})
}

func (p *msgBusPublisherImpl) Event(eventName string, data string) {
	p.qm.publishMessage(exchange.string(serviceEventExchange), eventName, amqp.Publishing{Body: []byte(data)})
}
