package msgbus

import (
	"github.com/kuritka/onho.io/common/qutils"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)

type (
	IMsgBus interface {
		Register(name string) (*msgBusListenerImpl, *msgBusPublisherImpl)
		Close()
	}

	BusImpl struct {
		channel    *amqp.Channel
		connection *amqp.Connection
		exmgr      *exchangeManagerImpl
		registry map[string]*queueManagerImpl
	}
)

func NewMsgBus(connectionString string) *BusImpl {
	utils.FailOnEmptyString(connectionString, "connection string cannot be empty")
	conn, ch := qutils.GetChannel(connectionString)
	utils.FailOnNil(conn, "connection")
	utils.FailOnNil(ch, "channel")
	exmgr := newExchangeManager(conn, ch)
	return &BusImpl{
		ch,
		conn,
		exmgr,
		make(map[string]*queueManagerImpl),
	}
}

// Register service and attach it to the bus
func (mb *BusImpl) Register(name string) (*msgBusListenerImpl, *msgBusPublisherImpl) {
	utils.FailOnEmptyString(name, "name cannot be nil")
	guid, _ := utils.Guid()
	queueDiscoveryName := name + "_" + "discovery" + "_" + guid

	mb.exmgr.createEventExchangeIfNotExists()

	discoveryRequest := newGobMessageProvider().EncodeDisco(DiscoveryRequest{CommandQueue: register, ServiceGuid: guid })
	discos, err := mb.exmgr.createDiscoveryExchangeIfNotExists().
		createQueueIfNotExists(queueDiscoveryName, true).
		bindToQueue("", serviceDiscoveryExchange).consumeFromChannel()
	utils.FailOnError(err, "discovery exchange")

	err = mb.channel.Publish(exchange.string(serviceDiscoveryExchange),
		"", false, false, discoveryRequest)
	utils.DisposeOnError(err, "Unable publish to "+exchange.string(serviceDiscoveryExchange), mb.Close)

	return newMsgBusListener(mb,  name,  discos, guid, mb.registry),
		newMessageBusPublisher(mb, name,guid, mb.registry)
}


func (mb *BusImpl) Close() {
	mb.exmgr.close()
}