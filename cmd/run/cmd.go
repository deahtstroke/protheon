package run

import (
	"log"

	"github.com/deahtstroke/protheon/commands"
	"github.com/deahtstroke/protheon/core/engine"
	"github.com/spf13/cobra"
)

func init() {
	commands.RegisterCommand(newRunCommand)
}

func newRunCommand() *cobra.Command {
	var config engine.ExecutorConfig

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs a Protheon ETL pipeline once",
		Long:  "Executes a set of steps with the associated Protheon ETL pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			executor := engine.NewProtheonExecutor(config)
			err := executor.Execute()
			if err != nil {
				log.Print(err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&config.Path, "input", "i", "", "The path of the input file to process")
	cmd.Flags().StringVarP(&config.Format, "format", "f", "", "The extension/format of the input file, e.g., 'json', 'jsonl', 'csv'")
	cmd.Flags().StringVarP(&config.Compress, "compression", "c", "", "The compression strategy (if any) for the input file")
	cmd.Flags().StringVarP(&config.Script, "script", "s", "", "The transformation script written in Lua")
	cmd.Flags().StringVarP(&config.Datasource, "dsn", "d", "", "The datasource name/URI, .e.g, postres://user:password@localhost:5432/db_name")
	cmd.Flags().StringVarP(&config.Table, "table", "t", "", "The table name to insert data to")

	return cmd
}
