package cli

import (
	"fmt"
	"path/filepath"

	"github.com/faridtriwicaksono/forgebe/internal/adapters"
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
	cmd.AddCommand(newExportAdapterCmd())
	cmd.AddCommand(newExportBundleCmd())
	return cmd
}

func newExportSummaryCmd() *cobra.Command {
	cmd := &cobra.Command{
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

			if OutputJSON(cmd) {
				return WriteOutput(cmd, "", map[string]string{"summary": summary, "project_id": projectID})
			}

			outPath := OutputPath(cmd)
			if outPath != "" {
				wrote, err := WriteOutputFile(cmd, summary, outPath)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Written to: %s\n", wrote)
				return nil
			}

			fmt.Fprint(cmd.OutOrStdout(), summary)
			return nil
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringP("output", "o", "", "Write output to file instead of stdout")
	return cmd
}

func newExportAdapterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adapter <tool> [project-id]",
		Short: "Export a tool-specific adapter file",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tool := args[0]
			projectID, err := resolveProjectIDArg(args[1:])
			if err != nil {
				return fmt.Errorf("export adapter: %w", err)
			}
			p, err := loadProfileByID(projectID)
			if err != nil {
				return fmt.Errorf("export adapter: %w", err)
			}
			out, filename, err := adapters.Render(tool, *p)
			if err != nil {
				return err
			}

			paths, err := storage.NewPaths()
			if err != nil {
				return err
			}
			if err := paths.EnsureBaseDirs(); err != nil {
				return err
			}

			// Use --output if provided, else default exports dir
			outPath := OutputPath(cmd)
			if outPath == "" {
				outPath = filepath.Join(paths.ExportsDir(), projectID+"-"+sanitizeToolName(tool)+"-"+filepath.Base(filename))
			}

			if err := storage.AtomicWrite(outPath, []byte(out), 0600); err != nil {
				return fmt.Errorf("export adapter: write file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Exported: %s\n", outPath)
			return nil
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringP("output", "o", "", "Write output to file instead of default exports dir")
	return cmd
}

func sanitizeToolName(v string) string {
	out := ""
	for _, r := range v {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			out += string(r)
		}
	}
	if out == "" {
		return "tool"
	}
	return out
}

func newExportBundleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bundle [project-id]",
		Short: "Export a portable .forgebe.zip bundle",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := storage.NewPaths()
			if err != nil {
				return fmt.Errorf("export bundle: %w", err)
			}

			var projectID string
			if len(args) == 1 {
				projectID = args[0]
			} else {
				entries, err := storage.ListProjectIDs(paths)
				if err != nil {
					return fmt.Errorf("export bundle: list projects: %w", err)
				}
				if len(entries) == 0 {
					return fmt.Errorf("export bundle: no projects found")
				}
				projectID = entries[len(entries)-1]
			}

			outPath := OutputPath(cmd)
			if outPath == "" {
				outPath = export.DefaultBundlePath(projectID, paths)
			}

			if err := export.ExportBundle(projectID, paths, outPath); err != nil {
				return fmt.Errorf("export bundle: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Exported bundle: %s\n", outPath)
			return nil
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringP("output", "o", "", "Write bundle to specified path")
	return cmd
}
