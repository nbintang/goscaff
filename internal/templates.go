package internal

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

const (
	TemplateFieldFramework    = "framework"
	TemplateFieldDatabase     = "database"
	TemplateFieldArchitecture = "architecture"
	TemplateFieldDI           = "di"

	FrameworkGin   = "gin"
	FrameworkFiber = "fiber"

	DatabasePostgreSQL = "postgresql"
	DatabaseMySQL      = "mysql"

	ArchitectureFullSetup = "full-setup"
	ArchitectureModular   = "modular"
	ArchitectureLayered   = "layered"

	DINone    = "none"
	DIUberDig = "uber-dig"
	DIUberFx  = "uber-fx"
)

type TemplateSpec struct {
	ID           string
	Framework    string
	Database     string
	Architecture string
	DI           string
}

type TemplateSelection struct {
	Framework    string
	Database     string
	Architecture string
	DI           string
}

type WizardConfig struct {
	ProjectName string
	ModulePath  string
	Template    TemplateSpec
}

var templateCatalog = []TemplateSpec{
	{
		ID:           "fiber-full-postgres",
		Framework:    FrameworkFiber,
		Database:     DatabasePostgreSQL,
		Architecture: ArchitectureFullSetup,
	},
	{
		ID:           "gin-full-postgres",
		Framework:    FrameworkGin,
		Database:     DatabasePostgreSQL,
		Architecture: ArchitectureFullSetup,
	},
	{
		ID:           "gin-full-mysql",
		Framework:    FrameworkGin,
		Database:     DatabaseMySQL,
		Architecture: ArchitectureFullSetup,
	},
	{
		ID:           "gin-postgres-layered",
		Framework:    FrameworkGin,
		Database:     DatabasePostgreSQL,
		Architecture: ArchitectureLayered,
	},
	{
		ID:           "gin-mysql-layered",
		Framework:    FrameworkGin,
		Database:     DatabaseMySQL,
		Architecture: ArchitectureLayered,
	},
	{
		ID:           "gin-postgres-modular",
		Framework:    FrameworkGin,
		Database:     DatabasePostgreSQL,
		Architecture: ArchitectureModular,
		DI:           DINone,
	},
	{
		ID:           "gin-mysql-modular",
		Framework:    FrameworkGin,
		Database:     DatabaseMySQL,
		Architecture: ArchitectureModular,
		DI:           DINone,
	},
	{
		ID:           "gin-postgres-uber-dig-modular",
		Framework:    FrameworkGin,
		Database:     DatabasePostgreSQL,
		Architecture: ArchitectureModular,
		DI:           DIUberDig,
	},
	{
		ID:           "gin-mysql-uber-dig-modular",
		Framework:    FrameworkGin,
		Database:     DatabaseMySQL,
		Architecture: ArchitectureModular,
		DI:           DIUberDig,
	},
	{
		ID:           "gin-postgres-uber-fx-modular",
		Framework:    FrameworkGin,
		Database:     DatabasePostgreSQL,
		Architecture: ArchitectureModular,
		DI:           DIUberFx,
	},
	{
		ID:           "gin-mysql-uber-fx-modular",
		Framework:    FrameworkGin,
		Database:     DatabaseMySQL,
		Architecture: ArchitectureModular,
		DI:           DIUberFx,
	},
}

func ListTemplates() ([]string, error) {
	fsys, err := templateFS()
	if err != nil {
		return nil, fmt.Errorf("open templates directory: %w", err)
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		// skip folder util kalau ada
		if name == "utils" || name == "_shared" {
			continue
		}
		out = append(out, name)
	}

	sort.Strings(out)
	return out, nil
}

func TemplateCatalog() []TemplateSpec {
	out := make([]TemplateSpec, len(templateCatalog))
	copy(out, templateCatalog)
	return out
}

func ResolveTemplate(selection TemplateSelection) (TemplateSpec, error) {
	sel := normalizeTemplateSelection(selection)
	if sel.Architecture != ArchitectureModular && sel.DI != "" && sel.DI != DINone {
		return TemplateSpec{}, fmt.Errorf("dependency injection %q is only available for modular architecture", selection.DI)
	}

	for _, spec := range templateCatalog {
		if spec.Framework != sel.Framework || spec.Database != sel.Database || spec.Architecture != sel.Architecture {
			continue
		}
		if spec.Architecture != ArchitectureModular {
			return spec, nil
		}
		if spec.DI == sel.DI {
			return spec, nil
		}
	}

	return TemplateSpec{}, fmt.Errorf("no template found for framework=%q database=%q architecture=%q di=%q", selection.Framework, selection.Database, selection.Architecture, selection.DI)
}

func TemplateByID(id string) (TemplateSpec, error) {
	for _, spec := range templateCatalog {
		if spec.ID == id {
			return spec, nil
		}
	}
	return TemplateSpec{}, fmt.Errorf("unknown template %q", id)
}

func OptionsFor(selection TemplateSelection, field string) []Option {
	sel := normalizeTemplateSelection(selection)
	field = normalizeKey(field)

	seen := make(map[string]bool)
	var opts []Option
	for _, spec := range templateCatalog {
		if sel.Framework != "" && spec.Framework != sel.Framework {
			continue
		}
		if sel.Database != "" && spec.Database != sel.Database {
			continue
		}
		if sel.Architecture != "" && spec.Architecture != sel.Architecture {
			continue
		}

		var value string
		switch field {
		case TemplateFieldFramework:
			value = spec.Framework
		case TemplateFieldDatabase:
			value = spec.Database
		case TemplateFieldArchitecture:
			value = spec.Architecture
		case TemplateFieldDI:
			value = spec.DI
		default:
			return nil
		}
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		opts = append(opts, Option{
			Label: PrettyChoiceLabel(value),
			Value: value,
		})
	}

	sort.SliceStable(opts, func(i, j int) bool {
		return optionRank(field, opts[i].Value) < optionRank(field, opts[j].Value)
	})
	return opts
}

func PrettyTemplateLabel(id string) string {
	spec, err := TemplateByID(id)
	if err != nil {
		return id
	}

	label := PrettyChoiceLabel(spec.Framework) + " + " + PrettyChoiceLabel(spec.Database) + " + " + PrettyChoiceLabel(spec.Architecture)
	if spec.DI != "" && spec.DI != DINone {
		label += " + " + PrettyChoiceLabel(spec.DI)
	}
	return label
}

func PrettyChoiceLabel(value string) string {
	switch normalizeValue(value) {
	case FrameworkGin:
		return "Gin"
	case FrameworkFiber:
		return "Fiber"
	case DatabasePostgreSQL:
		return "PostgreSQL"
	case DatabaseMySQL:
		return "MySQL"
	case ArchitectureFullSetup:
		return "Full Setup"
	case ArchitectureModular:
		return "Modular"
	case ArchitectureLayered:
		return "Layered"
	case DINone:
		return "None"
	case DIUberDig:
		return "Uber Dig"
	case DIUberFx:
		return "Uber Fx"
	default:
		return value
	}
}

func NormalizeFramework(value string) string {
	switch normalizeValue(value) {
	case "", FrameworkGin, FrameworkFiber:
		return normalizeValue(value)
	default:
		return normalizeValue(value)
	}
}

func NormalizeDatabase(value string) string {
	switch normalizeValue(value) {
	case "postgres", "postgresql", "pg":
		return DatabasePostgreSQL
	case "mysql", "my-sql":
		return DatabaseMySQL
	default:
		return normalizeValue(value)
	}
}

func NormalizeArchitecture(value string) string {
	switch normalizeValue(value) {
	case "full", "fullsetup", "full-setup":
		return ArchitectureFullSetup
	case "base", "modular":
		return ArchitectureModular
	case "layer", "layered":
		return ArchitectureLayered
	default:
		return normalizeValue(value)
	}
}

func NormalizeDI(value string) string {
	switch normalizeValue(value) {
	case "", "none", "no", "no-di":
		return DINone
	case "dig", "uber-dig", "uberdig":
		return DIUberDig
	case "fx", "uber-fx", "uberfx":
		return DIUberFx
	default:
		return normalizeValue(value)
	}
}

func normalizeTemplateSelection(selection TemplateSelection) TemplateSelection {
	sel := TemplateSelection{
		Framework:    NormalizeFramework(selection.Framework),
		Database:     NormalizeDatabase(selection.Database),
		Architecture: NormalizeArchitecture(selection.Architecture),
		DI:           NormalizeDI(selection.DI),
	}
	if sel.Architecture != ArchitectureModular && selection.DI == "" {
		sel.DI = ""
	}
	return sel
}

func normalizeKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeValue(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.Join(strings.Fields(value), "-")
	return value
}

func optionRank(field, value string) int {
	order := map[string][]string{
		TemplateFieldFramework:    {FrameworkGin, FrameworkFiber},
		TemplateFieldDatabase:     {DatabasePostgreSQL, DatabaseMySQL},
		TemplateFieldArchitecture: {ArchitectureModular, ArchitectureLayered, ArchitectureFullSetup},
		TemplateFieldDI:           {DINone, DIUberDig, DIUberFx},
	}
	for i, item := range order[field] {
		if item == value {
			return i
		}
	}
	return len(order[field])
}
