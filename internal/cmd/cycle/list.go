package cycle

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdList creates the cycle list command
func NewCmdList() *cobra.Command {
	var teamID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List cycles",
		Long:  "List all cycles in the Linear workspace.",
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
	cycles, err := factory.Client.GetCycles(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed to list cycles: %w", err)
	}

	// Warn if results might be truncated (API limit is 50)
	output.WarnIfTruncated(len(cycles), 50)

	headers := []string{"NUMBER", "NAME", "PROGRESS", "TEAM"}
	rows := make([][]string, len(cycles))
	for i, c := range cycles {
		teamName := "-"
		if c.Team != nil {
			teamName = c.Team.Key
		}
		rows[i] = []string{
			fmt.Sprintf("%d", c.Number),
			c.Name,
			output.FormatPercentage(c.Progress),
			teamName,
		}
	}

	return factory.Formatter.Print(headers, rows, cycles)
}
