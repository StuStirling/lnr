package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdView creates the project view command
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <project-id>",
		Short: "View project details",
		Long:  "View details of a specific project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runView(jsonOutput, args[0])
		},
	}

	return cmd
}

func runView(jsonOutput bool, projectID string) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	project, err := factory.Client.GetProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	leadName := "-"
	if project.Lead != nil {
		leadName = project.Lead.Name
	}

	teamKeys := make([]string, len(project.Teams))
	for i, t := range project.Teams {
		teamKeys[i] = t.Key
	}

	fields := []output.DetailField{
		{Label: "ID", Value: project.ID},
		{Label: "Name", Value: project.Name},
		{Label: "State", Value: project.State},
		{Label: "Progress", Value: output.FormatPercentage(project.Progress)},
		{Label: "Lead", Value: leadName},
		{Label: "Teams", Value: strings.Join(teamKeys, ", ")},
		{Label: "Start Date", Value: output.EmptyIfNil(project.StartDate)},
		{Label: "Target Date", Value: output.EmptyIfNil(project.TargetDate)},
		{Label: "Description", Value: project.Description},
		{Label: "URL", Value: project.URL},
	}

	return factory.Formatter.PrintDetail(fields, project)
}
