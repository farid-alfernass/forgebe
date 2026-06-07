package cli

import (
	"fmt"
	"path/filepath"

	"github.com/faridtriwicaksono/forgebe/internal/export"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import <bundle-path> [target-dir]",
		Short: "Import a .forgebe.zip bundle into local storage or a target directory",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			bundlePath := args[0]

			paths, err := storage.NewPaths()
			if err != nil {
				return fmt.Errorf("import: %w", err)
			}

			absBundle, err := filepath.Abs(bundlePath)
			if err != nil {
				return fmt.Errorf("import: resolve bundle path: %w", err)
			}

			targetDir := paths.ProjectsDir()
			if len(args) == 2 {
				targetDir = args[1]
			}

			if err := export.ImportBundle(absBundle, targetDir); err != nil {
				return fmt.Errorf("import: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Imported bundle: %s\n", absBundle)
			fmt.Fprintf(cmd.OutOrStdout(), "Destination: %s\n", targetDir)
			return nil
		},
		SilenceUsage: true,
	}
}
