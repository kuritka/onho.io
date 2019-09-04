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
		eventQueue      string
		commandQueue    string
		qm              *queueManagerImpl
		cmdEventAggreagtor *eventAggregator
		evntEventAggreagtor *eventAggregator
		registry         map[string]<-chan amqp.Delivery
		discos 			<-chan amqp.Delivery
		msgProvider *messageProvider
	}
)

func newMsgBusListener( msgBusImpl *BusImpl,  serviceEvent string, serviceCommand string, discos <-chan amqp.Delivery) *msgBusListenerImpl {
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	return &msgBusListenerImpl{
		serviceEvent,
		serviceCommand,
		qm,
		newEventAggregator(),
		newEventAggregator(),
		make(map[string]<-chan amqp.Delivery),
		discos,
		newGobMessageProvider(),
	}
}

func (l *msgBusListenerImpl) AddCommandHandler(command string, f func(Message)) *msgBusListenerImpl {
	l.cmdEventAggreagtor.AddListener(command, f)
	return l
}

func (l *msgBusListenerImpl) AddEventHandler(name string, f func(Message)) *msgBusListenerImpl {
	l.evntEventAggreagtor.AddListener(name, f)
	return l
}


func (l *msgBusListenerImpl) Listen() {
	events, err := l.bindHandlersToQueue(l.eventQueue, l.evntEventAggreagtor, serviceEventExchange)
	utils.FailOnError(err, fmt.Sprintf("%s %s", l.eventQueue, exchange.string(serviceEventExchange)))

	cmds, err := l.qm.channel.Consume(l.commandQueue, "", true, false, false, false, nil)
	utils.FailOnError(err, "consuming name queue "+l.commandQueue)

	go l.listenForEvents(events)

	go l.listenForCommands(cmds)

	//go l.listenForDiscoveryRequests(l.commandQueue, l.discos)

}

func (l *msgBusListenerImpl) listenForEvents(messages <-chan amqp.Delivery) {
	for value := range messages {
		m := Message{Name: value.RoutingKey,Message: string(value.Body)}
		l.evntEventAggreagtor.Publish(m)
	}
}

func (l *msgBusListenerImpl) listenForCommands(messages <-chan amqp.Delivery) {
	for msg := range messages{
		cmd := l.msgProvider.Decode(msg)
		l.cmdEventAggreagtor.Publish(cmd)
	}
}



func (l *msgBusListenerImpl) bindHandlersToQueue(queueName string, aggregator *eventAggregator, exchange exchange) (<-chan amqp.Delivery, error){
	var q = l.qm.createQueueIfNotExists(queueName, true)
	for handler := range aggregator.listeners {
		q.bindToQueue(handler, exchange)
	}
	return  q.consumeFromChannel()
}

func (l *msgBusListenerImpl) listenForDiscoveryRequests(queueCommandName string, discoveryChannel <-chan amqp.Delivery) {
	for msg := range discoveryChannel {
		workerQueue := string(msg.Body)
		if l.registry[workerQueue] == nil {
			l.registry[workerQueue] = make(<-chan amqp.Delivery)
			err := l.qm.channel.Publish(exchange.string(serviceDiscoveryExchange),
				"", false, false, amqp.Publishing{Body: []byte(queueCommandName )})
			utils.FailOnError(err, "Unable publish to "+exchange.string(serviceDiscoveryExchange))
			//for listener := range l.cmdEventAggreagtor.listeners {
			//	if listener ==
				cmds, err := l.qm.channel.Consume(workerQueue, "", true, false, false, false, nil)
				utils.FailOnError(err, "consuming name queue "+workerQueue)
				go l.ProcessCommand(cmds)
		//	}
		}
	}
}

func (l *msgBusListenerImpl) ProcessCommand(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		cmd := l.msgProvider.Decode(msg)
		l.cmdEventAggreagtor.Publish(cmd)
	}
}
