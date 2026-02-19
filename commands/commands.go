package commands

import "github.com/spf13/cobra"

var commands []func() *cobra.Command

func RegisterCommand(f func() *cobra.Command) {
	commands = append(commands, f)
}

func Commands() []func() *cobra.Command {
	return commands
}
