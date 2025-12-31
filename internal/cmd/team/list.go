package team

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdList creates the team list command
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List teams",
		Long:  "List all teams in the Linear workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runList(jsonOutput)
		},
	}

	return cmd
}

func runList(jsonOutput bool) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	teams, err := factory.Client.GetTeams(ctx)
	if err != nil {
		return fmt.Errorf("failed to list teams: %w", err)
	}

	headers := []string{"KEY", "NAME", "DESCRIPTION"}
	rows := make([][]string, len(teams))
	for i, t := range teams {
		rows[i] = []string{
			t.Key,
			t.Name,
			output.Truncate(t.Description, 50),
		}
	}

	return factory.Formatter.Print(headers, rows, teams)
}
