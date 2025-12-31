package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stustirling/lnr/internal/cmd/auth"
	"github.com/stustirling/lnr/internal/cmd/cycle"
	"github.com/stustirling/lnr/internal/cmd/initiative"
	"github.com/stustirling/lnr/internal/cmd/issue"
	"github.com/stustirling/lnr/internal/cmd/label"
	"github.com/stustirling/lnr/internal/cmd/project"
	"github.com/stustirling/lnr/internal/cmd/state"
	"github.com/stustirling/lnr/internal/cmd/team"
	"github.com/stustirling/lnr/internal/cmd/user"
)

var rootCmd = &cobra.Command{
	Use:   "lnr",
	Short: "Linear CLI - A command-line interface for Linear",
	Long: `lnr is a command-line tool for interacting with Linear.

It provides read-only access to your Linear workspace, allowing you to
view issues, projects, initiatives, teams, and more from the terminal.

To get started, set your Linear API key:
  export LINEAR_API_KEY=your_api_key

Then verify your authentication:
  lnr auth status`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")

	// Add commands
	rootCmd.AddCommand(auth.NewCmdAuth())
	rootCmd.AddCommand(cycle.NewCmdCycle())
	rootCmd.AddCommand(initiative.NewCmdInitiative())
	rootCmd.AddCommand(issue.NewCmdIssue())
	rootCmd.AddCommand(label.NewCmdLabel())
	rootCmd.AddCommand(project.NewCmdProject())
	rootCmd.AddCommand(state.NewCmdState())
	rootCmd.AddCommand(team.NewCmdTeam())
	rootCmd.AddCommand(user.NewCmdUser())
}
