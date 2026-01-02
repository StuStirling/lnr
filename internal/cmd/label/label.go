package label

import (
	"github.com/spf13/cobra"
)

// NewCmdLabel creates the label parent command
func NewCmdLabel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "View labels",
		Long:  "Commands for viewing Linear issue labels.",
	}

	cmd.AddCommand(NewCmdList())

	return cmd
}
