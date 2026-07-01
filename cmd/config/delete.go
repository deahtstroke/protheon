package config

import (
	"fmt"

	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

func NewDeleteConfigCommand(cli *app.Protheon) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [ARGS]",
		Short: "Delete Protheon ETL configurations",
		Example: `
# Delete a single configuration
protheon config rm 05d1fd988ca1426fc059b83499019fff

# Delete several configurations
# Note you can use aliases or Ids interchangeably
protheon config rm 05d1fd988ca1426fc059b83499019fff my_config
		`,
		Aliases: []string{"rm", "del"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 1 {
				return nil
			}

			return fmt.Errorf("%[1]s requires at least %[2]d %[3]s\n\nUsage: %[4]s\n\nSee '%[1]s --help' for more info",
				cmd.CommandPath(), 1, "argument", cmd.UseLine())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.ConfigService.DeleteConfigs(cli.GlobalCtx, args)
		},
	}
	return cmd
}
