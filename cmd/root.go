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

// Version is set by goreleaser via ldflags
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "lnr",
	Short:   "Linear CLI - A command-line interface for Linear",
	Version: Version,
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

	// Add completion command (built into Cobra)
	rootCmd.AddCommand(newCompletionCmd())
}

func newCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for lnr.

To load completions:

Bash:
  $ source <(lnr completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ lnr completion bash > /etc/bash_completion.d/lnr
  # macOS:
  $ lnr completion bash > $(brew --prefix)/etc/bash_completion.d/lnr

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ lnr completion zsh > "${fpath[1]}/_lnr"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ lnr completion fish | source

  # To load completions for each session, execute once:
  $ lnr completion fish > ~/.config/fish/completions/lnr.fish

PowerShell:
  PS> lnr completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> lnr completion powershell > lnr.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			}
			return nil
		},
	}
}
