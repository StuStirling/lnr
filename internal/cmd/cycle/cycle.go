package cycle

import (
	"github.com/spf13/cobra"
)

// NewCmdCycle creates the cycle parent command
func NewCmdCycle() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cycle",
		Short: "View cycles",
		Long:  "Commands for viewing Linear cycles (sprints).",
	}

	cmd.AddCommand(NewCmdList())
	cmd.AddCommand(NewCmdActive())
	cmd.AddCommand(NewCmdView())

	return cmd
}
