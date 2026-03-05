package config

import (
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

func NewConfigListCommand(cli *app.ProtheonCLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all ETL run configurations",
		Long:    "Lists all created ETL run configurations created",
		Aliases: []string{"ls", "ll"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.ConfigService.ListConfigs(cli.GlobalCtx)
		},
	}
	return cmd
}
