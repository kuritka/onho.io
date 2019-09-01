package msgbus

import (
	"fmt"
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
		registry         map[string]<-chan amqp.Delivery
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
		make(map[string]<-chan amqp.Delivery),
	}
}

// Register service and attach it to the bus
func (mb *BusImpl) Register(name string) (*msgBusListenerImpl, *msgBusPublisherImpl) {
	utils.FailOnEmptyString(name, "name cannot be nil")
	guid, _ := getGuid()
	queueDiscoveryName := name + "_" + "discovery" + "_" + guid
	queueEventName := name + "_" + "event" + "_" + guid
	queueCommandName := name + "_" + "command" + "_" + guid

	mb.exmgr.createQueueIfNotExists(queueCommandName,true)

	mb.exmgr.createEventExchangeIfNotExists()

	//Queuebinding for discos must complete before first request start, otherwise there will be still one service without
	//knowledge of other services
	discos, err := mb.exmgr.createDiscoveryExchangeIfNotExists().
		createQueueIfNotExists(queueDiscoveryName, true).
		bindToQueue("", serviceDiscoveryExchange).consumeFromChannel()
	utils.FailOnError(err, "discovery exchange")
	mb.exmgr.sendDiscoveryRequest(amqp.Publishing{Body: []byte(queueCommandName)})
	go mb.listenForDiscoveryRequests(queueCommandName, discos)

	return newMsgBusListener(name, mb, queueDiscoveryName, queueEventName, queueCommandName),
		newMessageBusPublisher(name, mb)
}

func (mb *BusImpl) listenForDiscoveryRequests(queueCommandName string, discoveryChannel <-chan amqp.Delivery) {
	for msg := range discoveryChannel {
		workerQueue := string(msg.Body)
		if mb.registry[workerQueue] == nil {
			fmt.Println(workerQueue)
			mb.registry[workerQueue] = make(<-chan amqp.Delivery)
			fmt.Println("sending " + queueCommandName + " " + exchange.string(serviceDiscoveryExchange))
			mb.exmgr.channel.
				Publish(exchange.string(serviceDiscoveryExchange),
					"", false, false,
					amqp.Publishing{Body: []byte(queueCommandName )})
			fmt.Println("\n\nREGISTRY:")
			for x := range mb.registry {
				fmt.Println(x)
			}
		}
	}
}


func (mb *BusImpl) Close() {
	mb.exmgr.close()
}