/*
Copyright © 2025 Daniel Villavicencio dvm3099@pm.me
*/
package cmd

import (
	"context"
	"os"

	"github.com/deahtstroke/protheon/cli"
	"github.com/deahtstroke/protheon/commands"
	"github.com/spf13/cobra"
)

func createRootCommand(protheonCLI cli.ProtheonCli) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "protheon",
		Short: "Distributed task executor",
		Long: `Protheon is a distributed system that allows the user
			to manage tasks that are running in other hosts`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	for _, c := range commands.Commands() {
		cmd := c(protheonCLI)
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}

// ProtheonMain adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ProtheonMain(ctx context.Context) {
	cli := cli.NewCli(ctx)
	rCmd := createRootCommand(cli)

	err := rCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
