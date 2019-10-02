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
		name        string
		guid 		string
		mb          *BusImpl
		qm          *queueManagerImpl
		msgProvider *messageProvider
		registry   map[string]*queueManagerImpl
	}
)

func newMessageBusPublisher( msgBusImpl *BusImpl,  name string, guid string, registry map[string]*queueManagerImpl) *msgBusPublisherImpl {
	utils.FailOnNil(msgBusImpl, "MessageBusImpl")
	utils.DisposeOnEmptyString(name, "missing name", msgBusImpl.Close)
	utils.DisposeOnEmptyString(guid, "guid", msgBusImpl.Close)
	utils.DisposeOnNil(registry,"registry", msgBusImpl.Close)
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	msgProvider := newGobMessageProvider()
	return &msgBusPublisherImpl{  name,guid, msgBusImpl, qm, msgProvider, registry}
}

func (p *msgBusPublisherImpl) Command(command string, data string) {
	utils.DisposeOnEmptyString(command, "command" , p.mb.Close)
	utils.DisposeOnEmptyString(data, "command data" , p.mb.Close)
	if p.registry[command] == nil {
		disco := DiscoveryRequest{CommandQueue: command, ServiceGuid: p.guid}
		p.registry[command] = p.mb.exmgr.createQueueIfNotExists(command, false)
		err := p.qm.channel.Publish(exchange.string(serviceDiscoveryExchange),
			"", false, false, p.msgProvider.EncodeDisco(disco))
		utils.DisposeOnError(err, "Unable publish to "+exchange.string(serviceDiscoveryExchange), p.mb.Close)
	}
	//if queue was closed by consumer than remove. If no one enqueued messages, keep queue in alive
	msg := p.msgProvider.Encode(Message{Name: command, Message:data })
	p.registry[command].publishMessage(exchange.string(exchangeWorkerQueue), command, msg)
}

func (p *msgBusPublisherImpl) Event(event string, data string) {
	utils.DisposeOnEmptyString(event, "event" , p.mb.Close)
	utils.DisposeOnEmptyString(data, "event data" , p.mb.Close)
	p.qm.publishMessage(exchange.string(serviceEventExchange), event, amqp.Publishing{Body: []byte(data)})
}
