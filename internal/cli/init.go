package cli

import (
	"fmt"
	"path/filepath"

	"github.com/faridtriwicaksono/forgebe/internal/contract"
	"github.com/faridtriwicaksono/forgebe/internal/discovery"
	"github.com/faridtriwicaksono/forgebe/internal/onboarding"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var nonInteractive bool

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize ForgeBE for the current project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "."
			if len(args) == 1 {
				target = args[0]
			}

			absTarget, err := filepath.Abs(target)
			if err != nil {
				return fmt.Errorf("init: resolve target path: %w", err)
			}

			scanner := discovery.NewScanner()
			report, err := scanner.Scan(absTarget)
			if err != nil {
				return fmt.Errorf("init: discovery scan failed: %w", err)
			}

			projectID := profile.ProjectID(absTarget, "")
			repoName := filepath.Base(absTarget)

			var answers *onboarding.Answers
			if nonInteractive {
				answers = onboarding.DefaultAnswersFromDiscovery(repoName, report)
			} else {
				answers = onboarding.DefaultAnswersFromDiscovery(repoName, report)
				guided, err := onboarding.RunGuided()
				if err != nil {
					return err
				}
				answers = guided
			}

			profileData := onboarding.BuildProfile(projectID, absTarget, answers, report)

			paths, err := storage.NewPaths()
			if err != nil {
				return fmt.Errorf("init: get home dir: %w", err)
			}
			if err := paths.EnsureBaseDirs(); err != nil {
				return fmt.Errorf("init: ensure forgebe directories: %w", err)
			}

			store := profile.NewStore(paths)
			if err := store.SaveProfile(profileData); err != nil {
				return fmt.Errorf("init: save profile: %w", err)
			}
			if err := store.SaveDiscovery(projectID, *report); err != nil {
				return fmt.Errorf("init: save discovery: %w", err)
			}

			contractMD, err := contract.Render(profileData)
			if err != nil {
				return fmt.Errorf("init: render contract: %w", err)
			}
			if err := storage.AtomicWrite(paths.ProjectContractPath(projectID), []byte(contractMD), 0600); err != nil {
				return fmt.Errorf("init: save contract: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ForgeBE initialized successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "Project ID: %s\n", projectID)
			fmt.Fprintf(cmd.OutOrStdout(), "Profile: %s\n", paths.ProjectProfilePath(projectID))
			fmt.Fprintf(cmd.OutOrStdout(), "Contract: %s\n", paths.ProjectContractPath(projectID))
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "Use auto-discovery defaults without prompts")
	return cmd
}
