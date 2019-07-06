package datamanager

import (
	"errors"
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/utils"
)

type SensorRepository struct {
	sensors map[string]int
}


func NewSensorRepository() *SensorRepository{
	sr := SensorRepository{}
	sr.sensors = sr.GetAllSensors()
	return &sr
}

func (sr *SensorRepository) SaveSensorMessage(msg dto.SensorMessage) error {
	if sr.sensors[msg.Name] == 0  {
		return errors.New("Unable to find sensor for name '"+msg.Name+"'")
	}
	_, err := db.Exec("INSERT INTO sensor_reading (value, sensor_id, taken_on) values ($1, $2, $3)",
						msg.Value, sr.sensors[msg.Name], msg.Timestamp)
	return err
}

func (sr *SensorRepository) GetAllSensors() map[string]int {
	sensors := make(map[string]int)
	rows, err := db.Query("SELECT id, name from sensor")
	utils.FailOnError(err, "reading sensors from DB")
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		utils.FailOnError(err, "reading sensor from DB")
		sensors[name] = id
	}
	return sensors
}