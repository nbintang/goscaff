package scaffold

import (
	"embed"
	"fmt"
	"os"

	"github.com/nbintang/goscaff/pkg"
)

//go:embed all:templates/**
var templateFS embed.FS

type Scaffold interface {
	Generate() error
}


type DBType string

const (
	DBTypePostgres DBType = "postgres"
	DBTypeMySQL    DBType = "mysql"
)

type PresetType string

const (
	PresetBase   PresetType = "base"
	PresetFull   PresetType = "full"
)

type ScaffoldOptions struct {
	ProjectName string
	ModulePath  string
	DB          DBType
	Preset      PresetType
	OutDir      string
}

type scaffoldImpl struct {
	templateFS embed.FS
	opts       ScaffoldOptions
	renderer   Renderer
}
type Overlay struct {
	Src string // file template path
	Dst string // output file path
}

func buildDbOverlayFiles(dbRoot string, db DBType) []Overlay {
	files := []Overlay{
		{Src: dbRoot + "/entity.go.tmpl", Dst: "internal/user/entity.go"},
		{Src: dbRoot + "/repository.go.tmpl", Dst: "internal/user/repository.go"},
		{Src: dbRoot + "/standalone.go.tmpl", Dst: "internal/infra/database/standalone.go"},

		{Src: dbRoot + "/migrate.go.tmpl", Dst: "cmd/migrate/init.go"},
		{Src: dbRoot + "/seed.go.tmpl", Dst: "cmd/seed/init.go"},
	}

	// Postgres-only helper (MySQL doesn't have CREATE TYPE enums).
	if db == DBTypePostgres {
		files = append(files, Overlay{
			Src: dbRoot + "/create_enums.go.tmpl",
			Dst: "cmd/migrate/create_enums.go",
		})
	}

	return files
}



// Keep env templates separate from db overlays to avoid confusion.
func envTemplatePath(db DBType, preset PresetType) string {
	base := "templates/utils/env"

	switch preset {
	case PresetFull:
		switch db {
		case DBTypeMySQL:
			return base + "/full/mysql.env.example.tmpl"
		default:
			return base + "/full/postgres.env.example.tmpl"
		}
	default: // PresetBase
		switch db {
		case DBTypeMySQL:
			return base + "/base/mysql.env.example.tmpl"
		default:
			return base + "/base/postgres.env.example.tmpl"
		}
	}
}


func (s *scaffoldImpl) applyOverlayFiles(files []Overlay) error {
	for _, f := range files {
		if err := s.renderer.RenderFileTo(f.Src, f.Dst, s.opts); err != nil {
			return err
		}
	}
	return nil
}
func NewScaffold(opts ScaffoldOptions, renderer Renderer) Scaffold {
	return &scaffoldImpl{
		templateFS: templateFS,
		opts:       opts,
		renderer:   renderer,
	}
} 

func (s *scaffoldImpl) Generate() error {
	// Preset source
	presetRoot := "templates/base"
	if s.opts.Preset == PresetFull {
		presetRoot = "templates/full"
	}

	// DB source
	dbRoot := "templates/utils/db/postgres"
	if s.opts.DB == DBTypeMySQL {
		dbRoot = "templates/utils/db/mysql"
	}

	fmt.Println()
	pkg.Header("Goscaff • Project Generator")
	pkg.Info("Folder : " + s.opts.OutDir)
	pkg.Info("Preset : " + string(s.opts.Preset))
	pkg.Info("DB     : " + string(s.opts.DB))
	fmt.Println()

	pkg.Action("Creating project directory")
	if err := os.MkdirAll(s.opts.OutDir, 0o755); err != nil {
		return err
	}
	pkg.Success("Directory created")

	// 1) Render preset to project root
	pkg.Action("Rendering preset (" + string(s.opts.Preset) + ")")
	if err := s.renderer.RenderDir(presetRoot, s.opts); err != nil {
		return err
	}
	pkg.Success("Preset rendered")

	// 2) Always generate .env.example (base/full)
	pkg.Action("Generating environment template")
	if err := s.renderer.RenderFileTo(
		envTemplatePath(s.opts.DB, s.opts.Preset),
		".env.example",
		s.opts,
	); err != nil {
		return err
	}
	pkg.Success(".env.example generated")

	// 3) Apply DB overlays ONLY for full preset
	if s.opts.Preset == PresetFull {
		pkg.Action("Applying database overlays (" + string(s.opts.DB) + ")")
		if err := s.applyOverlayFiles(buildDbOverlayFiles(dbRoot, s.opts.DB)); err != nil {
			return err
		}
		pkg.Success("Database applied")
	} else {
		pkg.Action("Skipping database overlays (base preset)")
		pkg.Success("Skipped")
	}

	pkg.Action("Running: go mod tidy")
	if err := runVerbose(s.opts.OutDir, "go", "mod", "tidy"); err != nil {
		return err
	}
	pkg.Success("Dependencies installed")

	pkg.Action("Initializing git repository")
	_ = runQuiet(s.opts.OutDir, "git", "init")
	pkg.Success("Git initialized")

	return nil
}