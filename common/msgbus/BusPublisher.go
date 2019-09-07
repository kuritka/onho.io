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
		commandQueueName string
		mb          *BusImpl
		qm          *queueManagerImpl
		msgProvider *messageProvider
		registry map[string]string
	}
)

func newMessageBusPublisher(commandQueueName string, msgBusImpl *BusImpl, registry map[string]string) *msgBusPublisherImpl {
	utils.FailOnEmptyString(commandQueueName, "service name")
	utils.FailOnNil(msgBusImpl, "MessageBusImpl")
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	msgProvider := newGobMessageProvider()
	return &msgBusPublisherImpl{ commandQueueName, msgBusImpl, qm, msgProvider, registry}
}

func (p *msgBusPublisherImpl) Command(targetService string, command string, data string) {
	queue := p.registry[command]
	msg := p.msgProvider.Encode(Message{Name: command, Message:data })
	p.qm.publishMessage(exchange.string(exchangeWorkerQueue), queue, msg)
}

func (p *msgBusPublisherImpl) Event(eventName string, data string) {
	p.qm.publishMessage(exchange.string(serviceEventExchange), eventName, amqp.Publishing{Body: []byte(data)})
}
