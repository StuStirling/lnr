package user

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdList creates the user list command
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Long:  "List all users in the Linear workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runList(jsonOutput)
		},
	}

	return cmd
}

func runList(jsonOutput bool) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		return err
	}

	ctx := context.Background()
	users, err := factory.Client.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	// Warn if results might be truncated (API limit is 100)
	output.WarnIfTruncated(len(users), 100)

	headers := []string{"NAME", "EMAIL", "ACTIVE", "ADMIN"}
	rows := make([][]string, len(users))
	for i, u := range users {
		active := "Yes"
		if !u.Active {
			active = "No"
		}
		admin := ""
		if u.Admin {
			admin = "Yes"
		}
		rows[i] = []string{u.Name, u.Email, active, admin}
	}

	return factory.Formatter.Print(headers, rows, users)
}
