package issue

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/api"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdList creates the issue list command
func NewCmdList() *cobra.Command {
	var teamID string
	var assigneeID string
	var stateID string
	var projectID string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		Long:  "List issues in the Linear workspace with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			opts := api.IssueListOptions{First: limit}
			if teamID != "" {
				opts.TeamID = &teamID
			}
			if assigneeID != "" {
				opts.AssigneeID = &assigneeID
			}
			if stateID != "" {
				opts.StateID = &stateID
			}
			if projectID != "" {
				opts.ProjectID = &projectID
			}
			return runList(jsonOutput, opts)
		},
	}

	cmd.Flags().StringVar(&teamID, "team", "", "Filter by team ID")
	cmd.Flags().StringVar(&assigneeID, "assignee", "", "Filter by assignee ID")
	cmd.Flags().StringVar(&stateID, "state", "", "Filter by state ID")
	cmd.Flags().StringVar(&projectID, "project", "", "Filter by project ID")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of issues to return")

	return cmd
}

func runList(jsonOutput bool, opts api.IssueListOptions) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	issues, err := factory.Client.GetIssues(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}

	headers := []string{"ID", "TITLE", "STATE", "ASSIGNEE", "PRIORITY"}
	rows := make([][]string, len(issues))
	for i, issue := range issues {
		stateName := "-"
		if issue.State != nil {
			stateName = issue.State.Name
		}
		assigneeName := "-"
		if issue.Assignee != nil {
			assigneeName = issue.Assignee.Name
		}
		rows[i] = []string{
			issue.Identifier,
			output.Truncate(issue.Title, 50),
			stateName,
			assigneeName,
			output.PriorityLabel(issue.Priority),
		}
	}

	return factory.Formatter.Print(headers, rows, issues)
}
