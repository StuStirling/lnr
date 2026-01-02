package issue

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/api"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdSearch creates the issue search command
func NewCmdSearch() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search issues",
		Long:  "Search for issues matching the given query.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			opts := api.IssueListOptions{First: limit}
			return runSearch(jsonOutput, args[0], opts)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of issues to return")

	return cmd
}

func runSearch(jsonOutput bool, query string, opts api.IssueListOptions) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	issues, err := factory.Client.SearchIssues(ctx, query, opts)
	if err != nil {
		return fmt.Errorf("failed to search issues: %w", err)
	}

	// Warn if results might be truncated
	output.WarnIfTruncated(len(issues), opts.First)

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
