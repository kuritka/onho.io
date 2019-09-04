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
	}
}

// Register service and attach it to the bus
func (mb *BusImpl) Register(name string) (*msgBusListenerImpl, *msgBusPublisherImpl) {
	utils.FailOnEmptyString(name, "name cannot be nil")
	guid, _ := getGuid()
	queueDiscoveryName := name + "_" + "discovery" + "_" + guid
	queueEventName := name + "_" + "event" + "_" + guid
	queueCommandName := name + "_" + "command"

	mb.exmgr.createQueueIfNotExists(queueCommandName,true)

	mb.exmgr.createEventExchangeIfNotExists()

	//Queuebinding for discos must complete before first request start, otherwise there will be still one service without
	//knowledge of other services
	discos, err := mb.exmgr.createDiscoveryExchangeIfNotExists().
		createQueueIfNotExists(queueDiscoveryName, true).
		bindToQueue("", serviceDiscoveryExchange).consumeFromChannel()
	utils.FailOnError(err, "discovery exchange")
	mb.exmgr.sendDiscoveryRequest(amqp.Publishing{Body: []byte(queueCommandName)})

	return newMsgBusListener(mb,  queueEventName, queueCommandName, discos),
		newMessageBusPublisher(queueCommandName, mb)
}

func (mb *BusImpl) Close() {
	mb.exmgr.close()
}