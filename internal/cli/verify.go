package cli

import (
	"fmt"

	"github.com/faridtriwicaksono/forgebe/internal/verify"
	"github.com/spf13/cobra"
)

func newVerifyCmd() *cobra.Command {
	var (
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "verify [project-id]",
		Short: "Verify project profile against engineering policies",
		Long: `Check that a ForgeBE project profile is properly configured and aligned
with engineering policies.

Validates testing strategy, dependency management, sensitive areas,
change policies, architecture patterns, and more.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectID, err := resolveProjectIDArg(args)
			if err != nil {
				return fmt.Errorf("cannot resolve project: %w", err)
			}

			proj, err := loadProfileByID(projectID)
			if err != nil {
				return fmt.Errorf("project %q not found: %w", projectID, err)
			}

			v, err := verify.NewValidator(proj)
			if err != nil {
				return err
			}

			results := v.Validate()

			if jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), verify.ResultsToJSON(results))
			} else {
				fmt.Fprint(cmd.OutOrStdout(), verify.ResultsToText(results))
			}

			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")

	return cmd
}
