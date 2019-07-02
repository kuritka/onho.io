package cmd

import (
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/services/sensorMock"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)


var (
	logger                = log.Log
	sensorOptions sensorMock.Options
)

var sensorCmd = &cobra.Command{
	Use:   "sensor",
	Short: "SensorFake fakes  webapp collecting face data",
	Long: `A Fast and Flexible face data collector. Authentication is done by github`,

	Run: func(cmd *cobra.Command, args []string) {

		logger.Info().Msg("sensor mock started")

		sensorMock.NewService(sensorOptions).Run()
	},
}

func init(){
	sensorCmd.Flags().StringVarP(&sensorOptions.Name, "name", "n", "", "Name of the sensor")
	sensorCmd.Flags().Uint8VarP(&sensorOptions.Freq,  "frequency","f", 1, "update frequency in cycles/sec ")
	sensorCmd.Flags().Float64VarP(&sensorOptions.Min, "min","i", 1., "minimum value for generated readings")
	sensorCmd.Flags().Float64VarP(&sensorOptions.Max, "max","x", 5., "maximum value for generated readings")
	sensorCmd.Flags().Float64VarP(&sensorOptions.StepSize, "step-size","s", 0.1, "maximum allowable change per measurement")
	sensorCmd.MarkFlagRequired("Name")
	rootCmd.AddCommand(sensorCmd)

	BuildContainer()
}


func BuildContainer() *dig.Container{
	container := dig.New()


	return container
}