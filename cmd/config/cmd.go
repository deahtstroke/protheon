package config

import (
	"github.com/deahtstroke/protheon/cmd/commands"
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

func init() {
	commands.RegisterCommand(newConfigCommand)
}

func newConfigCommand(cli *app.Protheon) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config COMMAND",
		Short: "Manage saved Protheon ETL configs",
		Long:  "Manage saved Protheon ETL configurations that are saved by the user",
		Args:  cobra.ExactArgs(0),
	}

	cmd.AddCommand(NewCreateConfigCommand(cli))
	cmd.AddCommand(NewListConfigCommand(cli))
	cmd.AddCommand(NewDeleteConfigCommand(cli))
	cmd.AddCommand(NewEditConfigCommand(cli))
	return cmd
}
