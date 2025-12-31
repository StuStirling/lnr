package user

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdMe creates the user me command
func NewCmdMe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "View current user",
		Long:  "Display information about the currently authenticated user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runMe(jsonOutput)
		},
	}

	return cmd
}

func runMe(jsonOutput bool) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	user, err := factory.Client.GetViewer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	fields := []output.DetailField{
		{Label: "ID", Value: user.ID},
		{Label: "Name", Value: user.Name},
		{Label: "Email", Value: user.Email},
		{Label: "Display Name", Value: user.DisplayName},
		{Label: "Active", Value: fmt.Sprintf("%t", user.Active)},
		{Label: "Admin", Value: fmt.Sprintf("%t", user.Admin)},
	}

	return factory.Formatter.PrintDetail(fields, user)
}
