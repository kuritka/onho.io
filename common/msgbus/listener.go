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

const (
	register = "register"
)

func newMsgBusListener( msgBusImpl *BusImpl,  name string,  discos <-chan amqp.Delivery, guid string, registry map[string]*queueManagerImpl) *msgBusListenerImpl {
	utils.FailOnNil(msgBusImpl, "MessageBusImpl")
	utils.DisposeOnEmptyString(name, "missing name", msgBusImpl.Close)
	utils.DisposeOnEmptyString(guid, "missing guid", msgBusImpl.Close)
	utils.DisposeOnNil(discos,"discovery channel", msgBusImpl.Close)
	utils.DisposeOnNil(registry,"registry", msgBusImpl.Close)
	qm := createQueueManager(msgBusImpl.connection, msgBusImpl.channel)
	return &msgBusListenerImpl{
		name + "_" + "event" + "_" + guid,
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
	utils.DisposeOnEmptyString(command,"empty command",l.qm.close )
	utils.DisposeOnNil(f,"empty command",l.qm.close )
	l.cmdEventAggreagtor.AddListener(command, f)
	return l
}

func (l *msgBusListenerImpl) AddEventHandler(event string, f func(Message)) *msgBusListenerImpl {
	utils.DisposeOnEmptyString(event,"empty command",l.qm.close )
	utils.DisposeOnNil(f,"empty command",l.qm.close )
	l.evntEventAggreagtor.AddListener(event, f)
	return l
}

func (l *msgBusListenerImpl) Listen() {
	events, err := l.bindHandlersToQueue(l.eventQueue, l.evntEventAggreagtor, serviceEventExchange)
	utils.DisposeOnError(err, fmt.Sprintf("%s %s", l.eventQueue, exchange.string(serviceEventExchange)), l.qm.close)

	go l.listenForEvents(events)

	go l.listenForDiscoveryRequests(l.discos)
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

		if discoMessage.ServiceGuid == l.guid {
			continue
		}

		if discoMessage.CommandQueue == register {
			l.publishCommandRegistry()
			continue
		}

		if l.registry[discoMessage.CommandQueue] == nil {
			stream , err := l.qm.channel.Consume(discoMessage.CommandQueue,"", true, false, false, false, nil)
			l.registry[discoMessage.CommandQueue] = stream
			utils.DisposeOnError(err, "cannot consume from " + discoMessage.CommandQueue, l.qm.close)
			go l.listenForCommands(stream)
		}
	}
}

func (l *msgBusListenerImpl) publishCommandRegistry(){
	for cq := range l.publishedCommands {
		discoPublishing :=  l.msgProvider.EncodeDisco(DiscoveryRequest{CommandQueue: cq,  ServiceGuid: l.guid })
		err := l.qm.channel.Publish(exchange.string(serviceDiscoveryExchange),
			"", false, false, discoPublishing)
		utils.DisposeOnError(err, "Unable publish to "+exchange.string(serviceDiscoveryExchange), l.qm.close)
	}
}
