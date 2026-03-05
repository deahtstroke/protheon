package config

import (
	"context"
	"fmt"
	"log"

	"github.com/deahtstroke/protheon/internal/config"
	"github.com/deahtstroke/protheon/internal/db"
	"github.com/spf13/cobra"
)

func NewConfigListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all ETL run configurations",
		Long:    "Lists all created ETL run configurations created",
		Aliases: []string{"ls", "ll"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListConfigs()
		},
	}
	return cmd
}

func runListConfigs() error {
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

	if err := cfgService.ListConfigs(ctx); err != nil {
		return err
	}

	return nil
}
