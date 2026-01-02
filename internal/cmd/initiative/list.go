package initiative

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stustirling/lnr/internal/output"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

// NewCmdList creates the initiative list command
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List initiatives",
		Long:  "List all initiatives in the Linear workspace.",
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
	initiatives, err := factory.Client.GetInitiatives(ctx)
	if err != nil {
		return fmt.Errorf("failed to list initiatives: %w", err)
	}

	// Warn if results might be truncated (API limit is 50)
	output.WarnIfTruncated(len(initiatives), 50)

	headers := []string{"NAME", "OWNER", "TARGET DATE"}
	rows := make([][]string, len(initiatives))
	for i, init := range initiatives {
		ownerName := "-"
		if init.Owner != nil {
			ownerName = init.Owner.Name
		}
		rows[i] = []string{
			output.Truncate(init.Name, 50),
			ownerName,
			output.EmptyIfNil(init.TargetDate),
		}
	}

	return factory.Formatter.Print(headers, rows, initiatives)
}
