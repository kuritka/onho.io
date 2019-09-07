package msgbus

import (
	"fmt"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
	"time"
)

type (
	IMsgBusListener interface {
		AddCommandHandler(name string, f func()) *msgBusListenerImpl
		AddEventHandler(name string, f func()) *msgBusListenerImpl
		Listen()
	}


	msgBusListenerImpl struct {
		eventQueue          string
		commandQueue        string
		guid                string
		qm                  *queueManagerImpl
		cmdEventAggreagtor  *eventAggregator
		evntEventAggreagtor *eventAggregator
		discos              <-chan amqp.Delivery
		msgProvider         *messageProvider
		commandRegistry     map[string]string
		serviceRegistry     map[string]bool
	}
)

func newMsgBusListener( msgBusImpl *BusImpl,  serviceEvent string, serviceCommand string, discos <-chan amqp.Delivery, registry map[string]string, guid string) *msgBusListenerImpl {
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	return &msgBusListenerImpl{
		serviceEvent,
		serviceCommand,
		guid,
		qm,
		newEventAggregator(),
		newEventAggregator(),
		discos,
		newGobMessageProvider(),
		registry,
		make(map[string]bool),
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

	discoPublishing := l.prepareDiscoveryRequest()
	l.sendDiscoveryRequest(discoPublishing)

	go l.listenForEvents(events)

	go l.listenForCommands(cmds)

	go l.listenForDiscoveryRequests(l.commandQueue, l.discos)
	//whole command bus must be changed. This delay is here because we want to
	// listeners wait for discos.Than commands can be sent
	time.Sleep(2 * time.Second)
}

func (l *msgBusListenerImpl) sendDiscoveryRequest(discoPublishing amqp.Publishing) {
	err := l.qm.channel.Publish(exchange.string(serviceDiscoveryExchange),
		"", false, false, discoPublishing)
	utils.FailOnError(err, "Unable publish to "+exchange.string(serviceDiscoveryExchange))
}

func (l *msgBusListenerImpl) prepareDiscoveryRequest() amqp.Publishing{
	var arr []string
	for k := range l.cmdEventAggreagtor.listeners {
		arr = append(arr, k)
	}
	disco := DiscoveryRequest{CommandQueue: l.commandQueue, CommandHandlers: arr, ServiceGuid: l.guid }
	return l.msgProvider.EncodeDisco(disco)
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

func (l *msgBusListenerImpl) listenForDiscoveryRequests(queueCommandName string, discoveryChannel <-chan amqp.Delivery) {
	for msg := range discoveryChannel {
		discoMessage := l.msgProvider.DecodeDisco(msg)
		if !l.serviceRegistry[discoMessage.ServiceGuid] {
			l.serviceRegistry[discoMessage.ServiceGuid] = true
			for _,cmd := range discoMessage.CommandHandlers {
				if q, found := l.commandRegistry[cmd]; found {
					if q != discoMessage.CommandQueue {
						utils.Fail(fmt.Sprintf("command %s is not unique", cmd ))
					}
					continue
				}
				l.commandRegistry[cmd] = discoMessage.CommandQueue
			}
			discoPublishing := l.prepareDiscoveryRequest()
			l.sendDiscoveryRequest(discoPublishing)
		}
	}
}
