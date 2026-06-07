package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of ForgeBE
	version = "dev"
	// Commit is the git commit hash
	commit = "none"
	// BuildDate is the date the binary was built
	buildDate = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print ForgeBE version",
		Run: func(cmd *cobra.Command, args []string) {
			if OutputJSON(cmd) {
				WriteOutput(cmd, "", map[string]string{
					"version":    version,
					"commit":     commit,
					"build_date": buildDate,
				})
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ForgeBE %s\n", version)
			fmt.Fprintf(cmd.OutOrStdout(), "Commit: %s\n", commit)
			fmt.Fprintf(cmd.OutOrStdout(), "Built:  %s\n", buildDate)
		},
	}
}
