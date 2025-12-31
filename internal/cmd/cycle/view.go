package cycle

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdView creates the cycle view command
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <cycle-id>",
		Short: "View cycle details",
		Long:  "View details of a specific cycle.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runView(jsonOutput, args[0])
		},
	}

	return cmd
}

func runView(jsonOutput bool, cycleID string) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cycle, err := factory.Client.GetCycle(ctx, cycleID)
	if err != nil {
		return fmt.Errorf("failed to get cycle: %w", err)
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
