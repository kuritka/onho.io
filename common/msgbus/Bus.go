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

	BusImpl struct{
		channel    *amqp.Channel
		connection *amqp.Connection
		exmgr      *exchangeManagerImpl
	}





	IMsgBusPublisher interface {
		PublishCommand(msg string)
		PublishEvent(msg string)
	}

	msgBusPublisherImpl struct {
		msgBusImpl *BusImpl
	}
)


func NewMsgBus(connectionString string) *BusImpl {
	utils.FailOnEmptyString(connectionString,"connection string cannot be empty")
	conn, ch := qutils.GetChannel(connectionString)
	utils.FailOnNil(conn,"connection")
	utils.FailOnNil(ch,"channel")
	exmgr := newExchangeManager(conn, ch)
	return &BusImpl{
		ch,
		conn,
		exmgr,
	}
}


// Register service and attach it to the bus
func (mb *BusImpl) Register(name string) (*msgBusListenerImpl, *msgBusPublisherImpl){
	utils.FailOnEmptyString(name, "name cannot be nil")
	guid,_ := getGuid()
	serviceDiscoveryName := name + "_"+ "discovery" +"_" + guid
	serviceEventName := name + "_"+ "event" +"_" + guid
	serviceCommandName := name + "_"+ "command" +"_" + guid
	mb.exmgr.
		createDiscoveryExchangeIfNotExists().
		sendDiscoveryRequest(amqp.Publishing{Body: []byte(serviceDiscoveryName) })

	mb.exmgr.createEventExchangeIfNotExists()

	mb.exmgr.createCommandExchangeIfNotExists()

	return  newMsgBusListener(name,mb, serviceDiscoveryName,serviceEventName,serviceCommandName),
			&msgBusPublisherImpl{mb,}
}


func (mb *BusImpl) Close(){
	mb.exmgr.close()
}





