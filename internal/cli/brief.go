package cli

import (
	"fmt"
	"os"

	"github.com/faridtriwicaksono/forgebe/internal/adapters"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"github.com/spf13/cobra"
)

func newBriefCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "brief <tool> [project-id]",
		Short: "Render a compact ready-to-paste AI brief",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tool := args[0]
			projectID, err := resolveProjectIDArg(args[1:])
			if err != nil {
				return fmt.Errorf("brief: %w", err)
			}
			p, err := loadProfileByID(projectID)
			if err != nil {
				return fmt.Errorf("brief: %w", err)
			}
			out, err := adapters.RenderBrief(tool, *p)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), out)
			return nil
		},
		SilenceUsage: true,
	}
}

func resolveProjectIDArg(rest []string) (string, error) {
	if len(rest) == 1 && rest[0] != "" {
		return rest[0], nil
	}
	wd, _ := os.Getwd()
	return profile.ProjectID(wd, ""), nil
}

func loadProfileByID(projectID string) (*profile.ProjectProfile, error) {
	paths, err := storage.NewPaths()
	if err != nil {
		return nil, err
	}
	store := profile.NewStore(paths)
	return store.LoadProfile(projectID)
}
