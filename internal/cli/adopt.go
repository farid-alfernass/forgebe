package cli

import (
	"fmt"
	"os"

	"github.com/faridtriwicaksono/forgebe/internal/adopt"
	"github.com/spf13/cobra"
)

func newAdoptCmd() *cobra.Command {
	var (
		force      bool
		dryRunMode bool
	)

	cmd := &cobra.Command{
		Use:   "adopt [path]",
		Short: "Adopt an existing project into ForgeBE",
		Long: `Adopt an existing project directory into ForgeBE workflow.

Scans the project, detects stack/architecture/testing setup, and creates
a ForgeBE profile automatically — no full interactive onboarding required.

Use --dry-run to preview what would be created without writing anything.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath := "."
			if len(args) > 0 {
				repoPath = args[0]
			}

			// Resolve to absolute path
			if repoPath == "." {
				wd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("cannot resolve current directory: %w", err)
				}
				repoPath = wd
			}

			result, err := adopt.Run(adopt.Options{
				RepoPath: repoPath,
				DryRun:   dryRunMode,
				Force:    force,
			})
			if err != nil {
				return err
			}

			if dryRunMode {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry-run: adopt %s\n", result.RepoPath)
				fmt.Fprintf(cmd.OutOrStdout(), "  Project ID: %s\n", result.ProjectID)
				fmt.Fprintf(cmd.OutOrStdout(), "  Language:   %s\n", result.Language)
				fmt.Fprintln(cmd.OutOrStdout(), "\nActions that would be performed:")
				for _, a := range result.Actions {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", a)
				}
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Project adopted successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Project ID: %s\n", result.ProjectID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Profile:    %s\n", result.ProfilePath)
			fmt.Fprintf(cmd.OutOrStdout(), "  Language:   %s\n", result.Language)
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing profile")
	cmd.Flags().BoolVar(&dryRunMode, "dry-run", false, "Preview adopt without writing")

	return cmd
}
