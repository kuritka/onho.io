package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	name string
	freq uint8
	max float64
	min float64
	stepSize float64
)

var sensorFakeCmd = &cobra.Command{
	Use:   "sensor-fake",
	Short: "SensorFake is authenticated webapp collecting face data",
	Long: `A Fast and Flexible face data collector. Authentication is done by github`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting SESNSOR FAKE CMD...")
	},
}

func init(){
	sensorFakeCmd.Flags().StringVarP(&name, "name", "n", "", "name of the sensor")
	sensorFakeCmd.Flags().Uint8VarP(&freq, "frequency","f", 1, "update frequency in cycles/sec ")
	sensorFakeCmd.Flags().Float64VarP(&min, "min","i", 1., "minimum value for generated readings")
	sensorFakeCmd.Flags().Float64VarP(&max, "max","x", 5., "maximum value for generated readings")
	sensorFakeCmd.Flags().Float64VarP(&stepSize, "stepSize","s", 0.1, "maximum allowable change per measurement")
	sensorFakeCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(sensorFakeCmd)

}