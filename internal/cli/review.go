package cli

import (
	"fmt"
	"os"

	"github.com/faridtriwicaksono/forgebe/internal/git"
	"github.com/faridtriwicaksono/forgebe/internal/review"
	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {
	var (
		jsonOutput bool
		staged     bool
		since      string
		strict     bool
	)

	cmd := &cobra.Command{
		Use:   "review [project-id]",
		Short: "Review code changes against your policy (awareness report)",
		Long: `Compare the git diff (what your AI produced) against this project's
ForgeBE policy and print an awareness report: what changed, and what needs
your conscious review.

Informational by default (exit 0). Use --strict to exit non-zero on any
FAIL finding, for pre-commit hooks and CI. Works with any AI tool — it
inspects the diff, not the agent.`,
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

			repoPath := proj.Metadata.RepoPath
			if repoPath == "" {
				repoPath, _ = os.Getwd()
			}
			if !git.IsRepo(repoPath) {
				return fmt.Errorf("review: %s is not a git repository", repoPath)
			}

			spec := git.RangeSpec{Staged: staged, Since: since}
			changes, err := git.ChangedFiles(repoPath, spec)
			if err != nil {
				return fmt.Errorf("review: %w", err)
			}

			reviewer, err := review.NewReviewer(proj, repoPath, changes)
			if err != nil {
				return err
			}

			report := review.NewReport(reviewer.Run(), rangeLabel(spec))
			if jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), report.JSON())
			} else {
				fmt.Fprint(cmd.OutOrStdout(), report.Text())
			}

			if strict && report.HasFail() {
				os.Exit(1)
			}
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	cmd.Flags().BoolVar(&staged, "staged", false, "Review staged changes (git diff --cached)")
	cmd.Flags().StringVar(&since, "since", "", "Review changes since a git ref (e.g. main)")
	cmd.Flags().BoolVar(&strict, "strict", false, "Exit non-zero if any FAIL finding")

	return cmd
}

func rangeLabel(spec git.RangeSpec) string {
	switch {
	case spec.Staged:
		return "staged (vs HEAD)"
	case spec.Since != "":
		return "since " + spec.Since
	default:
		return "working tree (vs HEAD)"
	}
}
