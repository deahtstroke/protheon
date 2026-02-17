package commands

import (
	"github.com/deahtstroke/protheon/cli"
	"github.com/spf13/cobra"
)

type ProtheonCommand func(cli cli.ProtheonCli) *cobra.Command

var commands []ProtheonCommand

func RegisterCommand(f ProtheonCommand) {
	commands = append(commands, f)
}

func Commands() []ProtheonCommand {
	return commands
}
