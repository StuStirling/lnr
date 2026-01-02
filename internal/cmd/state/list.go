package state

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdList creates the state list command
func NewCmdList() *cobra.Command {
	var teamID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workflow states",
		Long:  "List all workflow states in the Linear workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			var teamPtr *string
			if teamID != "" {
				teamPtr = &teamID
			}
			return runList(jsonOutput, teamPtr)
		},
	}

	cmd.Flags().StringVar(&teamID, "team", "", "Filter by team ID")

	return cmd
}

func runList(jsonOutput bool, teamID *string) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	states, err := factory.Client.GetWorkflowStates(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed to list workflow states: %w", err)
	}

	// Warn if results might be truncated (API limit is 100)
	output.WarnIfTruncated(len(states), 100)

	headers := []string{"NAME", "TYPE", "TEAM"}
	rows := make([][]string, len(states))
	for i, s := range states {
		teamName := "-"
		if s.Team != nil {
			teamName = s.Team.Key
		}
		rows[i] = []string{s.Name, s.Type, teamName}
	}

	return factory.Formatter.Print(headers, rows, states)
}
