package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nbintang/goscaff/internal"
	"github.com/spf13/cobra"
)

var (
	flagModule   string
	flagTemplate string
)

func isInteractive(cmd *cobra.Command) bool {
	if os.Getenv("GOSCAFF_NON_INTERACTIVE") == "1" {
		return false
	}
	return !cmd.Flags().Changed("template") &&
		!cmd.Flags().Changed("module")
}

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new project from embedded templates",
	Long: `Create a new Go backend project.

This command will:
  1) Create a new directory using [project-name]
  2) Render embedded templates based on preset (base|full)
  3) Run "go mod tidy"
  4) Initialize git repository (git init)

Preset:
  base - minimal template
  full - complete template (default)

Database:
  postgres (default) or mysql
`,
	Args: cobra.ExactArgs(1),
	Example: `  
  # Quick start
  goscaff new myapp

  # Full preset (default) with module path
  goscaff new myapp --module github.com/you/myapp

  # Base preset (minimal)
  goscaff new myapp --preset base --module github.com/you/myapp

  # Full preset + MySQL
  goscaff new myapp --preset full --db mysql --module github.com/you/myapp
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectArg := args[0]
		outDir := projectArg
		if !filepath.IsAbs(projectArg) {
			outDir = filepath.Join(".", projectArg)
		}
		outDir = filepath.Clean(outDir)
		projectName := filepath.Base(outDir)

		if _, err := os.Stat(outDir); err == nil {
			return fmt.Errorf("directory %s already exists", outDir)
		}

		tpl := flagTemplate
		modulePath := flagModule

		if isInteractive(cmd) {
			w := internal.NewWizard()
			ctx := cmd.Context()

			templates, err := internal.ListTemplates()
			if err != nil {
				return err
			}
			if len(templates) == 0 {
				return fmt.Errorf("tidak ada template di embedded FS")
			}

			opts := make([]internal.Option, 0, len(templates))
			for _, id := range templates {
				opts = append(opts, internal.Option{
					Label: internal.PrettyTemplateLabel(id),
					Value: id,
				})
			}

			tpl, err = w.SelectOption(ctx, "Template", opts, templates[0])
			if err != nil {
				return err
			}

			modulePath, err = w.Input(ctx, "Module path", projectName)
			if err != nil {
				return err
			}
		}

		if tpl == "" {
			return fmt.Errorf("template wajib diisi. Lihat pilihan via `goscaff templates` atau pakai mode interaktif")
		}

		if modulePath == "" {
			modulePath = projectName
		}

		opts := internal.ScaffoldOptions{
			ProjectName: projectName,
			ModulePath:  modulePath,
			Template:    tpl,
			OutDir:      outDir,
		}

		renderer := internal.NewRenderer()
		s := internal.NewScaffold(opts, renderer)

		if err := s.Generate(); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		internal.NewPrinter(opts).PrintNextSteps()
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVar(&flagModule, "module", "", "Go module path (default: project-name)")
	newCmd.Flags().StringVar(&flagTemplate, "template", "", "Template ID to use (skips interactive template selection)")
}
