package cli

import (
	"fmt"

	"github.com/faridtriwicaksono/forgebe/internal/scaffold"
	"github.com/spf13/cobra"
)

var supportedTemplates = []string{
	"go-service",
	"go-api",
	"node-express",
	"node-nestjs",
	"python-fastapi",
}

func newScaffoldCmd() *cobra.Command {
	var (
		outputDir    string
		templateType string
		force        bool
		dryRunMode   bool
	)

	cmd := &cobra.Command{
		Use:   "scaffold [project-id]",
		Short: "Generate project boilerplate from profile or template",
		Long: `Generate a project structure based on an existing ForgeBE profile
or a specified template.

Templates:
  go-service      Go service with layered architecture
  go-api          Go API with handlers, models, repository
  node-express    Node.js + TypeScript + Express
  node-nestjs     Node.js + TypeScript + NestJS
  python-fastapi  Python + FastAPI
`,
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

			// Resolve template
			var tmpl scaffold.TemplateType
			if templateType != "" {
				tmpl = scaffold.TemplateType(templateType)
			}

			if outputDir == "" {
				outputDir = proj.Metadata.RepoPath
			}

			// Generate scaffolding
			gen := scaffold.NewGenerator(proj, outputDir, tmpl, force, dryRunMode)
			files, err := gen.Generate()
			if err != nil {
				return fmt.Errorf("scaffold failed: %w", err)
			}

			if dryRunMode {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry-run: %d file(s) to generate\n", len(files))
				for _, f := range files {
					if f.IsDir {
						fmt.Fprintf(cmd.OutOrStdout(), "  %s/\n", f.Path)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", f.Path)
					}
				}
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Scaffold complete: %d file(s) generated in %s\n", len(files), outputDir)
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: project repo path)")
	cmd.Flags().StringVarP(&templateType, "template", "t", "", fmt.Sprintf("Template type (%s)", joinStrings(supportedTemplates)))
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")
	cmd.Flags().BoolVar(&dryRunMode, "dry-run", false, "Preview files to generate without writing")

	return cmd
}

func joinStrings(items []string) string {
	result := ""
	for i, s := range items {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}
