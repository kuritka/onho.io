package msgbus

import (
	"bytes"
	"encoding/gob"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/streadway/amqp"
)

type (
	IMessageProvider interface {
		EncodeMessage(sm Message) amqp.Publishing
		DecodeMessage(msg amqp.Delivery) Message
		EncodeDisco(DiscoveryRequest) amqp.Publishing
		DecodeDisco(disco amqp.Delivery) DiscoveryRequest
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

func (p *messageProvider) EncodeMessage(sm Message) amqp.Publishing {
	return p.Encode(sm)
}

func (p *messageProvider) EncodeDisco(sm DiscoveryRequest) amqp.Publishing {
	return p.Encode(sm)
}

func  (p *messageProvider) DecodeDisco(msg amqp.Delivery) DiscoveryRequest {
	utils.FailOnNil(msg,"gob nil message")
	buf := bytes.NewBuffer(msg.Body)
	dec := gob.NewDecoder(buf)
	sm := new(DiscoveryRequest)
	err := dec.Decode(sm)
	utils.FailOnError(err,"gob decoding")
	return *sm
}


func (p *messageProvider) Encode(sm interface{}) amqp.Publishing{
	utils.FailOnNil(sm,"gob nil message")
	p.buf.Reset()
	err := gob.NewEncoder(p.buf).Encode(sm)
	utils.FailOnError(err,"gob encoding message")
	return amqp.Publishing{ Body: p.buf.Bytes()}
}

func  (p *messageProvider) DecodeMessage(msg amqp.Delivery) Message {
	utils.FailOnNil(msg,"gob nil message")
	buf := bytes.NewBuffer(msg.Body)
	dec := gob.NewDecoder(buf)
	sm := new(Message)
	err := dec.Decode(sm)
	utils.FailOnError(err,"gob decoding")
	return *sm
}

