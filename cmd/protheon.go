/*
Copyright © 2025 Daniel Villavicencio dvm3099@pm.me
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/deahtstroke/protheon/cmd/commands"
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/deahtstroke/protheon/internal/config"
	"github.com/deahtstroke/protheon/internal/db"
	"github.com/spf13/cobra"

	_ "github.com/deahtstroke/protheon/cmd/config"
	_ "github.com/deahtstroke/protheon/cmd/run"
)

func createRootCommand(cli *app.Protheon) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "protheon",
		Short: "CLI-based ETL application",
		Long:  `Protheon is an ETL application primarily designed to be used on the terminal`,
	}

	for _, c := range commands.Commands() {
		cmd := c(cli)
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}

// ProtheonMain adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ProtheonMain(ctx context.Context) {
	db, err := db.Connect(config.GlobalDatabaseUrl())
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}

	defer db.Close()

	cli, err := app.New(ctx, db, false)
	if err != nil {
		log.Fatalf("Error while creating CLI: %s", err)
	}

	rCmd := createRootCommand(cli)
	err = rCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
