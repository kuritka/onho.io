package sensorBusMock

import (
	"fmt"
	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/msgbus"
	"github.com/kuritka/onho.io/common/utils"
	"math/rand"
	"strconv"
	"time"
)


type SensorBusMock struct {
	options Options
	rand    *rand.Rand
	value   float64
	nom     float64
}

func NewService(options Options) *SensorBusMock {

	utils.FailOnEmptyString(options.Name,"empty name")
	utils.FailOnEmptyString(options.ConnectionString,"empty connection string")

	return &SensorBusMock{
		options,
		rand.New(rand.NewSource(time.Now().UnixNano())),
		rand.Float64()*(options.Max-options.Min) + options.Min,
		(options.Max-options.Min)/2 + options.Min,
	}
}



func (sm *SensorBusMock) Run() error {

	mb := msgbus.NewMsgBus(sm.options.ConnectionString)
	defer mb.Close()

	listener, publisher := mb.Register("serviceBus")
	listener.Listen()

	duration,err := time.ParseDuration(strconv.Itoa(1000/int(sm.options.Freq)) + "ms")

	utils.FailOnError(err,"parse duration")

	signal := time.Tick(duration)


	for range signal {
		sm.calcValue()

		reading := dto.SensorMessage{
			Name:      sm.options.Name,
			Value:     sm.value,
			Timestamp: time.Now(),
			Face: "HAPPY-FACE",
			Session: "-",
		}
		data := fmt.Sprintf("NAME=%s VALUE=%v",reading.Name, reading.Value)
		publisher.Command("CMD-A", data)
		fmt.Println(data)
	}
	return nil
}



func (sm *SensorBusMock) calcValue() {
	var maxStep, minStep float64

	if sm.value< sm.nom {
		maxStep = sm.options.StepSize
		minStep = -1 * sm.options.StepSize * (sm.value - sm.options.Min) / (sm.nom - sm.options.Min)
	} else {
		maxStep = sm.options.StepSize * (sm.options.Max - sm.value) / (sm.options.Max  - sm.nom)
		minStep = -1 * sm.options.StepSize
	}
	sm.value += sm.rand.Float64()*(maxStep-minStep) + minStep
}