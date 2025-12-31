package auth

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/config"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdStatus creates the auth status command
func NewCmdStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "View authentication status",
		Long:  "Verify your Linear API key and display account information.",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			return runStatus(jsonOutput)
		},
	}

	return cmd
}

func runStatus(jsonOutput bool) error {
	factory, err := cmdutil.NewFactory(jsonOutput)
	if err != nil {
		if err == config.ErrNoAPIKey {
			fmt.Println("Not authenticated.")
			fmt.Println("")
			fmt.Println("To authenticate, set your Linear API key:")
			fmt.Println("  export LINEAR_API_KEY=your_api_key")
			fmt.Println("")
			fmt.Println("You can create an API key at:")
			fmt.Println("  Settings > Account > Security & Access > Personal API keys")
			return nil
		}
		return err
	}

	ctx := context.Background()

	// Get current user
	user, err := factory.Client.GetViewer(ctx)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Get organisation info
	org, err := factory.Client.GetOrganisation(ctx)
	if err != nil {
		return fmt.Errorf("failed to get organisation: %w", err)
	}

	if jsonOutput {
		data := map[string]interface{}{
			"authenticated": true,
			"user":          user,
			"organisation":  org,
		}
		return factory.Formatter.PrintJSON(data)
	}

	fmt.Println("Authenticated!")
	fmt.Println("")
	fmt.Printf("User:         %s (%s)\n", user.Name, user.Email)
	fmt.Printf("Organisation: %s\n", org.Name)
	fmt.Printf("Users:        %d\n", org.UserCount)

	if user.Admin {
		fmt.Println("Role:         Admin")
	}

	return nil
}
