package datamanager

import (
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/qutils"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)



type DbController struct {
	options          Options
	consumer         *qutils.MessageConsumer
	sensorRepository *SensorRepository
}


func NewDbController(options Options) *DbController{
	sensorRepository := NewSensorRepository()
	return &DbController{options:options, sensorRepository: sensorRepository}
}


func (dbc *DbController) ReadPersistentQueue(){
	conn, channel := qutils.GetChannel(dbc.options.QueueConnectionString)
	defer conn.Close()
	defer channel.Close()
	consumer := qutils.NewMessageConsumer(conn, channel)

	//I want consume exclusive because only one datamanager is here (no massive concurent data streams)
	msgs, err  := consumer.GetPersistentQueue(qutils.PersistentReadingsQueue).ConsumeExclusiveFromChannel(false)
	utils.FailOnError(err,"reading from "+qutils.PersistentReadingsQueue)

	dbc.readMessages(msgs)
}

func (dbc *DbController) readMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		sensorMessage, err := dto.FromQueueMessage(msg)
		err = dbc.sensorRepository.SaveSensorMessage(sensorMessage)
		if err != nil {
			logger.Error().Msgf("failed to save message from sensor %v; %s", sensorMessage.Name, err.Error())
		} else {
			msg.Ack(false)
			logger.Debug().Msgf("Save: %v\n", sensorMessage)
		}
	}
}
