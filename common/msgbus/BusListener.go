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
		eventQueue          string
		guid                string
		qm                  *queueManagerImpl
		cmdEventAggreagtor  *eventAggregator
		evntEventAggreagtor *eventAggregator
		discos              <-chan amqp.Delivery
		msgProvider         *messageProvider
		registry      		map[string]<-chan amqp.Delivery
		publishedCommands   map[string]*queueManagerImpl
	}
)

func newMsgBusListener( msgBusImpl *BusImpl,  serviceEvent string,  discos <-chan amqp.Delivery, guid string, registry map[string]*queueManagerImpl) *msgBusListenerImpl {

	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	return &msgBusListenerImpl{
		serviceEvent,
		guid,
		qm,
		newEventAggregator(),
		newEventAggregator(),
		discos,
		newGobMessageProvider(),
		make(map[string]<-chan amqp.Delivery),
		registry,
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

	discoPublishing :=  l.msgProvider.EncodeDisco(DiscoveryRequest{CommandQueue: "register", CommandHandlers: []string{}, ServiceGuid: l.guid })
	l.sendDiscoveryRequest(discoPublishing)

	go l.listenForEvents(events)

	go l.listenForDiscoveryRequests(l.discos)
}

func (l *msgBusListenerImpl) sendDiscoveryRequest(discoPublishing amqp.Publishing) {
	l.qm.publishMessage( exchange.string(serviceDiscoveryExchange), "",discoPublishing )
}

func (l *msgBusListenerImpl) listenForEvents(messages <-chan amqp.Delivery) {
	for value := range messages {
		m := Message{Name: value.RoutingKey,Message: string(value.Body)}
		l.evntEventAggreagtor.Publish(m)
	}
}

func (l *msgBusListenerImpl) listenForCommands(messages <-chan amqp.Delivery) {
	for msg := range messages{
		cmd := l.msgProvider.DecodeMessage(msg)
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

func (l *msgBusListenerImpl) listenForDiscoveryRequests(discoveryChannel <-chan amqp.Delivery) {
	for msg := range discoveryChannel {
		discoMessage := l.msgProvider.DecodeDisco(msg)

		if discoMessage.CommandQueue == "register" {
			//fmt.Println("REGISTER call " + discoMessage.CommandQueue)
			//l.publishCommandRegistry()
			continue
		}

		if l.registry[discoMessage.CommandQueue] == nil {
			fmt.Println("REGISTERING: "+ discoMessage.CommandQueue)
			stream , err := l.qm.channel.Consume(discoMessage.CommandQueue,"", true, false, false, false, nil)
			l.registry[discoMessage.CommandQueue] = stream
			utils.DisposeOnError(err, "cannot consume from " + discoMessage.CommandQueue, l.qm.close)
			go l.listenForCommands(stream)
		}
	}
}


func (l *msgBusListenerImpl) publishCommandRegistry(){
	commands := []string{}
	for command,_ := range l.publishedCommands {
		discoPublishing :=  l.msgProvider.EncodeDisco(DiscoveryRequest{CommandQueue: command+"_"+l.guid, CommandHandlers: commands, ServiceGuid: l.guid })
		l.sendDiscoveryRequest(discoPublishing)
	}
}
