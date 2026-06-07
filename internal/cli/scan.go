package cli

import (
	"encoding/json"
	"fmt"

	"github.com/faridtriwicaksono/forgebe/internal/discovery"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan the current repository and infer a project profile",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "."
			if len(args) == 1 {
				target = args[0]
			}

			scanner := discovery.NewScanner()
			report, err := scanner.Scan(target)
			if err != nil {
				return err
			}

			if OutputJSON(cmd) {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(report)
			}

			data, err := yaml.Marshal(report)
			if err != nil {
				return fmt.Errorf("scan: marshal report: %w", err)
			}
			_, err = cmd.OutOrStdout().Write(data)
			return err
		},
		SilenceUsage:  true,
		SilenceErrors: false,
	}
}
