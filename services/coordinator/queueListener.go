package coordinator

import (
	"bytes"
	"encoding/gob"
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/qutils"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)

type IQueueListener interface {
	ListenForNewSource()
	ProcessMessages(msgs <-chan amqp.Delivery)
}

type QueueListener struct {
	conn             *amqp.Connection
	ch               *amqp.Channel
	registry         map[string]<-chan amqp.Delivery
	name             string
	connectionString string
	ea               IEventAggregator
}

func NewQueueListener(options Options) *QueueListener {
	listener := QueueListener{
		registry: make(map[string]<-chan amqp.Delivery),
		name: options.Name,
		connectionString: options.ConnectionString,
		ea: NewEventAggregator(),
	}
	listener.conn,listener.ch = qutils.GetChannel(options.ConnectionString)
	return &listener
}


//ListenForNewSource is method discover new sensors
func (l *QueueListener) ListenForNewSource() {

	fanoutChannel, err := qutils.
		NewMessageConsumer(l.conn, l.ch).
		GetUniqueQueue().
		BindToFanout().
		ConsumeFromChannel()
	utils.FailOnError(err, "unable prepare fanout channel for reading")

	l.DiscoverSensors()

	for msg := range fanoutChannel {

		sensorId := string(msg.Body)
		dataChannel, _ := l.ch.Consume(
			sensorId,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		if l.registry[sensorId] == nil {
			l.registry[sensorId] = dataChannel

			go l.ProcessMessages(dataChannel)
		}
	}
}


func (l *QueueListener) ProcessMessages(msgs <-chan amqp.Delivery){
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sensorMessage := new(dto.SensorMessage)
		err := d.Decode(sensorMessage)
		utils.FailOnError(err, "decoding message")
		logger.Debug().Msgf("Received message: %v\n", sensorMessage)

		data := EventData{
			Name: sensorMessage.Name,
			Value: sensorMessage.Value,
			Session: sensorMessage.Session,
			Face: sensorMessage.Face,
			Timestamp: sensorMessage.Timestamp,
		}
		l.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, data)
	}
}


func (l *QueueListener) DiscoverSensors(){
	err := l.ch.ExchangeDeclare(qutils.SensorDiscoveryExchange,
		qutils.FanoutKind,
		false,
		false,
		false,
		false,
		nil,
		)
	utils.FailOnError(err,"can't create SensorDiscoveryExchange")
	err = l.ch.Publish(qutils.SensorDiscoveryExchange,
		"",
		false,
		false,
		amqp.Publishing{})
	utils.FailOnError(err,"can't publish request to SensorDiscoveryExchange")
}