package team

import (
	"github.com/spf13/cobra"
)

// NewCmdTeam creates the team parent command
func NewCmdTeam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "team",
		Short: "Manage teams",
		Long:  "Commands for viewing Linear teams.",
	}

	cmd.AddCommand(NewCmdList())

	return cmd
}
