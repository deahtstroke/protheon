package run

import (
	"github.com/deahtstroke/protheon/cli"
	"github.com/deahtstroke/protheon/commands"
	"github.com/spf13/cobra"
)

func init() {
	commands.RegisterCommand(newRunCommand)
}

func newRunCommand(cli cli.ProtheonCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [OPTIONS]",
		Short: "Runs a Protheon ETL pipeline once",
		Long:  "Executes a set of steps with the associated Protheon ETL pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	cmd.Flags().StringP("input", "i", "", "The path of the input file to process")
	cmd.Flags().StringP("format", "f", "", "The extension/format of the input file, e.g., 'json', 'jsonl', 'csv'")
	cmd.Flags().StringP("compression", "c", "", "The compression strategy (if any) for the input file")
	cmd.Flags().StringP("script", "s", "", "The transformation script written in Lua")
	cmd.Flags().StringP("dsn", "d", "", "The datasource name/URI, .e.g, postres://user:password@localhost:5432/db_name")
	cmd.Flags().StringP("table", "t", "", "The table name to insert data to")

	return cmd
}
