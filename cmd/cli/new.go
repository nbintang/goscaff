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
	flagModule       string
	flagTemplate     string
	flagFramework    string
	flagDatabase     string
	flagArchitecture string
	flagDI           string
	flagPreset       string
)

func isInteractive(cmd *cobra.Command) bool {
	if os.Getenv("GOSCAFF_NON_INTERACTIVE") == "1" {
		return false
	}
	return !cmd.Flags().Changed("template") && !hasSelectionFlags(cmd)
}

func hasSelectionFlags(cmd *cobra.Command) bool {
	return cmd.Flags().Changed("framework") ||
		cmd.Flags().Changed("db") ||
		cmd.Flags().Changed("architecture") ||
		cmd.Flags().Changed("di") ||
		cmd.Flags().Changed("preset")
}

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new project from embedded templates",
	Long: `Create a new Go backend project.

This command will:
  1) Create a new directory using [project-name]
  2) Ask for framework, database, architecture, and dependency injection
  3) Run "go mod tidy"
  4) Initialize git repository (git init)

Use the wizard for the recommended experience, or pass flags for automation.
`,
	Args: cobra.ExactArgs(1),
	Example: `  
  # Quick start
  goscaff new myapp

  # With module path
  goscaff new myapp --module github.com/you/myapp

  # Non-interactive Gin + PostgreSQL + Modular + Uber Fx
  goscaff new myapp --framework gin --db postgres --architecture modular --di uber-fx

  # Legacy preset flags remain supported
  goscaff new myapp --preset full --db mysql
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
		var spec internal.TemplateSpec

		if tpl != "" {
			var err error
			spec, err = templateSpecFromID(tpl)
			if err != nil {
				return err
			}
		} else if isInteractive(cmd) {
			w := internal.NewWizard()
			cfg, ok, err := runNewWizard(cmd.Context(), w, projectName, modulePath, selectionFromFlags())
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			spec = cfg.Template
			tpl = spec.ID
			modulePath = cfg.ModulePath
		} else {
			var err error
			spec, err = internal.ResolveTemplate(defaultSelection(selectionFromFlags()))
			if err != nil {
				return err
			}
			tpl = spec.ID
		}

		if tpl == "" {
			return fmt.Errorf("template wajib diisi. Lihat pilihan via `goscaff templates` atau pakai mode interaktif")
		}

		if modulePath == "" {
			modulePath = projectName
		}

		opts := internal.ScaffoldOptions{
			ProjectName:  projectName,
			ModulePath:   modulePath,
			Template:     tpl,
			OutDir:       outDir,
			Framework:    spec.Framework,
			Database:     spec.Database,
			Architecture: spec.Architecture,
			DI:           spec.DI,
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

func selectionFromFlags() internal.TemplateSelection {
	sel := internal.TemplateSelection{
		Framework:    flagFramework,
		Database:     flagDatabase,
		Architecture: flagArchitecture,
		DI:           flagDI,
	}

	preset := internal.NormalizeArchitecture(flagPreset)
	switch preset {
	case internal.ArchitectureModular:
		if sel.Architecture == "" {
			sel.Architecture = internal.ArchitectureModular
		}
		if sel.DI == "" {
			sel.DI = internal.DINone
		}
	case internal.ArchitectureFullSetup:
		if sel.Architecture == "" {
			sel.Architecture = internal.ArchitectureFullSetup
		}
	case "":
	default:
		if flagPreset != "" {
			sel.Architecture = preset
		}
	}

	return sel
}

func defaultSelection(sel internal.TemplateSelection) internal.TemplateSelection {
	if sel.Framework == "" {
		sel.Framework = firstValue(internal.OptionsFor(sel, internal.TemplateFieldFramework))
	}
	sel.Framework = internal.NormalizeFramework(sel.Framework)

	if sel.Database == "" {
		sel.Database = firstValue(internal.OptionsFor(sel, internal.TemplateFieldDatabase))
	}
	sel.Database = internal.NormalizeDatabase(sel.Database)

	if sel.Architecture == "" {
		sel.Architecture = firstValue(internal.OptionsFor(sel, internal.TemplateFieldArchitecture))
	}
	sel.Architecture = internal.NormalizeArchitecture(sel.Architecture)

	if sel.Architecture == internal.ArchitectureModular {
		if sel.DI == "" {
			sel.DI = firstValue(internal.OptionsFor(sel, internal.TemplateFieldDI))
		}
		sel.DI = internal.NormalizeDI(sel.DI)
	} else if sel.DI != "" {
		sel.DI = internal.NormalizeDI(sel.DI)
	}

	return sel
}

func runNewWizard(ctx context.Context, w *internal.Wizard, projectName, modulePath string, sel internal.TemplateSelection) (internal.WizardConfig, bool, error) {
	var err error

	sel.Framework, err = selectTemplateValue(ctx, w, "Select framework", internal.TemplateFieldFramework, sel, sel.Framework)
	if err != nil {
		return internal.WizardConfig{}, false, err
	}

	sel.Database, err = selectTemplateValue(ctx, w, "Select database", internal.TemplateFieldDatabase, sel, sel.Database)
	if err != nil {
		return internal.WizardConfig{}, false, err
	}

	sel.Architecture, err = selectTemplateValue(ctx, w, "Select architecture", internal.TemplateFieldArchitecture, sel, sel.Architecture)
	if err != nil {
		return internal.WizardConfig{}, false, err
	}

	if len(internal.OptionsFor(sel, internal.TemplateFieldDI)) > 0 {
		sel.DI, err = selectTemplateValue(ctx, w, "Dependency Injection", internal.TemplateFieldDI, sel, sel.DI)
		if err != nil {
			return internal.WizardConfig{}, false, err
		}
	}

	spec, err := internal.ResolveTemplate(sel)
	if err != nil {
		return internal.WizardConfig{}, false, err
	}

	if modulePath == "" {
		modulePath, err = w.Input(ctx, "Module path", projectName)
		if err != nil {
			return internal.WizardConfig{}, false, err
		}
	}

	cfg := internal.WizardConfig{
		ProjectName: projectName,
		ModulePath:  modulePath,
		Template:    spec,
	}
	printConfiguration(cfg)

	ok, err := w.Confirm(ctx, "Continue?", true)
	if err != nil {
		return internal.WizardConfig{}, false, err
	}
	return cfg, ok, nil
}

func selectTemplateValue(ctx context.Context, w *internal.Wizard, label, field string, sel internal.TemplateSelection, current string) (string, error) {
	opts := internal.OptionsFor(sel, field)
	if current != "" {
		value := normalizeByField(field, current)
		if !hasOption(opts, value) {
			return "", fmt.Errorf("%s %q is not available for the current selection", field, current)
		}
		return value, nil
	}

	def := firstValue(opts)
	return w.SelectOption(ctx, label, opts, def)
}

func normalizeByField(field, value string) string {
	switch field {
	case internal.TemplateFieldFramework:
		return internal.NormalizeFramework(value)
	case internal.TemplateFieldDatabase:
		return internal.NormalizeDatabase(value)
	case internal.TemplateFieldArchitecture:
		return internal.NormalizeArchitecture(value)
	case internal.TemplateFieldDI:
		return internal.NormalizeDI(value)
	default:
		return value
	}
}

func firstValue(opts []internal.Option) string {
	if len(opts) == 0 {
		return ""
	}
	return opts[0].Value
}

func hasOption(opts []internal.Option, value string) bool {
	for _, opt := range opts {
		if opt.Value == value {
			return true
		}
	}
	return false
}

func templateSpecFromID(id string) (internal.TemplateSpec, error) {
	spec, err := internal.TemplateByID(id)
	if err == nil {
		return spec, nil
	}

	templates, listErr := internal.ListTemplates()
	if listErr != nil {
		return internal.TemplateSpec{}, listErr
	}
	for _, tpl := range templates {
		if tpl == id {
			return internal.TemplateSpec{ID: id}, nil
		}
	}
	return internal.TemplateSpec{}, err
}

func printConfiguration(cfg internal.WizardConfig) {
	fmt.Println()
	fmt.Println("Configuration")
	fmt.Printf("%-13s: %s\n", "Project Name", cfg.ProjectName)
	fmt.Printf("%-13s: %s\n", "Module Path", cfg.ModulePath)
	fmt.Printf("%-13s: %s\n", "Framework", internal.PrettyChoiceLabel(cfg.Template.Framework))
	fmt.Printf("%-13s: %s\n", "Database", internal.PrettyChoiceLabel(cfg.Template.Database))
	fmt.Printf("%-13s: %s\n", "Architecture", internal.PrettyChoiceLabel(cfg.Template.Architecture))
	if cfg.Template.DI != "" {
		fmt.Printf("%-13s: %s\n", "DI", internal.PrettyChoiceLabel(cfg.Template.DI))
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVar(&flagModule, "module", "", "Go module path (default: project-name)")
	newCmd.Flags().StringVar(&flagTemplate, "template", "", "Template ID to use (skips interactive template selection)")
	newCmd.Flags().StringVar(&flagFramework, "framework", "", "Framework to use (gin, fiber)")
	newCmd.Flags().StringVar(&flagDatabase, "db", "", "Database to use (postgres, mysql)")
	newCmd.Flags().StringVar(&flagArchitecture, "architecture", "", "Architecture to use (modular, layered, full-setup)")
	newCmd.Flags().StringVar(&flagDI, "di", "", "Dependency injection to use (none, uber-dig, uber-fx)")
	newCmd.Flags().StringVar(&flagPreset, "preset", "", "Legacy preset (base, full)")
	_ = newCmd.Flags().MarkHidden("preset")
}
