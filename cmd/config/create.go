package config

import (
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

type CreateConfigOpts struct {
	// The format of the config to be created for Protheon
	// Valid values are: TOML and YAML
	Format string

	// An alias that these configuration should have
	//
	// Aliases can be used in place of the auto-generated identifier to call a config
	// with the CLI
	Alias string
}

func NewCreateConfigCommand(cli *app.Protheon) *cobra.Command {
	var opts CreateConfigOpts
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an ETL run configuration",
		Example: `
# Create a YAML configuration
protheon config create --format YAML

# Create a TOML configuration with an alias name of 'skibidi'
protheon config create -f toml -a skibidi`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.ConfigService.CreateConfig(cli.GlobalCtx, opts.Format, opts.Alias)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Format, "format", "f", "TOML", "Config file format to create config. Valid values are YAML and TOML")
	flags.StringVarP(&opts.Alias, "alias", "a", "", "Alias for this configuration")

	cmd.SilenceUsage = true
	return cmd
}
