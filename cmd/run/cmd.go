package run

import (
	"fmt"
	"log"
	"os"

	"github.com/deahtstroke/protheon/cmd/commands"
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
	inputFlag           string = "input"
	extensionFlag       string = "extension"
	compressionFlag     string = "compression"
	transformationsFlag string = "transformations"
	datasourceUrlFlag   string = "url"
	datasourceTableFlag string = "table"
	configFileFlag      string = "config-file"
)

func newRunCommand(cli *app.ProtheonCLI) *cobra.Command {
	var config engine.ETLConfig
	var configPath string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs a Protheon ETL pipeline once",
		Long:  "Executes a set of steps with the associated Protheon ETL pipeline",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(configFileFlag)
			if err != nil {
				return fmt.Errorf("Error retrieving flag `config`: %s", err)
			}

			if configPath == "" {
				required := []string{inputFlag, extensionFlag, transformationsFlag, datasourceTableFlag, datasourceUrlFlag}
				for _, req := range required {
					f, err := cmd.Flags().GetString(req)
					if err != nil {
						return fmt.Errorf("Error retrieving required flag `%s`: %s", req, err)
					}

					if f == "" {
						return fmt.Errorf("Required flag --%s not set (or provide a config file with --f)", req)
					}
				}
				return nil
			}

			// Parse config if config is not ""
			f, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("Unable to open file with config path %s: %v", configPath, err)
			}

			err = yaml.Unmarshal(f, &config)
			if err != nil {
				err = toml.Unmarshal(f, &config)
			}

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			engine, err := engine.NewEngine()
			if err != nil {
				return err
			}
			defer engine.Cleanup()

			if err := engine.Run(config); err != nil {
				log.Print(err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&config.Input.Path, inputFlag, "i", "", "Path of the input file to process")
	cmd.Flags().StringVarP(&config.Input.Extension, extensionFlag, "e", "", "Extension/format of the input file, e.g., 'json', 'jsonl', 'csv'")
	cmd.Flags().StringVarP(&config.Input.Compress, compressionFlag, "c", "", "Compression strategy of the input file")
	cmd.Flags().StringVarP(&config.Transformations, transformationsFlag, "s", "", "List of Transformation scripts written in Lua")
	cmd.Flags().StringVar(&config.Datasource.URL, datasourceUrlFlag, "", "Datasource URL, .e.g, postres://user:password@localhost:5432/db_name")
	cmd.Flags().StringVarP(&config.Datasource.Table, datasourceTableFlag, "t", "", "Datasource table to insert data to")
	cmd.Flags().StringVarP(&configPath, configFileFlag, "f", "", "Path to config to a YAML/TOML config file for the current ETL run")
	return cmd
}
