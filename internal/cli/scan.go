package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "Scan the current repository and infer a project profile",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "ForgeBE scan is not implemented yet. Planned in Day 2 discovery.")
		},
	}
}
