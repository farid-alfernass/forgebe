package cli

import (
	"fmt"
	"os"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"github.com/spf13/cobra"
)

func newProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage ForgeBE project profiles",
	}
	cmd.AddCommand(newProfileShowCmd())
	return cmd
}

func newProfileShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [project-id]",
		Short: "Show the current project profile",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := storage.NewPaths()
			if err != nil {
				return fmt.Errorf("profile show: %w", err)
			}

			projectID := ""
			if len(args) == 1 {
				projectID = args[0]
			} else {
				wd, _ := os.Getwd()
				projectID = profile.ProjectID(wd, "")
			}

			if projectID == "" {
				return fmt.Errorf("profile show: specify project id or run from a ForgeBE-initialized directory")
			}

			store := profile.NewStore(paths)
			p, err := store.LoadProfile(projectID)
			if err != nil {
				return fmt.Errorf("profile show: profile not found for %s: %w", projectID, err)
			}

			if OutputJSON(cmd) {
				return WriteOutput(cmd, "", p)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Project:    %s\n", p.Project.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Profile ID: %s\n", p.Metadata.ProfileID)
			fmt.Fprintf(cmd.OutOrStdout(), "Language:   %s\n", p.Stack.PrimaryLanguage)
			fmt.Fprintf(cmd.OutOrStdout(), "Framework:  %s\n", p.Stack.Framework)
			fmt.Fprintf(cmd.OutOrStdout(), "Arch:       %s\n", p.Stack.Architecture)
			fmt.Fprintf(cmd.OutOrStdout(), "Maturity:   %s\n", p.Project.Maturity)
			fmt.Fprintf(cmd.OutOrStdout(), "Type:       %s\n", p.Project.Type)
			fmt.Fprintf(cmd.OutOrStdout(), "Testing:    %s\n", p.Policy.Testing.Strategy)
			fmt.Fprintf(cmd.OutOrStdout(), "Delivery:   %s (%s priority)\n", p.Policy.Delivery.Mode, p.Policy.Delivery.Priority)
			fmt.Fprintf(cmd.OutOrStdout(), "Repo:       %s\n", p.Metadata.RepoPath)
			fmt.Fprintf(cmd.OutOrStdout(), "Sensitive:  %v\n", p.Areas.SensitiveAreas)
			return nil
		},
		SilenceUsage: true,
	}
}
