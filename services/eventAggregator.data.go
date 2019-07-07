package services

import "time"

type EventData struct {
	Session string
	Name string
	Face string
	Value float64
	Timestamp time.Time
}
