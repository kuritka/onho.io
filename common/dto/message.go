package dto

import (
	"encoding/gob"
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

