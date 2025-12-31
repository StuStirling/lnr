package initiative

import (
	"github.com/spf13/cobra"
)

// NewCmdInitiative creates the initiative parent command
func NewCmdInitiative() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "initiative",
		Short: "Manage initiatives",
		Long:  "Commands for viewing Linear initiatives.",
	}

	cmd.AddCommand(NewCmdList())
	cmd.AddCommand(NewCmdView())

	return cmd
}
