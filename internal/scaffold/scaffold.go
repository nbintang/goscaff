package scaffold

import (
	"embed"
	"fmt"
	"os"

	"github.com/nbintang/goscaff/pkg"
)


type Scaffold interface {
	Generate() error
}

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
 
func (s *scaffoldImpl)  envTemplatePath(db DBType, preset PresetType) string {
	base := "templates/utils/env"

	switch preset {
	case PresetFull:
		switch db {
		case DBTypeMySQL:
			return base + "/full/mysql.env.example.tmpl"
		default:
			return base + "/full/postgres.env.example.tmpl"
		}
	default:  
		switch db {
		case DBTypeMySQL:
			return base + "/base/mysql.env.example.tmpl"
		default:
			return base + "/base/postgres.env.example.tmpl"
		}
	}
}


func NewScaffold(opts ScaffoldOptions, renderer Renderer) Scaffold {
	return &scaffoldImpl{
		templateFS: templateFS,
		opts:       opts,
		renderer:   renderer,
	}
} 

func (s *scaffoldImpl) Generate() error { 
	presetRoot := "templates/base"
	if s.opts.Preset == PresetFull {
		presetRoot = "templates/full"
	}
  
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
 
	pkg.Action("Rendering preset (" + string(s.opts.Preset) + ")")
	if err := s.renderer.RenderDir(presetRoot, s.opts); err != nil {
		return err
	}
	pkg.Success("Preset rendered")
 
	pkg.Action("Generating environment template")
	if err := s.renderer.RenderFileTo(
		s.envTemplatePath(s.opts.DB, s.opts.Preset),
		".env.example",
		s.opts,
	); err != nil {
		return err
	}
	pkg.Success(".env.example generated")
 
	if s.opts.Preset == PresetFull {
		pkg.Action("Applying database overlays (" + string(s.opts.DB) + ")")
		overlays := buildDbOverlayFiles(dbRoot, s.opts.DB)
		if err := applyOverlayFiles(overlays, s.renderer, s.opts); err != nil {
			return err
		}
		pkg.Success("Database applied")
	} else {
		pkg.Action("Skipping database overlays (base preset)")
		pkg.Success("Skipped")
	}

	pkg.Action("Running: go mod tidy")
	if err := pkg.RunVerbose(s.opts.OutDir, "go", "mod", "tidy"); err != nil {
		return err
	}
	pkg.Success("Dependencies installed")

	pkg.Action("Initializing git repository")
	_ = pkg.RunQuiet(s.opts.OutDir, "git", "init")
	pkg.Success("Git initialized")

	return nil
}