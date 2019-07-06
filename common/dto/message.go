package dto

import (
	"bytes"
	"encoding/gob"
	"github.com/streadway/amqp"
	"time"
)

type SensorMessage struct {
	Session string
	Name string
	Face string
	Value float64
	Timestamp time.Time
}


//any consumer can rely that message will be encoded to gob (much more powerful than json for instance)
func init() {
	gob.Register(SensorMessage{})
}


func FromQueueMessage(msg amqp.Delivery) (SensorMessage, error) {
	r := bytes.NewReader(msg.Body)
	d := gob.NewDecoder(r)
	sensorMessage := new(SensorMessage)
	err := d.Decode(sensorMessage)
	return *sensorMessage, err
}

