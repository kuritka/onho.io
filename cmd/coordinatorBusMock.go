package cmd


import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services"
	"github.com/kuritka/onho.io/services/coordinatorBusMock"
	"github.com/spf13/cobra"
)

var (
	coordinatorBusOptions coordinatorBusMock.Options
)

var coordinatorBusMockCmd = &cobra.Command{
	Use:   "coordinator-bus-mock",
	Short: "CoordinatorBusMock is mocking sensor service",
	Long: `A Fast and Flexible face data collector. Authentication is done by github`,

	Run: func(cmd *cobra.Command, args []string) {
		var svc services.IService
		svc = coordinatorBusMock.NewService(coordinatorBusOptions)
		err := svc.Run()
		//err := sensorMock.NewService(sensorOptions).Run()
		utils.FailOnError(err,"Error in SENSOR MOCK")
	},
}

func init(){

	coordinatorBusMockCmd.Flags().StringVarP(&coordinatorBusOptions.Name, "name", "n", "", "Source directory to read from")
	coordinatorBusMockCmd.Flags().StringVarP(&coordinatorBusOptions.QueueConnectionString, "connection-string", "c", "", "connectionString i.e. amqp://guest:guest@localhost:5672")
	coordinatorBusMockCmd.MarkFlagRequired("Name")
	coordinatorBusMockCmd.MarkFlagRequired("connection-string")
	rootCmd.AddCommand(coordinatorBusMockCmd)

}
