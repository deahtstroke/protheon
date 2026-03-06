package config

import (
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

func NewListConfigCommand(cli *app.ProtheonCLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all ETL run configurations",
		Long:    "Lists all Protheon ETL run configurations",
		Aliases: []string{"ls", "ll"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.ConfigService.ListConfigs(cli.GlobalCtx)
		},
	}
	return cmd
}
