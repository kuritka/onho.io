package cmd

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services"
	"github.com/kuritka/onho.io/services/sensorBusMock"
	"github.com/spf13/cobra"
)

var (
	sensorBusOptions sensorBusMock.Options
)

var sensorBusMockCmd = &cobra.Command{
	Use:   "sensor-bus-mock",
	Short: "SensorBusMock is mocking sensor service",
	Long: `A Fast and Flexible face data collector. Authentication is done by github`,

	Run: func(cmd *cobra.Command, args []string) {
		var svc services.IService
		svc = sensorBusMock.NewService(sensorBusOptions)
		err := svc.Run()
		//err := sensorMock.NewService(sensorOptions).Run()
		utils.FailOnError(err,"Error in SENSOR MOCK")
	},
}

func init(){
	sensorBusMockCmd.Flags().StringVarP(&sensorBusOptions.Name, "name", "n", "", "Name of the sensor")
	sensorBusMockCmd.Flags().Uint8VarP(&sensorBusOptions.Freq, "frequency","f", 1, "update frequency in cycles/sec ")
	sensorBusMockCmd.Flags().Float64VarP(&sensorBusOptions.Min, "min","i", 1., "minimum value for generated readings")
	sensorBusMockCmd.Flags().Float64VarP(&sensorBusOptions.Max, "max","x", 5., "maximum value for generated readings")
	sensorBusMockCmd.Flags().Float64VarP(&sensorBusOptions.StepSize, "step-size","s", 0.1, "maximum allowable change per measurement")
	sensorBusMockCmd.Flags().StringVarP(&sensorBusOptions.ConnectionString, "connection-string", "c", "", "connectionString i.e. amqp://guest:guest@localhost:5672")
	sensorBusMockCmd.MarkFlagRequired("Name")
	sensorBusMockCmd.MarkFlagRequired("connection-string")
	rootCmd.AddCommand(sensorBusMockCmd)
}