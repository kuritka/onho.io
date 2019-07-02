package sensorMock

import (
	"github.com/streadway/amqp"
	"math/rand"
	"strconv"
	"time"

	"github.com/kuritka/onho.io/common/dto"
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/common/qutils"
	"github.com/kuritka/onho.io/common/utils"
)

var logger   = log.Log

type SensorMock struct {
	options Options
	rand    *rand.Rand
	value   float64
	nom     float64

}

func NewService(options Options) *SensorMock {

	utils.FailOnEmptyString(options.Name,"empty name")
	utils.FailOnEmptyString(options.ConnectionString,"empty connection string")

	return &SensorMock{
		options,
		rand.New(rand.NewSource(time.Now().UnixNano())),
		rand.Float64()*(options.Max-options.Min) + options.Min,
		(options.Max-options.Min)/2 + options.Min,
	}
}


func (sm *SensorMock) Run() error {

	duration,err := time.ParseDuration(strconv.Itoa(1000/int(sm.options.Freq)) + "ms")

	utils.FailOnError(err,"parse duration")

	signal := time.Tick(duration)

	conn, ch := qutils.GetChannel(sm.options.ConnectionString)
	defer conn.Close()
	defer ch.Close()

	provider :=  qutils.NewGobMessageProvider()

	discoveryRequests, err := qutils.NewMessageConsumer(conn, ch).
	 	GetUniqueQueue().
	 	BindToExchange(qutils.SensorDiscoveryExchange).
		ConsumeFromChannel()
	utils.FailOnError(err, "discovery exchange")

	go sm.listenForDiscoverRequests(discoveryRequests, provider ,ch)

	provider.AsAmqpMessage(sm.options.Name).PublishQueueNameToFanout(ch)

	dataQueue := qutils.GetQueue(sm.options.Name,ch,false)

	for range signal {
		sm.calcValue()

		reading := dto.SensorMessage{
			Name:      sm.options.Name,
			Value:     sm.value,
			Timestamp: time.Now(),
			Face: "HAPPY-FACE",
			Session: "-",
		}

		provider.Encode(reading).AsAmqpMessage().PublishDefault(ch,dataQueue)

		logger.Info().Msgf("NAME=%s VALUE=%v",reading.Name, reading.Value)
	}

	logger.Info().Msgf("Reading sent. %s Value: %v\n",sm.options.Name, sm.value)

	return nil
}

func (sm *SensorMock) listenForDiscoverRequests(deliveries <-chan amqp.Delivery,provider qutils.MessageProvider, ch *amqp.Channel) {
	for range deliveries {
		provider.AsAmqpMessage(sm.options.Name).PublishQueueNameToFanout(ch)
	}
}

func (sm *SensorMock) calcValue() {
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