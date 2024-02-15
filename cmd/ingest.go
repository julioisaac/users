package cmd

import (
	"github.com/julioisaac/users/ioc"
	"github.com/spf13/cobra"
)

var (
	ingestCommand = &cobra.Command{
		Use:   "ingest",
		Short: "Initializes the codebase running as Ingest app",
		Long:  "Initializes the codebase running as Ingest app",
		RunE:  ingestExecute,
	}
)

func init() {
	rootCmd.AddCommand(ingestCommand)

	ingestCommand.Flags().StringP("path", "p", "", "Specify the path (required)")

	err := ingestCommand.MarkFlagRequired("path")
	if err != nil {
		return
	}
}

func ingestExecute(cmd *cobra.Command, args []string) error {
	pathFlag, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}

	is := ioc.IngestService()

	return is.Run(pathFlag)
}
