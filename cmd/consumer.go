package cmd

import (
	"github.com/julioisaac/users/consumer"
	"github.com/spf13/cobra"
)

var (
	consumerCommand = &cobra.Command{
		Use:   "consumer",
		Short: "Initializes the codebase running as Consumer app",
		Long:  "Initializes the codebase running as Consumer app",
		RunE:  consumerExecute,
	}
)

func init() {
	rootCmd.AddCommand(consumerCommand)
}

func consumerExecute(cmd *cobra.Command, args []string) error {
	return consumer.Start()
}
