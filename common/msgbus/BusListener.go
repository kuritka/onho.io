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
		eventQueue      string
		commandQueue    string
		qm              *queueManagerImpl
		cmdEventAggreagtor *eventAggregator
		evntEventAggreagtor *eventAggregator
		registry         map[string]<-chan amqp.Delivery
	}
)

func newMsgBusListener(serviceName string, msgBusImpl *BusImpl, discoveryQueue string, serviceEvent string, serviceCommand string) *msgBusListenerImpl {
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	return &msgBusListenerImpl{
		serviceName,
		serviceEvent,
		serviceCommand,
		qm,
		newEventAggregator(),
		newEventAggregator(),
		make(map[string]<-chan amqp.Delivery),
	}
}

func (l *msgBusListenerImpl) AddCommandHandler(name string, f func(amqp.Delivery)) *msgBusListenerImpl {
	l.cmdEventAggreagtor.AddListener(name, f)
	return l
}

func (l *msgBusListenerImpl) AddEventHandler(name string, f func(amqp.Delivery)) *msgBusListenerImpl {
	l.evntEventAggreagtor.AddListener(name, f)
	return l
}


func (l *msgBusListenerImpl) Listen() {
	evnts, err := l.bindHandlersToQueue(l.eventQueue, l.evntEventAggreagtor, serviceEventExchange)
	utils.FailOnError(err, fmt.Sprintf("%s %s", l.eventQueue, exchange.string(serviceEventExchange)))
	go l.consume(evnts, l.evntEventAggreagtor)
	//
	//cmds, err := l.bindHandlersToQueue(l.commandQueue, l.cmdEventAggreagtor, serviceCommandExchange)
	//utils.FailOnError(err, fmt.Sprintf("%s %s", l.commandQueue, exchange.string(serviceCommandExchange)))
	//go l.consume(cmds, l.cmdEventAggreagtor)


}

func (l *msgBusListenerImpl) consume(input <-chan amqp.Delivery, aggregator *eventAggregator) {
	for value := range input {
		aggregator.Publish(value.RoutingKey, value)
	}
}

func (l *msgBusListenerImpl) bindHandlersToQueue(queueName string, aggregator *eventAggregator, exchange exchange) (<-chan amqp.Delivery, error){
	var q = l.qm.createQueueIfNotExists(queueName, true)
	for handler := range aggregator.listeners {
		q.bindToQueue(handler, exchange)
	}
	return  q.consumeFromChannel()
}

func (l *msgBusListenerImpl) listenForDiscoverRequests(discoveryChannel <-chan amqp.Delivery) {
	for d := range discoveryChannel {
		fmt.Println("hit")
		workerQueue := string(d.Body)
			if l.registry[workerQueue] == nil {
				fmt.Println(workerQueue)
				l.registry[workerQueue] = make(<-chan amqp.Delivery)
				fmt.Println("sending " + l.commandQueue + " " + exchange.string(serviceDiscoveryExchange))
				//l.qm.publishMessage(exchange.string(serviceDiscoveryExchange),
				//	"", amqp.Publishing{Body: []byte(l.commandQueue)})
				l.qm.channel.
					Publish(exchange.string(serviceDiscoveryExchange),
						"", false, false,
						amqp.Publishing{Body: []byte(l.commandQueue)})
			}
		}
		//if l.registry[workerQueue] == nil {
		//	msgs,err := l.qm.channel.Consume(
		//		workerQueue,
		//		"",
		//		true,
		//		false,
		//		false,
		//		false,
		//		nil,
		//	)
		//	utils.FailOnError(err,"hovno")
		//  l.registry[workerQueue] = msgs
		//go func() {
		//	for d := range msgs {
		//		fmt.Printf("Received a message: %s", d.Body)
		//	}
		//}()
		//if l.registry[workerQueue] == nil {
		//	l.registry[workerQueue] = commandChannel
		//	go l.ProcessMessages(commandChannel)
		//}
	//}
}

func (l *msgBusListenerImpl) ProcessMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		fmt.Println("COMMAND: "+string(msg.Body))
	}
}
