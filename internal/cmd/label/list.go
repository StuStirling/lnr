package label

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdList creates the label list command
func NewCmdList() *cobra.Command {
	var teamID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels",
		Long:  "List all issue labels in the Linear workspace.",
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
	labels, err := factory.Client.GetLabels(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed to list labels: %w", err)
	}

	// Warn if results might be truncated (API limit is 100)
	output.WarnIfTruncated(len(labels), 100)

	headers := []string{"NAME", "COLOR", "TEAM"}
	rows := make([][]string, len(labels))
	for i, l := range labels {
		teamName := "-"
		if l.Team != nil {
			teamName = l.Team.Key
		}
		rows[i] = []string{l.Name, l.Color, teamName}
	}

	return factory.Formatter.Print(headers, rows, labels)
}
