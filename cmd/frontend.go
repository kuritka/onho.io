package cmd

import (
	"github.com/kuritka/onho.io/common/manager/depresolver"
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services/frontend"
	"github.com/spf13/cobra"
)


var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "frontend web client",
	Long: `frontend web client`,

	Run: func(cmd *cobra.Command, args []string) {

		//ctx := contextType.Service.SetInContext(rootContext,"FRONTEND SERVICE")
		options := depresolver.NewFromEnvDependencyResolver().
			MustResolveGithubOAuth().
			MustResolveRabbitMQ().
			MustResolveEnvironment().
			MustResolveCookieStore().
			MustResolvePort().
			Options

		//TODO: app needs to be refactored to something like this
		//runner.NewServiceRunner().
		//	WithRest(runner.NewRestRunner("8089").WithDefaultMiddleware().ExportHealthEndpoint()).
		//	MustRun(ctx,options)

		err := frontend.NewService(options).Run()

		utils.FailFastOnError(err)
	},
}

func init(){
	rootCmd.AddCommand(frontendCmd)
}