package auth

import (
	"github.com/spf13/cobra"
)

// NewCmdAuth creates the auth parent command
func NewCmdAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
		Long:  "Commands for managing Linear authentication.",
	}

	cmd.AddCommand(NewCmdStatus())

	return cmd
}
