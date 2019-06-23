package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var coordinator string

var coordinatorCmd = &cobra.Command{
	Use:   "coordinator",
	Short: "coordinator is authenticated webapp collecting face data",
	Long: `A Fast and Flexible face data collector. Authentication is done by github`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("starting coordinator CMD...")
	},
}

func init(){
	coordinatorCmd.Flags().StringVarP(&coordinator, "coordinator", "s", "", "Source directory to read from")
	coordinatorCmd.MarkFlagRequired("coordinator")
	rootCmd.AddCommand(coordinatorCmd)
}