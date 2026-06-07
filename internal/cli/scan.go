package cli

import (
	"fmt"
	"os"

	"github.com/faridtriwicaksono/forgebe/internal/discovery"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
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

			data, err := yaml.Marshal(report)
			if err != nil {
				return fmt.Errorf("scan: marshal report: %w", err)
			}
			_, err = cmd.OutOrStdout().Write(data)
			if err != nil {
				return fmt.Errorf("scan: write output: %w", err)
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: false,
	}
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
