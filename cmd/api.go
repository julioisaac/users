package cmd

import (
	"github.com/julioisaac/users/api"
	"github.com/spf13/cobra"
)

var (
	apiCommand = &cobra.Command{
		Use:   "api",
		Short: "Initializes the codebase running as HTTP API",
		Long:  "Initializes the codebase running as HTTP API",
		RunE:  apiExecute,
	}
)

func init() {
	rootCmd.AddCommand(apiCommand)
}

func apiExecute(cmd *cobra.Command, args []string) error {
	return api.StartAPI()
}
