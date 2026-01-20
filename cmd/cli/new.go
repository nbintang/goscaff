package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nbintang/goscaff/internal/scaffold"
	"github.com/spf13/cobra"
)

var (
	flagModule string
	flagDB     string
	flagPreset string
)

func isInteractive(cmd *cobra.Command) bool {
	// interaktif kalau user TIDAK set flag ini
	// (yang penting: beda antara default vs "explicitly set")
	return !cmd.Flags().Changed("preset") &&
		!cmd.Flags().Changed("db") &&
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

		preset := flagPreset
		db := flagDB
		modulePath := flagModule

		if isInteractive(cmd) {
			w := scaffold.NewWizard()
			ctx := cmd.Context()

			_, err := w.SelectOption(
				ctx,
				"Preset",
				[]scaffold.Option{
					{Label: "Base (minimal, default)", Value: "base"},
					{Label: "Full (production-ready)", Value: "full"},
				},
				"base",
			)
			if err != nil {
				return err
			}

			db, err = w.SelectOption(
				ctx,
				"Database",
				[]scaffold.Option{
					{Label: "Postgres (default)", Value: "postgres"},
					{Label: "MySQL", Value: "mysql"},
				},
				"postgres",
			)
			if err != nil {
				return err
			}

			modulePath, err = w.Input(ctx, "Module path", projectName)
			if err != nil {
				return err
			}
		}

		if modulePath == "" {
			modulePath = projectName
		}

		if preset != "base" && preset != "full" {
			return fmt.Errorf("invalid --preset=%s (use base|full)", preset)
		}
		if db != "postgres" && db != "mysql" {
			return fmt.Errorf("invalid --db=%s (use postgres|mysql)", db)
		}

		opts := scaffold.ScaffoldOptions{
			ProjectName: projectName,
			ModulePath:  modulePath,
			DB:          scaffold.DBType(db),
			Preset:      scaffold.PresetType(preset),
			OutDir:      outDir,
		}

		renderer := scaffold.NewRenderer()
		s := scaffold.NewScaffold(opts, renderer)
		if err := s.Generate(); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		scaffold.NewPrinter(opts).PrintNextSteps()
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVar(&flagModule, "module", "", "Go module path (default: project-name)")
	newCmd.Flags().StringVar(&flagPreset, "preset", "base", "Template preset: base|full")
	newCmd.Flags().StringVar(&flagDB, "db", "postgres", "Database driver: postgres|mysql")
}
