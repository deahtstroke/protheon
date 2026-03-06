package config

import (
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

func NewEditConfigCommand(cli *app.ProtheonCLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <id_or_alias> [key=value...]",
		Short: "Edits a Protheon ETL configuration",
		Long: `Let's the user edit an ETL configuration based on the Id or alias argument passed. 
If no key-value pairs are passed to directly mutate the file then the user's editor will open the file for editing instead`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			identifier := args[0]
			keyValues := args[:1]
			return cli.ConfigService.EditConfig(cli.GlobalCtx, identifier, keyValues)
		},
	}
	return cmd
}
