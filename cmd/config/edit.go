package config

import (
	"github.com/deahtstroke/protheon/cmd/completion"
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/deahtstroke/protheon/internal/config"
	"github.com/spf13/cobra"
)

func NewEditConfigCommand(cli *app.Protheon) *cobra.Command {
	opts := EditOpts{}

	cmd := &cobra.Command{
		Use:   "edit <id_or_alias> [FLAGS]",
		Short: "Edits a Protheon ETL configuration",
		Long: `Let's the user edit an ETL configuration based on the Id or alias argument passed. 
If no flags are passed to directly manipulate the target config file,
then the user's editor will open the file for interactive editing instead`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.ListConfigs(cli, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			idOrAlias := args[0]
			return runEditConfig(cli, idOrAlias, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Transformations, "transformations", "", "Path to the transformation(s) scripts")
	flags.StringVar(&opts.InputExtension, "extension", "", "Extension of the input, e.g., '.zst', 'json', 'jsonl', etc.")
	flags.StringVar(&opts.InputPath, "input-path", "", "Path to the input source files")
	flags.StringVar(&opts.DatasourceTable, "table", "", "Datasource table destination")
	flags.StringVar(&opts.DatasourceUrl, "url", "", "Datasource URL")

	return cmd
}

func runEditConfig(cli *app.Protheon, id string, opts EditOpts) error {
	editOpts := config.EditConfigurationOptions{
		Input:      &config.ConfigurationInput{},
		Datasource: &config.ConfigurationDatasource{},
	}

	if opts.DatasourceUrl != "" {
		editOpts.Datasource.Url = opts.DatasourceUrl
	}

	if opts.DatasourceTable != "" {
		editOpts.Datasource.Table = opts.DatasourceTable
	}

	if opts.InputPath != "" {
		editOpts.Input.Path = opts.InputPath
	}

	if opts.InputExtension != "" {
		editOpts.Input.Extension = opts.InputExtension
	}

	if opts.Transformations != "" {
		editOpts.TransformationPath = opts.Transformations
	}

	return cli.ConfigService.EditConfig(cli.GlobalCtx, id, editOpts)
}
