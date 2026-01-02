package project

import (
	"github.com/spf13/cobra"
)

// NewCmdProject creates the project parent command
func NewCmdProject() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "View projects",
		Long:  "Commands for viewing Linear projects.",
	}

	cmd.AddCommand(NewCmdList())
	cmd.AddCommand(NewCmdView())

	return cmd
}
