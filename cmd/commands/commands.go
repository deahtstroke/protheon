package commands

import (
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

type ProtheonCmdFunc func(*app.ProtheonCLI) *cobra.Command

var commands []ProtheonCmdFunc

func RegisterCommand(f ProtheonCmdFunc) {
	commands = append(commands, f)
}

func Commands() []ProtheonCmdFunc {
	return commands
}
