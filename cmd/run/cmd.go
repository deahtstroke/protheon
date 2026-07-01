package run

import (
	"log/slog"
	"os"

	"github.com/deahtstroke/protheon/cmd/commands"
	"github.com/deahtstroke/protheon/cmd/completion"
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/deahtstroke/protheon/internal/core/engine"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

func init() {
	commands.RegisterCommand(newRunCommand)
}

const (
	dryRunFlag          string = "dry-run"
	sampleSizeFlag      string = "sample"
	inputFlag           string = "input"
	extensionFlag       string = "extension"
	compressionFlag     string = "compression"
	transformationsFlag string = "transformations"
	datasourceUrlFlag   string = "url"
	datasourceTableFlag string = "table"
	configFileFlag      string = "config-file"
)

func newRunCommand(cli *app.Protheon) *cobra.Command {
	var opts engine.RunOpts
	var dryRun bool
	var pipelinePath string

	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Runs a Protheon ETL pipeline once",
		Long:              "Executes a set of steps with the associated Protheon ETL pipeline",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.ListConfigs(cli, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			engine, err := engine.NewEngine()
			if err != nil {
				return err
			}

			defer engine.Cleanup()

			aliasOrId := args[0]
			if aliasOrId == "" {
				slog.Error("Expecting Id or Alias as an argument")
				return err
			}

			runConfig, err := cli.ConfigService.GetConfigByAliasOrId(cli.GlobalCtx, aliasOrId)
			if err != nil {
				slog.Error("Unable to retrieve run configuration", "AliasOrID", aliasOrId)
				return err
			}

			f, err := os.ReadFile(runConfig.Path)
			if err != nil {
				slog.Error("Unable to read run configuration", "Config", runConfig, "Error", err)
				return err
			}

			err = yaml.Unmarshal(f, &opts)
			if err != nil {
				err = toml.Unmarshal(f, &opts)
				if err != nil {
					return err
				}
			}

			if err := engine.Run(opts, dryRun); err != nil {
				slog.ErrorContext(cli.GlobalCtx, "There was an error running the pipeline", "idOrAlias", aliasOrId, "error", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Input.Path, inputFlag, "i", "", "Path of the input file to process")
	flags.StringVarP(&opts.Input.Extension, extensionFlag, "e", "", "Extension/format of the input file, e.g., 'json', 'jsonl', 'csv'")
	flags.StringVarP(&opts.Input.Compress, compressionFlag, "c", "", "Compression strategy of the input file")
	flags.StringVarP(&opts.Transformations, transformationsFlag, "s", "", "List of Transformation scripts written in Lua")
	flags.StringVar(&opts.Datasource.URL, datasourceUrlFlag, "", "Datasource URL, .e.g, postres://user:password@localhost:5432/db_name")
	flags.StringVarP(&opts.Datasource.Table, datasourceTableFlag, "t", "", "Datasource table to insert data to")
	flags.StringVarP(&pipelinePath, configFileFlag, "f", "", "Path to config to a YAML/TOML config file for the current ETL run")
	flags.BoolVar(&dryRun, dryRunFlag, false, "Whether to execute a dry run")
	flags.Int64Var(&opts.Input.SampleSize, sampleSizeFlag, 0, "Input sample size to run for")

	return cmd
}
