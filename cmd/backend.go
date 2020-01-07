package cmd

import (
	"github.com/kuritka/onho.io/common/manager/depresolver"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services/backend"
	"github.com/spf13/cobra"
)

var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "backend web client",
	Long: `backend web client`,

	Run: func(cmd *cobra.Command, args []string) {

		//ctx := contextType.Service.SetInContext(rootContext,"FRONTEND SERVICE")
		options := depresolver.NewFromEnvDependencyResolver().
			MustResolveEnvironment().
			MustResolvePort().
			Options

		err := backend.NewService(options).Run()

		utils.FailFastOnError(err)
	},
}

func init(){
	rootCmd.AddCommand(backendCmd)
}