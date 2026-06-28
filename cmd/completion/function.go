package completion

import (
	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

func ListConfigs(cli *app.ProtheonCLI, limit int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if limit > 0 && len(args) >= limit {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		res, err := cli.ConfigService.ListConfigs(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var namesOrAliases []string
		for _, config := range res {
			namesOrAliases = append(namesOrAliases, config.Id)
		}

		return namesOrAliases, cobra.ShellCompDirectiveNoFileComp
	}
}
