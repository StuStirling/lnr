package issue

import (
	"github.com/spf13/cobra"
)

// NewCmdIssue creates the issue parent command
func NewCmdIssue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Long:  "Commands for viewing Linear issues.",
	}

	cmd.AddCommand(NewCmdList())
	cmd.AddCommand(NewCmdView())
	cmd.AddCommand(NewCmdSearch())

	return cmd
}
