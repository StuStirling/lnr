package auth

import (
	"github.com/spf13/cobra"
)

// NewCmdAuth creates the auth parent command
func NewCmdAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
		Long:  "Commands for managing Linear authentication.",
	}

	cmd.AddCommand(NewCmdStatus())

	return cmd
}
