package cli

import (
	"fmt"
	"os"

	"github.com/faridtriwicaksono/forgebe/internal/check"
	"github.com/spf13/cobra"
)

func newCheckCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "check [project-id]",
		Short: "Run pre-commit/CI health checks on a project",
		Long: `Aggregate health checks for CI/pre-commit workflows.

Checks profile freshness, repository existence, verification status,
and sync state. Returns non-zero exit code on failure.

Designed for use in scripts, pre-commit hooks, and CI pipelines.`,
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

			c, err := check.NewChecker(proj)
			if err != nil {
				return err
			}

			result := c.Check()

			if jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), result.JSON())
			} else {
				fmt.Fprint(cmd.OutOrStdout(), result.Text())
			}

			if result.Status == "FAIL" {
				os.Exit(1)
			}

			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")

	return cmd
}
