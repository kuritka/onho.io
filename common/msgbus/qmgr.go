package msgbus

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type (
	iExchangeManager interface {
		createDiscoveryExchangeIfNotExists() *exchangeManagerImpl
		createQueueIfNotExists(name string, autoDelete bool) *queueManagerImpl
		sendDiscoveryRequest(message amqp.Publishing) *exchangeManagerImpl
		close()
	}

	exchangeManagerImpl struct {
		connection *amqp.Connection
		channel    *amqp.Channel
	}
)

type (
	iQueueManager interface {
		bindToQueue(exchange exchange) *queueManagerImpl
		publishMessage(exchange string, queueName string, message amqp.Publishing)
		close()
	}

	queueManagerImpl struct {
		connection *amqp.Connection
		channel    *amqp.Channel
		queue      *amqp.Queue
	}
)

func newExchangeManager(connection *amqp.Connection, channel *amqp.Channel) *exchangeManagerImpl {
	return &exchangeManagerImpl{
		connection,
		channel,
	}
}

func (em *exchangeManagerImpl) createDiscoveryExchangeIfNotExists() *exchangeManagerImpl {
	err := em.channel.ExchangeDeclare(
		exchange.string(serviceDiscoveryExchange),
		amqp.ExchangeFanout,
		false,
		true,
		false,
		false,
		nil,
	)
	utils.DisposeOnError(err, "cannot create discovery exchange", em.close)
	return em
}

func (em *exchangeManagerImpl) createEventExchangeIfNotExists() *exchangeManagerImpl {
	err := em.channel.ExchangeDeclare(
		exchange.string(serviceEventExchange),
		amqp.ExchangeDirect,
		false,
		true,
		false,
		false,
		nil,
	)
	utils.DisposeOnError(err, "cannot create event exchange", em.close)
	return em
}

func (em *exchangeManagerImpl) createCommandExchangeIfNotExists() *exchangeManagerImpl {
	err := em.channel.ExchangeDeclare(
		exchange.string(serviceCommandExchange),
		amqp.ExchangeDirect,
		false,
		true,
		false,
		false,
		nil,
	)
	utils.DisposeOnError(err, "cannot create event exchange", em.close)
	return em
}

func (em *exchangeManagerImpl) sendDiscoveryRequest(message amqp.Publishing) *exchangeManagerImpl {
	err := em.channel.Publish(exchange.string(serviceDiscoveryExchange),
		"",
		false,
		false, //if true than throws error when no consumers on the q
		message)
	utils.FailOnError(err, "default publishing")
	return em
}

func (em *exchangeManagerImpl) createQueueIfNotExists(name string, autoDelete bool) *queueManagerImpl {
	q, err := em.channel.QueueDeclare( //automatically creates queue if doesnt exists
		name,       //queue name
		false,      //determines if the message should be saved to disk, messages will survive servere restart
		autoDelete, //what to do with messages if they doesnt have any consumer, true = message will be deleted automatically from the queue, false = keep it
		false,      //exclusive - true = accessible oly from connection that created queue. False queue could be accessible from different connections
		false,      //only return preexisting queue which matches providing configuration, (true) server receive an error when not find, (false) create new queue if doesnt exists on the server
		nil)        //declaring headers
	utils.FailOnError(err, "Failed to declare queue")
	return &queueManagerImpl{
		em.connection,
		em.channel,
		&q,
	}
}

func (em *exchangeManagerImpl) close() {
	dispose(em.connection, em.channel)
}

func createQueueManager(conn *amqp.Connection, ch *amqp.Channel) *queueManagerImpl {
	return &queueManagerImpl{
		conn,
		ch,
		nil,
	}
}

func (qm *queueManagerImpl) createQueueIfNotExists(name string, autoDelete bool) *queueManagerImpl {
	q, err := qm.channel.QueueDeclare( //automatically creates queue if doesnt exists
		name,       //queue name
		false,      //determines if the message should be saved to disk, messages will survive servere restart
		autoDelete, //what to do with messages if they doesnt have any consumer, true = message will be deleted automatically from the queue, false = keep it
		false,      //exclusive - true = accessible oly from connection that created queue. False queue could be accessible from different connections
		false,      //only return preexisting queue which matches providing configuration, (true) server receive an error when not find, (false) create new queue if doesnt exists on the server
		nil)        //declaring headers
	utils.FailOnError(err, "Failed to declare queue")
	qm.queue = &q
	return qm
}

func (qm *queueManagerImpl) bindToQueue(routingKey string, exchange exchange) *queueManagerImpl {
	utils.FailOnNil(qm.queue, "binding to nil queue")
	//rebinding queue to fanout exchange
	err := qm.channel.QueueBind(qm.queue.Name,
		routingKey, //one queue could be bounded to one exchange several times and all bounds will work
		exchange.string(),
		false,
		nil,
	)
	utils.DisposeOnError(err, "cannot bind exchange", qm.close)
	return qm
}

func (qm *queueManagerImpl) publishMessage(exchange string, routingKey string, message amqp.Publishing) {
	err := qm.channel.Publish(exchange, routingKey,
		false,
		false, //if true than throws error when no consumers on the q
		message)
	utils.FailOnError(err, "default publishing")
}

func (qm *queueManagerImpl) consumeFromChannel() (<-chan amqp.Delivery, error) {
	return qm.channel.Consume(qm.queue.Name,
		"", true, false, false, false, nil)
}

func (qm *queueManagerImpl) close() {
	dispose(qm.connection, qm.channel)
}

func dispose(connection *amqp.Connection, channel *amqp.Channel) {
	utils.FailOnNil(connection, "connection is nil")
	utils.FailOnNil(channel, "channel is nil")
	var err error
	err = channel.Close()
	utils.FailOnError(err, "unable to dispose channel")
	err = connection.Close()
	utils.FailOnError(err, "unable to dispose connection")
	log.Debug().Msg("connection closed")
}
