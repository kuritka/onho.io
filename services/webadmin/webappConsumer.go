package webadmin

import (
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/qutils"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services"
	"github.com/kuritka/onho.io/services/webadmin/controller"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type webAppConsumer struct {
	eventAggregator services.IEventAggregator
	conn            *amqp.Connection
	ch              *amqp.Channel
	sources         [] string
	consumer        *qutils.MessageConsumer
	provider        *qutils.MessageProvider
}



func NewConsumer(options controller.Options, eventAggregator services.IEventAggregator) *webAppConsumer {
	utils.FailOnNil(eventAggregator,"event aggregator")
	consumer := webAppConsumer{
		eventAggregator: eventAggregator,
	}
	consumer.conn, consumer.ch = qutils.GetChannel(options.QueueConnectionString)
	consumer.consumer = qutils.NewMessageConsumer(consumer.conn, consumer.ch)
	consumer.provider = qutils.NewGobMessageProvider()

	go consumer.ListenForDiscoveryRequest()

	consumer.eventAggregator.AddListener(services.DataSourceDiscovered, func (eventData interface{}){
		consumer.SubscribeToDataEvent(eventData.(string))
	})

	consumer.ch.ExchangeDeclare(
		qutils.WebAppSourceExchange,
		qutils.FanoutKind,
		false,
		false,
		false,
		false,
		nil,
		)

	consumer.ch.ExchangeDeclare(
		qutils.WebAppReadingsExchange,
		qutils.FanoutKind,
		false,
		false,
		false,
		false,
		nil,
	)
	return &consumer
}


func (c *webAppConsumer) ListenForDiscoveryRequest() {
	msgs, err := c.consumer.
		GetPersistentQueue(qutils.WebAppDiscoveryQueue).
		ConsumeFromChannel()
	utils.FailOnError(err, "cannot read from queue")
	for range msgs {

		for _, src := range c.sources {
			log.Info().Msgf("discovered %s", src)
			c.SendMessageSource(src)
		}
	}
}



func (c *webAppConsumer) SendMessageSource(s string) {
	c.provider.AsAmqpMessage(s).PublishToCustomExchange(c.ch, qutils.WebAppSourceExchange,qutils.Default)
}




func (c *webAppConsumer) SubscribeToDataEvent(eventName string) {
	for _,v:= range c.sources {
		if v == eventName {
			return
		}
	}
	c.sources = append(c.sources, eventName)
	c.SendMessageSource(eventName)
	c.eventAggregator.AddListener(services.MessageReceivedPrefix+eventName, func(eventData interface{}){
		ed := eventData.(services.EventData)
		sm := dto.SensorMessage{ Name: ed.Name, Timestamp:ed.Timestamp,Face:ed.Face, Session:ed.Session,  Value:ed.Value}
		c.provider.Encode(sm).AsAmqpMessage().PublishToCustomExchange(c.ch,qutils.WebAppReadingsExchange,qutils.Default)
	})
}