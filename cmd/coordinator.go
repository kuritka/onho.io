package cmd

import (
	"fmt"
	"github.com/kuritka/onho.io/common/utils"

	"github.com/kuritka/onho.io/services/coordinator"
	"github.com/spf13/cobra"
)

var
(
	coordinatorOptions coordinator.Options
)

var coordinatorCmd = &cobra.Command{
	Use:   "coordinator",
	Short: "coordinator is authenticated webapp collecting face data",
	Long: `A Fast and Flexible face data collector. Authentication is done by github`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("coordinator started..")

		err := coordinator.NewService(coordinatorOptions, sharedEventAggregator).Run()

		utils.FailFastOnError(err)
	},
}

func init(){
	coordinatorCmd.Flags().StringVarP(&coordinatorOptions.Name, "name", "n", "", "Source directory to read from")
	coordinatorCmd.Flags().StringVarP(&coordinatorOptions.QueueConnectionString, "connection-string", "c", "", "connectionString i.e. amqp://guest:guest@localhost:5672")
	coordinatorCmd.MarkFlagRequired("Name")
	coordinatorCmd.MarkFlagRequired("connection-string")
	rootCmd.AddCommand(coordinatorCmd)
}