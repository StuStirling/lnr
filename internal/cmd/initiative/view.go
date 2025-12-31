package initiative

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdView creates the initiative view command
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <initiative-id>",
		Short: "View initiative details",
		Long:  "View details of a specific initiative including linked projects.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runView(jsonOutput, args[0])
		},
	}

	return cmd
}

func runView(jsonOutput bool, initiativeID string) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	initiative, err := factory.Client.GetInitiative(ctx, initiativeID)
	if err != nil {
		return fmt.Errorf("failed to get initiative: %w", err)
	}

	ownerName := "-"
	if initiative.Owner != nil {
		ownerName = initiative.Owner.Name
	}

	projectNames := make([]string, len(initiative.Projects))
	for i, p := range initiative.Projects {
		projectNames[i] = p.Name
	}

	fields := []output.DetailField{
		{Label: "ID", Value: initiative.ID},
		{Label: "Name", Value: initiative.Name},
		{Label: "Owner", Value: ownerName},
		{Label: "Target Date", Value: output.EmptyIfNil(initiative.TargetDate)},
		{Label: "Description", Value: initiative.Description},
		{Label: "Projects", Value: strings.Join(projectNames, ", ")},
	}

	return factory.Formatter.PrintDetail(fields, initiative)
}
