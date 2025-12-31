package issue

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdView creates the issue view command
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <issue-id>",
		Short: "View issue details",
		Long:  "View details of a specific issue by ID or identifier (e.g., ENG-123).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runView(jsonOutput, args[0])
		},
	}

	return cmd
}

func runView(jsonOutput bool, issueID string) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	issue, err := factory.Client.GetIssue(ctx, issueID)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	stateName := "-"
	if issue.State != nil {
		stateName = issue.State.Name
	}

	assigneeName := "-"
	if issue.Assignee != nil {
		assigneeName = issue.Assignee.Name
	}

	creatorName := "-"
	if issue.Creator != nil {
		creatorName = issue.Creator.Name
	}

	projectName := "-"
	if issue.Project != nil {
		projectName = issue.Project.Name
	}

	cycleName := "-"
	if issue.Cycle != nil {
		cycleName = issue.Cycle.Name
	}

	labelNames := make([]string, len(issue.Labels))
	for i, l := range issue.Labels {
		labelNames[i] = l.Name
	}

	fields := []output.DetailField{
		{Label: "Identifier", Value: issue.Identifier},
		{Label: "Title", Value: issue.Title},
		{Label: "State", Value: stateName},
		{Label: "Priority", Value: output.PriorityLabel(issue.Priority)},
		{Label: "Assignee", Value: assigneeName},
		{Label: "Creator", Value: creatorName},
		{Label: "Team", Value: issue.Team.Name},
		{Label: "Project", Value: projectName},
		{Label: "Cycle", Value: cycleName},
		{Label: "Labels", Value: strings.Join(labelNames, ", ")},
		{Label: "Due Date", Value: output.EmptyIfNil(issue.DueDate)},
		{Label: "URL", Value: issue.URL},
		{Label: "Description", Value: issue.Description},
	}

	return factory.Formatter.PrintDetail(fields, issue)
}
