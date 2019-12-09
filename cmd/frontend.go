package cmd

import (
	"github.com/kuritka/onho.io/common/utils"
	"github.com/kuritka/onho.io/services/frontend"
	"github.com/spf13/cobra"
)

var
(
	frontendOptions frontend.Options
)

var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "frontend web client",
	Long: `frontend web client`,

	Run: func(cmd *cobra.Command, args []string) {

		//ctx := contextType.Service.SetInContext(rootContext,"FRONTEND SERVICE")
		//options := depresolver.NewFromEnvDependencyResolver().
		//	MustResolveGithubOAuth().
		//	MustResolveRabbitMQ().
		//	MustResolveEnvironment().
		//	Options


		//runner.NewServiceRunner().
		//	WithRest(runner.NewRestRunner("8089").WithWebServerMiddlewares().ExportHealthEndpoint()).
		//	MustRun(ctx,options)

		err := frontend.NewService(frontendOptions).Run()

		utils.FailFastOnError(err)
	},
}

func init(){
	frontendCmd.Flags().StringVarP(&frontendOptions.ClientID, "client-id", "n", "", "client id i.e. 4fc93087525410d7ae3b")
	frontendCmd.Flags().StringVarP(&frontendOptions.ClientSecret, "client-secret", "s", "", "client secret  i.e. 42efdd9c5afe4f31db0c38fa9b5c1ad443f2aab5")
	frontendCmd.Flags().StringVarP(&frontendOptions.CookieStoreKey, "cookie-key", "k", "", "cookie key i.e. asdifj9182u98u98eabc")
	frontendCmd.Flags().IntVarP(&frontendOptions.Port, "port", "p", 8080, "port i.e. 8081")
	frontendCmd.Flags().StringVarP(&frontendOptions.QueueConnectionString, "connection-string", "c", "", "connectionString i.e. amqp://guest:guest@localhost:5672")
	err := frontendCmd.MarkFlagRequired("client-id")
	utils.FailFastOnError(err)
	err = frontendCmd.MarkFlagRequired("client-secret")
	utils.FailFastOnError(err)
	err = frontendCmd.MarkFlagRequired("cookie-key")
	utils.FailFastOnError(err)
	err = frontendCmd.MarkFlagRequired("port")
	utils.FailFastOnError(err)
	rootCmd.AddCommand(frontendCmd)
}