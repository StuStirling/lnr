package state

import (
	"github.com/spf13/cobra"
)

// NewCmdState creates the state parent command
func NewCmdState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "View workflow states",
		Long:  "Commands for viewing Linear workflow states.",
	}

	cmd.AddCommand(NewCmdList())

	return cmd
}
