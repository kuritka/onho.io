package cmd

import (
	"fmt"
	"github.com/kuritka/onho.io/services/webadmin"
	"github.com/kuritka/onho.io/services/webadmin/controller"
	"github.com/spf13/cobra"
)

var
(
	webAdminOptions controller.Options
)

var webAdminCmd = &cobra.Command{
	Use:   "webadmin",
	Short: "web admin. It provides overview about users",
	Long: `web admin is used by group administrator. It provides overview about users`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Web admin started..")

		webadmin.NewService(webAdminOptions, sharedEventAggregator).Run()
	},
}

func init(){
	webAdminCmd.Flags().StringVarP(&webAdminOptions.Name, "name", "n", "", "Source directory to read from")
	webAdminCmd.Flags().StringVarP(&webAdminOptions.QueueConnectionString, "connection-string", "c", "", "queue connectionString i.e. amqp://guest:guest@localhost:5672")
	webAdminCmd.Flags().IntVarP(&webAdminOptions.Port, "port", "p", 3000, "Port")
	webAdminCmd.MarkFlagRequired("Name")
	webAdminCmd.MarkFlagRequired("connection-string")
	webAdminCmd.MarkFlagRequired("port")
	rootCmd.AddCommand(webAdminCmd)
}