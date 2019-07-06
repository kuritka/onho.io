package datamanager

import "github.com/kuritka/onho.io/common/dto"

type ISensorRepository interface{
	GetAllSensors() map[string]int
	SaveSensorMessage(msg dto.SensorMessage) error
}