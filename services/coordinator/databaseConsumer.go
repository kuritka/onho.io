package coordinator

import (
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/qutils"
	"github.com/kuritka/onho.io/services"
	"github.com/streadway/amqp"
	"time"
)

const maxRate = 10 * time.Second

type DatabaseConsumer struct {
	ea                      services.IEventAggregator
	conn                    *amqp.Connection
	ch                      *amqp.Channel
	persistentReadingsQueue *amqp.Queue
	messageProvider         *qutils.MessageProvider
	sources 				[]string
}


func NewDatabaseConsumer(ea services.IEventAggregator, connectionString string) *DatabaseConsumer{
	dc := DatabaseConsumer {ea: ea}
	dc.conn, dc.ch = qutils.GetChannel(connectionString)
	dc.persistentReadingsQueue = qutils.GetQueue(qutils.PersistentReadingsQueue,dc.ch,false)
	dc.messageProvider = qutils.NewGobMessageProvider()
	ea.AddListener(services.DataSourceDiscovered, func(eventName interface{}){
		dc.SubscribeToDataEvent(eventName.(string))
	})
	return &dc
}


//any time new sensor is added , SubscribeToDataEvent is called and adds
//new listener for deque items from data queue (MessageReceivedPrefix+sensorName
// i.e. calling this listener adds new listener MessageReceivedPrefixAAA publishing when new message in data channel comes )
func (dc *DatabaseConsumer) SubscribeToDataEvent(sensorName string) {
	for _,v := range dc.sources {
		if v == sensorName {
			return
		}
	}

	//self-executing method parameter func() func(interface{}){}()
	//if event is published it checks whether previous message is older than 5 sec and if yes, message is stored
	dc.ea.AddListener(services.MessageReceivedPrefix+sensorName, func() func(interface{}){
		prevTime  := time.Unix(0,0)
		return func(eventData interface{}){
			ed := eventData.(services.EventData)
			//only message newer than maxRate seconds is written. This counts for each sensor and coordinator.
			//if multiple coordinators exists than possibility that all of them will push message at same time.
			if time.Since(prevTime) > maxRate {
				prevTime = time.Now()
				sm := dto.SensorMessage{ Name: ed.Name, Timestamp:ed.Timestamp,Face:ed.Face, Session:ed.Session,  Value:ed.Value}
				dc.messageProvider.Encode(sm).AsAmqpMessage().PublishDefault(dc.ch,dc.persistentReadingsQueue)
			}
		}
	}())
}