package cycle

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/api"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdActive creates the cycle active command
func NewCmdActive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active <team-id>",
		Short: "Show active cycle",
		Long:  "Show the currently active cycle for a team.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runActive(jsonOutput, args[0])
		},
	}

	return cmd
}

func runActive(jsonOutput bool, teamID string) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cycle, err := factory.Client.GetActiveCycle(ctx, teamID)
	if err != nil {
		if errors.Is(err, api.ErrNoActiveCycle) {
			fmt.Println("No active cycle for this team.")
			return nil
		}
		return fmt.Errorf("failed to get active cycle: %w", err)
	}

	fields := []output.DetailField{
		{Label: "ID", Value: cycle.ID},
		{Label: "Name", Value: cycle.Name},
		{Label: "Number", Value: fmt.Sprintf("%d", cycle.Number)},
		{Label: "Progress", Value: output.FormatPercentage(cycle.Progress)},
		{Label: "Team", Value: cycle.Team.Name},
		{Label: "Description", Value: cycle.Description},
	}

	return factory.Formatter.PrintDetail(fields, cycle)
}
