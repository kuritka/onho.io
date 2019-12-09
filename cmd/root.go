package cmd

import (
	"context"
	"github.com/kuritka/onho.io/services"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Verbose bool
var sharedEventAggregator = services.NewEventAggregator()
var rootContext, _ = context.WithCancel(context.Background())

var rootCmd = &cobra.Command{
	Short: "onho.io",
	Long: `onho.io - Home Office For All`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error().Msg("No parameters included")
			_ = cmd.Help()
			os.Exit(0)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("done..")
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}


func Execute() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
