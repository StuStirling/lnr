package user

import (
	"github.com/spf13/cobra"
)

// NewCmdUser creates the user parent command
func NewCmdUser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "View users",
		Long:  "Commands for viewing Linear workspace users.",
	}

	cmd.AddCommand(NewCmdMe())
	cmd.AddCommand(NewCmdList())

	return cmd
}
