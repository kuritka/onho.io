package msgbus

import (
	"bytes"
	"encoding/gob"
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)

type (
	IMessageProvider interface {
		Encode(message dto.SensorMessage) *amqp.Publishing
		Decode(message amqp.Delivery) Message
	}

	messageProvider struct {
		buf *bytes.Buffer
	}
)

func newGobMessageProvider() *messageProvider{
	buf := new (bytes.Buffer)
	return &messageProvider{
		buf,
	}
}

func (p *messageProvider) Encode(sm Message) amqp.Publishing{
	utils.FailOnNil(sm,"gob nil message")
	p.buf.Reset()
	err := gob.NewEncoder(p.buf).Encode(sm)
	utils.FailOnError(err,"gob encoding message")
	return amqp.Publishing{ Body: p.buf.Bytes()}
}

func  (p *messageProvider) Decode(msg amqp.Delivery) Message {
	utils.FailOnNil(msg,"gob nil message")
	buf := bytes.NewBuffer(msg.Body)
	dec := gob.NewDecoder(buf)
	sm := new(Message)
	err := dec.Decode(sm)
	utils.FailOnError(err,"gob decoding")
	return *sm
}

