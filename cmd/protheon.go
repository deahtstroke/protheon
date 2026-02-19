/*
Copyright © 2025 Daniel Villavicencio dvm3099@pm.me
*/
package cmd

import (
	"context"
	"os"

	"github.com/deahtstroke/protheon/commands"
	"github.com/spf13/cobra"

	_ "github.com/deahtstroke/protheon/cmd/run"
)

func createRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "protheon",
		Short: "CLI-based ETL application",
		Long:  `Protheon is an ETL application primarily designed to be used on the terminal`,
	}

	for _, c := range commands.Commands() {
		cmd := c()
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}

// ProtheonMain adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ProtheonMain(ctx context.Context) {
	rCmd := createRootCommand()

	err := rCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
