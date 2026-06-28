package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/deahtstroke/protheon/internal/app"
	"github.com/spf13/cobra"
)

func NewListConfigCommand(cli *app.ProtheonCLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all ETL run configurations",
		Long:    "Lists all Protheon ETL run configurations",
		Aliases: []string{"ls", "ll"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRun(cmd.Context(), cli)
		},
	}
	return cmd
}

func runRun(ctx context.Context, cli *app.ProtheonCLI) error {
	configs, err := cli.ConfigService.ListConfigs(ctx)
	if err != nil {
		return err
	}
	headers := []string{"CONFIG ID", "ALIAS", "CREATED AT"}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	header := strings.Join(headers, "\t")
	fmt.Fprintln(w, header)

	for _, config := range configs {
		t := time.Unix(config.CreatedAt, 0)
		fmt.Fprintf(w, "%s\t%s\t%s\n", config.Id, config.Alias, t.Format(time.RFC1123))
	}

	return nil
}
