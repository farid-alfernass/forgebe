package cli

import (
	"fmt"

	"github.com/faridtriwicaksono/forgebe/internal/scaffold"
	"github.com/spf13/cobra"
)

func newScaffoldCmd() *cobra.Command {
	var (
		outputDir     string
		templateType  string
		force         bool
		dryRunMode    bool
		listTemplates bool
	)

	cmd := &cobra.Command{
		Use:   "scaffold [project-id]",
		Short: "Generate boilerplate code structure based on profile or template",
		Long: `Generate directory structure and starter files for a project.
By default, it infers the template from the project profile, but you can
explicitly specify a template using --template.

Supported templates:
  go-service      Go service with layered architecture
  go-api          Go API with handlers, models, repository
  node-express    Node.js + TypeScript + Express
  node-nestjs     Node.js + TypeScript + NestJS
  python-fastapi  Python + FastAPI
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if listTemplates {
				fmt.Fprintln(cmd.OutOrStdout(), "Available templates:")
				for _, t := range scaffold.Templates() {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", t)
				}
				return nil
			}

			projectID, err := resolveProjectIDArg(args)
			if err != nil {
				return fmt.Errorf("cannot resolve project: %w", err)
			}

			proj, err := loadProfileByID(projectID)
			if err != nil {
				return fmt.Errorf("project %q not found: %w", projectID, err)
			}

			gen := scaffold.NewGenerator(proj, outputDir, scaffold.TemplateType(templateType), force, dryRunMode)
			files, err := gen.Generate()
			if err != nil {
				return err
			}

			if dryRunMode {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry-run: %d file(s) to generate in %s\n", len(files), outputDir)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Scaffold complete: %d file(s) generated in %s\n", len(files), outputDir)
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for generated files")
	cmd.Flags().StringVarP(&templateType, "template", "t", "", "Specify template name (e.g., go-service, node-express). If omitted, infers from profile.")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing files")
	cmd.Flags().BoolVarP(&dryRunMode, "dry-run", "d", false, "Preview files to be generated without writing them")
	cmd.Flags().BoolVar(&listTemplates, "list-templates", false, "List all available templates")

	return cmd
}
