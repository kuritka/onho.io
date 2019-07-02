package qutils

import (
	"bytes"
	"encoding/gob"
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)

type MessageProvider struct{
	buf *bytes.Buffer
}

type AmqpMessage struct {
	buf *bytes.Buffer
}


type AmqpPublishing struct {
	message *amqp.Publishing
}

func NewGobMessageProvider() MessageProvider {
	buf := new (bytes.Buffer)
	return MessageProvider{
		buf,
	}
}

func (p *MessageProvider) Encode(m dto.SensorMessage) *AmqpMessage{
	return p.encode(m)
}


func (p *MessageProvider) AsAmqpMessage(message string) *AmqpPublishing {
	return &AmqpPublishing{
		&amqp.Publishing{
			Body: []byte(message),
		},
	}
}


func (a *AmqpMessage) AsAmqpMessage() *AmqpPublishing {
	return &AmqpPublishing {
		&amqp.Publishing {
			Body: a.buf.Bytes(),
		},
	}
}


func (pub *AmqpPublishing) PublishDefault(ch *amqp.Channel, q *amqp.Queue) {
	utils.FailOnNil(ch,"channel is nil")
	utils.FailOnNil(q,"queue is nil")
	err := ch.Publish(Default,q.Name,
		false,
		false, 	    //if true than throws error when no consumers on the q
		*pub.message)
	utils.FailOnError(err, "default publishing")
}

func (pub *AmqpPublishing) PublishQueueNameToFanout(ch *amqp.Channel) {
	utils.FailOnNil(ch,"channel is nil")
	err := ch.Publish(Fanout,"",
		false,
		false, 	    //if true than throws error when no consumers on the q
		*pub.message)
	utils.FailOnError(err, "default publishing")
}



func (p *MessageProvider) encode(m interface{}) *AmqpMessage{
	utils.FailOnNil(m,"nil message")
	p.buf.Reset()
	err := gob.NewEncoder(p.buf).Encode(m)
	utils.FailOnError(err,"encoding message")
	return &AmqpMessage{ p.buf }
}