package cli

import (
	"fmt"

	"github.com/faridtriwicaksono/forgebe/internal/prompt"
	"github.com/spf13/cobra"
)

func newPromptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prompt <mode> [task]",
		Short: "Generate a context-aware AI prompt",
		Long: `Generate a smart, context-aware prompt for AI tools based on the
project profile. More advanced than 'brief' — task-specific and mode-driven.

Modes:
  implement  Generate code implementation prompt
  review     Generate code review prompt
  debug      Generate debugging prompt
  plan       Generate implementation planning prompt

Examples:
  forgebe prompt implement "Add user authentication endpoint"
  forgebe prompt review "Check the payment handler for security issues"
  forgebe prompt debug "Fix panic in request middleware"
  forgebe prompt plan "Migrate database from MySQL to PostgreSQL"
`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := args[0]
			task := ""
			if len(args) > 1 {
				task = args[1]
			}

			projectID, err := resolveProjectIDArg(nil)
			if err != nil {
				return fmt.Errorf("cannot resolve project: %w", err)
			}

			proj, err := loadProfileByID(projectID)
			if err != nil {
				return fmt.Errorf("project %q not found: %w", projectID, err)
			}

			g, err := prompt.NewGenerator(proj)
			if err != nil {
				return err
			}

			output, err := g.Generate(mode, task)
			if err != nil {
				return err
			}

			fmt.Fprint(cmd.OutOrStdout(), output)
			return nil
		},
		SilenceUsage: true,
	}

	return cmd
}
