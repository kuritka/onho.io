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
		guid		string
		name        string
		mb          *BusImpl
		qm          *queueManagerImpl
		msgProvider *messageProvider
		registry   map[string]*queueManagerImpl
	}
)

func newMessageBusPublisher( msgBusImpl *BusImpl,  name string, guid string, registry map[string]*queueManagerImpl) *msgBusPublisherImpl {
	utils.FailOnNil(msgBusImpl, "MessageBusImpl")
	utils.DisposeOnEmptyString(name, "missing name", msgBusImpl.Close)
	utils.DisposeOnEmptyString(guid, "missing guid", msgBusImpl.Close)
	utils.DisposeOnNil(registry,"registry", msgBusImpl.Close)
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	msgProvider := newGobMessageProvider()
	return &msgBusPublisherImpl{  guid, name,msgBusImpl, qm, msgProvider, registry}
}

func (p *msgBusPublisherImpl) Command(command string, data string) {
	utils.DisposeOnEmptyString(command, "command" , p.mb.Close)
	utils.DisposeOnEmptyString(data, "command data" , p.mb.Close)
	cq := command + "_" + p.guid
	if p.registry[cq] == nil {
		disco := DiscoveryRequest{CommandQueue: cq, ServiceGuid: p.guid}
		p.registry[cq] = p.mb.exmgr.createQueueIfNotExists(cq, true)
		err := p.qm.channel.Publish(exchange.string(serviceDiscoveryExchange),
			"", false, false, p.msgProvider.EncodeDisco(disco))
		utils.DisposeOnError(err, "Unable publish to "+exchange.string(serviceDiscoveryExchange), p.mb.Close)
	}
	msg := p.msgProvider.Encode(Message{Name: command, Message:data })
	p.registry[cq].publishMessage(exchange.string(exchangeWorkerQueue), cq, msg)
}

func (p *msgBusPublisherImpl) Event(event string, data string) {
	utils.DisposeOnEmptyString(event, "evet" , p.mb.Close)
	utils.DisposeOnEmptyString(data, "event data" , p.mb.Close)
	p.qm.publishMessage(exchange.string(serviceEventExchange), event, amqp.Publishing{Body: []byte(data)})
}
