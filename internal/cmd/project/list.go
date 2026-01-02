package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/api"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdList creates the project list command
func NewCmdList() *cobra.Command {
	var teamID string
	var state string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		Long:  "List all projects in the Linear workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			opts := api.ProjectListOptions{}
			if teamID != "" {
				opts.TeamID = &teamID
			}
			if state != "" {
				opts.State = &state
			}
			return runList(jsonOutput, opts)
		},
	}

	cmd.Flags().StringVar(&teamID, "team", "", "Filter by team ID")
	cmd.Flags().StringVar(&state, "state", "", "Filter by state (e.g., started, completed, canceled)")

	return cmd
}

func runList(jsonOutput bool, opts api.ProjectListOptions) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	projects, err := factory.Client.GetProjects(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	// Warn if results might be truncated (API default limit is 50)
	limit := opts.First
	if limit == 0 {
		limit = 50
	}
	output.WarnIfTruncated(len(projects), limit)

	headers := []string{"NAME", "STATE", "PROGRESS", "LEAD", "TEAMS"}
	rows := make([][]string, len(projects))
	for i, p := range projects {
		leadName := "-"
		if p.Lead != nil {
			leadName = p.Lead.Name
		}
		teamKeys := make([]string, len(p.Teams))
		for j, t := range p.Teams {
			teamKeys[j] = t.Key
		}
		rows[i] = []string{
			output.Truncate(p.Name, 40),
			p.State,
			output.FormatPercentage(p.Progress),
			leadName,
			strings.Join(teamKeys, ", "),
		}
	}

	return factory.Formatter.Print(headers, rows, projects)
}
