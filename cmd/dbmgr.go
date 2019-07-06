package cmd

import (
	"fmt"
	"github.com/kuritka/onho.io/services/datamanger"
	"github.com/spf13/cobra"
)

var
(
	dataManagerOptions datamanager.Options
)

var dataManagerCmd = &cobra.Command{
	Use:   "dbmgr",
	Short: "data manager is wrapper over database",
	Long: `service is writing / providing data in distributed system`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("data manager started..")

		go datamanager.NewService(dataManagerOptions).Run()

		var a string

		fmt.Scanln(&a)

	},
}

func init(){
	dataManagerCmd.Flags().StringVarP(&dataManagerOptions.Name, "name", "n", "", "Source directory to read from")
	dataManagerCmd.Flags().StringVarP(&dataManagerOptions.QueueConnectionString, "connection-string", "c", "", "queue connectionString i.e. amqp://guest:guest@localhost:5672")
	dataManagerCmd.Flags().StringVarP(&dataManagerOptions.DatabaseConnectionString, "database-connection-string", "d", "", "database connectionString i.e. amqp://guest:guest@localhost:5672")
	dataManagerCmd.MarkFlagRequired("Name")
	dataManagerCmd.MarkFlagRequired("connection-string")
	rootCmd.AddCommand(dataManagerCmd)
}