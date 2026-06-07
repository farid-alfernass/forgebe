package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"github.com/spf13/cobra"
)

const doctorBanner = "ForgeBE Doctor — Project Health Check\n"

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor [project-id]",
		Short: "Validate the ForgeBE setup and project profile",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := storage.NewPaths()
			if err != nil {
				return err
			}
			if err := paths.EnsureBaseDirs(); err != nil {
				return fmt.Errorf("doctor: initialize directories: %w", err)
			}

			// Collect issues
			type checkItem struct {
				Check   string `json:"check"`
				Status  string `json:"status"`
				Message string `json:"message,omitempty"`
			}
			var checks []checkItem

			// Check root dir
			rootDir := paths.Root()
			if fi, err := os.Stat(rootDir); err == nil && fi.IsDir() {
				checks = append(checks, checkItem{Check: "ForgeBE root directory", Status: "ok", Message: rootDir})
			} else {
				checks = append(checks, checkItem{Check: "ForgeBE root directory", Status: "error", Message: "not found or inaccessible"})
			}

			// Check projects dir
			projDir := paths.ProjectsDir()
			if fi, err := os.Stat(projDir); err == nil && fi.IsDir() {
				checks = append(checks, checkItem{Check: "Projects store", Status: "ok", Message: projDir})
			} else {
				checks = append(checks, checkItem{Check: "Projects store", Status: "error", Message: "not found"})
			}

			// Resolve project ID
			projectID := ""
			if len(args) == 1 {
				projectID = args[0]
			} else {
				// Look for an existing project
				ids, _ := storage.ListProjectIDs(paths)
				if len(ids) > 0 {
					projectID = ids[0]
				}
			}

			if projectID == "" || !dirExists(paths.ProjectDir(projectID)) {
				if OutputJSON(cmd) {
					result := map[string]interface{}{
						"valid":  false,
						"checks": checks,
					}
					return WriteOutput(cmd, "", result)
				}
				fmt.Fprint(cmd.OutOrStdout(), doctorBanner)
				for _, c := range checks {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s %s\n", statusIcon(c.Status), c.Check)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "\n  No project profile found. Run `forgebe init` first.")
				return nil
			}

			// Validate profile
			store := profile.NewStore(paths)
			p, err := store.LoadProfile(projectID)
			if err != nil {
				if OutputJSON(cmd) {
					result := map[string]interface{}{
						"valid":   false,
						"project": projectID,
						"error":   err.Error(),
						"checks":  checks,
					}
					return WriteOutput(cmd, "", result)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\n  Profile: %s\n  Error loading profile: %v\n", projectID, err)
				return nil
			}

			result := profile.ValidateProjectProfile(p)

			if OutputJSON(cmd) {
				jsonResult := map[string]interface{}{
					"valid":    result.Valid,
					"project":  projectID,
					"name":     p.Project.Name,
					"language": p.Stack.PrimaryLanguage,
					"checks":   checks,
					"issues":   result.Issues,
				}
				return WriteOutput(cmd, "", jsonResult)
			}

			// Text output
			fmt.Fprint(cmd.OutOrStdout(), doctorBanner)
			for _, c := range checks {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s %s\n", statusIcon(c.Status), c.Check)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n  Profile: %s\n  Project: %s (%s)\n", projectID, p.Project.Name, p.Stack.PrimaryLanguage)
			if result.Valid {
				fmt.Fprintln(cmd.OutOrStdout(), "  Status: All checks passed")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "  Status: Issues found")
				sort.Slice(result.Issues, func(i, j int) bool { return result.Issues[i].Severity < result.Issues[j].Severity })
				for _, issue := range result.Issues {
					fmt.Fprintf(cmd.OutOrStdout(), "    %s [%s] %s\n", issue.Severity, issue.Field, issue.Message)
				}
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: false,
	}
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func statusIcon(s string) string {
	switch s {
	case "ok":
		return "✅"
	case "error":
		return "❌"
	case "warning":
		return "⚠️"
	default:
		return "•"
	}
}
