package config

import (
	"github.com/deahtstroke/protheon/commands"
	"github.com/spf13/cobra"
)

func init() {
	commands.RegisterCommand(newConfigCommand)
}

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config COMMAND",
		Short: "Manage saved Protheon ETL configs",
		Long:  "Manage saved Protheon ETL configurations that are saved by the user",
		Args:  cobra.ExactArgs(0),
	}

	cmd.AddCommand(NewCreateConfigCommand())
	return cmd
}
