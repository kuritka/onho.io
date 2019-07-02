package qutils

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)


//SensorListQueue is well known queue used cross whole application. New sensor is telling the others that exists....
//const SensorListQueue = "SensorListQueue"


func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	utils.FailOnError(err, "Failed connect to RabitMQ")
	channel, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	return conn, channel
}

//Autodelete  what to do with messages if they doesnt have any consumer, true = message will be deleted automatically from the queue, false = keep it
//Would be true for all discovery queues
func GetQueue(name string, channel *amqp.Channel, autoDelete bool) *amqp.Queue {
	q, err := channel.QueueDeclare( //automatically creates queue if doesnt exists
		name,       	  //queue name
		false,      //determines if the message should be saved to disk, messages will survive servere restart
		autoDelete,       //what to do with messages if they doesnt have any consumer, true = message will be deleted automatically from the queue, false = keep it
		false,     //exclusive - true = accessible oly from connection that created queue. False queue could be accessible from different connections
		false,      //only return preexisting queue which matches providing configuration, (true) server receive an error when not find, (false) create new queue if doesnt exists on the server
		nil)			//declaring headers
	utils.FailOnError(err, "Failed to declare queue")
	return &q
}

