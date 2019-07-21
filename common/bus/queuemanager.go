package bus

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)


type (
	iQueueManager interface {
		createExchangeIfNotExists(exchange exchange, k exchangeType, autoDelete bool) *queueProviderImpl
		close()
	}

	queueManagerImpl struct {
		connection *amqp.Connection
		channel *amqp.Channel
	}
)

type (
	iQueueProvider interface {
		createQueue(name string, autoDelete bool) *messageProviderImpl
	}
	queueProviderImpl struct {
		qm 	*queueManagerImpl
		exchange     exchange
		exchangeType exchangeType
	}
)

type (
	iExchangeProvider interface {
		sendDiscoveryEvent(serviceName string)
		bindToQueue() *queueProviderImpl
	}
	messageProviderImpl struct {
		qp *queueProviderImpl
		queue amqp.Queue
	}
)



func newQueueManager(connection *amqp.Connection,channel *amqp.Channel) *queueManagerImpl{
	return &queueManagerImpl{
		connection,
		channel,
	}
}



func (qm *queueManagerImpl) createExchangeIfNotExists(exchange exchange, exchangeType exchangeType, autoDelete bool) *queueProviderImpl{
	err := qm.channel.ExchangeDeclare(
		exchange.string(),
		exchangeType.string(),
		false,
		autoDelete,
		false,
		false,
		nil,
	)
	utils.DisposeOnError(err,"cannot create exchange",qm.close)
	return  &queueProviderImpl {
		qm:           qm,
		exchange:     exchange,
		exchangeType: exchangeType,
	}
}

func  (qp *queueProviderImpl) createQueue(name string, autoDelete bool) *messageProviderImpl {
	utils.DisposeOnEmptyString(name,"name is empty",qp.qm.close)
	q, err := qp.qm.channel.QueueDeclare( //automatically creates queue if doesnt exists
		name,       	  //queue name
		false,      //determines if the message should be saved to disk, messages will survive servere restart
		autoDelete,       //what to do with messages if they doesnt have any consumer, true = message will be deleted automatically from the queue, false = keep it
		false,     //exclusive - true = accessible oly from connection that created queue. False queue could be accessible from different connections
		false,      //only return preexisting queue which matches providing configuration, (true) server receive an error when not find, (false) create new queue if doesnt exists on the server
		nil)			//declaring headers
	utils.DisposeOnError(err, "Failed to declare queue", qp.qm.close)
	return &messageProviderImpl{
		qp,
		q,
	}
}


func (p *messageProviderImpl) bindToQueue() *messageProviderImpl {
	//rebinding queue to fanout exchange
	err := p.qp.qm.channel.QueueBind( p.queue.Name,
		"",				//one queue could be bounded to one exchange several times and all bounds will work
		p.qp.exchange.string(),
		false,
		nil,
	)
	utils.DisposeOnError(err,"cannot bind exchange",p.qp.qm.close)
	return p
}



func (p *messageProviderImpl) sendDiscoveryEvent(serviceName string) {
	utils.DisposeOnEmptyString(serviceName,"service name cannot be empty", p.qp.qm.close)
	err := p.qp.qm.channel.Publish(p.qp.exchange.string(), queueEnum(serviceDiscoveryQueue).string() , //serviceName+"_"+uid,
		false,
		false, 	    //if true than throws error when no consumers on the q
		toAmqpMessage(serviceName))
	utils.DisposeOnError(err, "default publishing", p.qp.qm.close)
}



func toAmqpMessage(message string) amqp.Publishing{
	return 	amqp.Publishing{
			Body: []byte(message),
		}
}

func (qm *queueManagerImpl) close(){
	errChannel := qm.channel.Close()
	errConnection := qm.connection.Close()
	utils.FailOnError(errChannel, "unable to close channel")
	utils.FailOnError(errConnection, "unable to close connection")
}

