package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deahtstroke/protheon/internal/config"
	"github.com/deahtstroke/protheon/internal/db"
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

func NewCreateConfigCommand() *cobra.Command {
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
			return runCreateConfig(opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Format, "format", "f", "TOML", "Config file format to create config. Valid values are YAML and TOML")
	flags.StringVarP(&opts.Alias, "alias", "a", "", "Alias for this configuration")

	cmd.SilenceUsage = true
	return cmd
}

func runCreateConfig(opts CreateConfigOpts) error {
	dbPath := config.GlobalDbDir()
	ctx := context.Background()

	if err := ensureDir(dbPath); err != nil {
		return err
	}

	if err := ensureDir(config.GlobalConfigDataDir()); err != nil {
		return err
	}

	conn, err := db.Connect(fmt.Sprintf("file:%s/%s", dbPath, "protheon.db"))
	if err != nil {
		log.Fatalf("Unable to connect to sqlite DB: %v", err)
	}
	repo := db.NewRepository(conn)
	cfgService := config.NewService(repo)

	res, err := cfgService.CreateConfig(ctx, opts.Format, opts.Alias)
	if err != nil {
		return err
	}

	if res.Alias != "" {
		fmt.Println(res.Alias)
	} else {
		fmt.Println(res.Id)
	}
	return nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(config.GlobalConfigDataDir()); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o775)
	}
	return nil
}
