package team

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdView creates the team view command
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <team-id>",
		Short: "View team details",
		Long:  "View details of a specific team.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runView(jsonOutput, args[0])
		},
	}

	return cmd
}

func runView(jsonOutput bool, teamID string) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	team, err := factory.Client.GetTeam(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed to get team: %w", err)
	}

	privateStr := "No"
	if team.Private {
		privateStr = "Yes"
	}

	fields := []output.DetailField{
		{Label: "ID", Value: team.ID},
		{Label: "Name", Value: team.Name},
		{Label: "Key", Value: team.Key},
		{Label: "Description", Value: team.Description},
		{Label: "Private", Value: privateStr},
	}

	return factory.Formatter.PrintDetail(fields, team)
}
