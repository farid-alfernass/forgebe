package cli

import (
	"fmt"

	"github.com/faridtriwicaksono/forgebe/internal/export"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export ForgeBE context to various formats",
	}
	cmd.AddCommand(newExportSummaryCmd())
	return cmd
}

func newExportSummaryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "summary [project-id]",
		Short: "Export a compact project summary as markdown",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := storage.NewPaths()
			if err != nil {
				return fmt.Errorf("export summary: %w", err)
			}

			var projectID string
			if len(args) == 1 {
				projectID = args[0]
			} else {
				// Find the most recent project
				entries, err := storage.ListProjectIDs(paths)
				if err != nil {
					return fmt.Errorf("export summary: list projects: %w", err)
				}
				if len(entries) == 0 {
					return fmt.Errorf("export summary: no projects found")
				}
				projectID = entries[len(entries)-1]
			}

			store := profile.NewStore(paths)
			p, err := store.LoadProfile(projectID)
			if err != nil {
				return fmt.Errorf("export summary: load profile: %w", err)
			}

			summary, err := export.RenderSummary(*p)
			if err != nil {
				return fmt.Errorf("export summary: render: %w", err)
			}
			fmt.Fprint(cmd.OutOrStdout(), summary)
			return nil
		},
		SilenceUsage: true,
	}
}
