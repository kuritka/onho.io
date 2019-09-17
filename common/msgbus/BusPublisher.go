package msgbus

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/rs/zerolog/log"
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
	utils.FailOnEmptyString(guid, "missing guid")
	utils.FailOnNil(msgBusImpl, "MessageBusImpl")
	utils.FailOnNil(registry,"registry")
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	msgProvider := newGobMessageProvider()
	return &msgBusPublisherImpl{  guid, name,msgBusImpl, qm, msgProvider, registry}
}

func (p *msgBusPublisherImpl) Command(command string, data string) {
	cq := command + "_" + p.guid
	if p.registry[cq] == nil {
		disco := DiscoveryRequest{CommandQueue: cq, CommandHandlers: []string{command}, ServiceGuid: p.guid}
		p.registry[cq] = p.mb.exmgr.createQueueIfNotExists(cq, true)
		err := p.qm.channel.Publish(exchange.string(serviceDiscoveryExchange),
			"", false, false, p.msgProvider.EncodeDisco(disco))
		utils.DisposeOnError(err, "Unable publish to "+exchange.string(serviceDiscoveryExchange), p.mb.Close)
		log.Info().Msg("PUBLISHED: " +cq )
	}
	msg := p.msgProvider.Encode(Message{Name: command, Message:data })
	p.registry[cq].publishMessage(exchange.string(exchangeWorkerQueue), cq, msg)
}

func (p *msgBusPublisherImpl) Event(eventName string, data string) {
	p.qm.publishMessage(exchange.string(serviceEventExchange), eventName, amqp.Publishing{Body: []byte(data)})
}
