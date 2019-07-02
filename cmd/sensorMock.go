package cmd

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services"
	"github.com/kuritka/onho.io/services/sensorMock"
	"github.com/spf13/cobra"
)

var (
	name string
)

var sensorMockCmd = &cobra.Command{
	Use:   "sensor-mock",
	Short: "SensorMock is mocking sensor service",
	Long: `A Fast and Flexible face data collector. Authentication is done by github`,

	Run: func(cmd *cobra.Command, args []string) {
		var svc services.IService
		svc = sensorMock.NewService(sensorOptions)
		err := svc.Run()
		//err := sensorMock.NewService(sensorOptions).Run()
		utils.FailOnError(err,"Error in SENSOR MOCK")
	},
}

func init(){
	sensorMockCmd.Flags().StringVarP(&sensorOptions.Name, "name", "n", "", "Name of the sensor")
	sensorMockCmd.Flags().Uint8VarP(&sensorOptions.Freq, "frequency","f", 1, "update frequency in cycles/sec ")
	sensorMockCmd.Flags().Float64VarP(&sensorOptions.Min, "min","i", 1., "minimum value for generated readings")
	sensorMockCmd.Flags().Float64VarP(&sensorOptions.Max, "max","x", 5., "maximum value for generated readings")
	sensorMockCmd.Flags().Float64VarP(&sensorOptions.StepSize, "step-size","s", 0.1, "maximum allowable change per measurement")
	sensorMockCmd.Flags().StringVarP(&sensorOptions.ConnectionString, "connection-string", "c", "", "connectionString i.e. amqp://guest:guest@localhost:5672")
	sensorMockCmd.MarkFlagRequired("Name")
	sensorMockCmd.MarkFlagRequired("connection-string")
	rootCmd.AddCommand(sensorMockCmd)
}