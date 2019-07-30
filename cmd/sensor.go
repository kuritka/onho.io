package cmd

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/services/sensor"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)


var (
	logger                = log.Log
	options 			sensor.Options
)

var sensorCmd = &cobra.Command{
	Use:   "sensor",
	Short: "SensorFake fakes  webapp collecting face data",
	Long: `A Fast and Flexible face data collector. Authentication is done by github`,

	Run: func(cmd *cobra.Command, args []string) {

		logger.Info().Msg("sensor mock started")

		sensor.NewService(options).Run()

		var a string

		fmt.Scanln(&a)
	},
}

func init(){
	sensorCmd.Flags().StringVarP(&options.Name, "name", "n", "", "Name of the sensor")
	sensorCmd.Flags().StringVarP(&options.ConnectionString, "connection-string", "c", "", "connectionString i.e. amqp://guest:guest@localhost:5672")
	sensorCmd.MarkFlagRequired("Name")
	sensorCmd.MarkFlagRequired("connection-string")
	rootCmd.AddCommand(sensorCmd)
	BuildContainer()
}


func BuildContainer() *dig.Container{
	container := dig.New()


	return container
}