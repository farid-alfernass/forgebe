package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"github.com/spf13/cobra"
)

const doctorBanner = "ForgeBE Doctor — Project Health Check"

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor [project-id]",
		Short: "Validate ForgeBE local state and project profile",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), doctorBanner)

			paths, err := storage.NewPaths()
			if err != nil {
				return fmt.Errorf("doctor: %w", err)
			}

			forgebeOk := checkDir(cmd, "ForgeBE root", paths.Root())
			projectsOk := checkDir(cmd, "Projects store", paths.ProjectsDir())
			if !forgebeOk || !projectsOk {
				fmt.Fprintln(cmd.OutOrStdout(), "\n  Run 'forgebe init' to set up ForgeBE for this project.")
				return nil
			}

			if len(args) == 1 {
				return validateSingle(cmd, paths, args[0])
			}
			return validateDefaultOrList(cmd, paths)
		},
		SilenceUsage: true,
	}
}

func checkDir(cmd *cobra.Command, label, path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "  ❌ %s: %v\n", label, err)
		return false
	}
	if info.IsDir() {
		fmt.Fprintf(cmd.OutOrStdout(), "  ✅ %s: %s\n", label, path)
		return true
	}
	fmt.Fprintf(cmd.OutOrStdout(), "  ❌ %s: not a directory\n", label)
	return false
}

func validateSingle(cmd *cobra.Command, paths *storage.Paths, projectID string) error {
	store := profile.NewStore(paths)
	p, err := store.LoadProfile(projectID)
	if err != nil {
		return fmt.Errorf("doctor: load profile: %w", err)
	}

	result := profile.ValidateProjectProfile(p)
	fmt.Fprintf(cmd.OutOrStdout(), "\n  Profile: %s\n", projectID)
	fmt.Fprintf(cmd.OutOrStdout(), "  Project: %s (%s)\n", p.Project.Name, p.Stack.PrimaryLanguage)
	fmt.Fprintf(cmd.OutOrStdout(), "  Status: %s\n\n", result.Summary)

	icons := map[string]string{"error": "❌", "warning": "⚠️", "info": "ℹ️"}
	for _, issue := range result.Issues {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s [%s] %s\n", icons[issue.Severity], issue.Field, issue.Message)
	}
	return nil
}

func validateDefaultOrList(cmd *cobra.Command, paths *storage.Paths) error {
	entries, err := os.ReadDir(paths.ProjectsDir())
	if err != nil {
		return fmt.Errorf("doctor: list projects: %w", err)
	}
	var projectDirs []string
	for _, e := range entries {
		if e.IsDir() {
			projectDirs = append(projectDirs, e.Name())
		}
	}
	sort.Strings(projectDirs)

	if len(projectDirs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "  No projects found.")
		fmt.Fprintln(cmd.OutOrStdout(), "  Run 'forgebe init' to initialize a project.")
		return nil
	}
	if len(projectDirs) == 1 {
		return validateSingle(cmd, paths, projectDirs[0])
	}

	fmt.Fprintln(cmd.OutOrStdout(), "\n  Multiple projects found. Specify one:")
	for _, dir := range projectDirs {
		store := profile.NewStore(paths)
		p, err := store.LoadProfile(dir)
		if err == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "    %s  (%s)\n", dir, p.Project.Name)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "    %s\n", dir)
		}
	}
	return nil
}
