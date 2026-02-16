package internal

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
	Template    string // <-- NEW
	OutDir      string
}

type scaffoldImpl struct {
	templateFS embed.FS
	opts       ScaffoldOptions
	renderer   Renderer
}

func NewScaffold(opts ScaffoldOptions, renderer Renderer) Scaffold {
	return &scaffoldImpl{
		templateFS: templateFS,
		opts:       opts,
		renderer:   renderer,
	}
}

func (s *scaffoldImpl) Generate() error {
	templateRoot := "templates/" + s.opts.Template

	fmt.Println()
	pkg.Header("Goscaff • Project Generator")
	pkg.Info("Folder   : " + s.opts.OutDir)
	pkg.Info("Template : " + s.opts.Template)
	fmt.Println()

	pkg.Action("Creating project directory")
	if err := os.MkdirAll(s.opts.OutDir, 0o755); err != nil {
		return err
	}
	pkg.Success("Directory created")

	pkg.Action("Rendering template (" + s.opts.Template + ")")
	if err := s.renderer.RenderDir(templateRoot, s.opts); err != nil {
		return err
	}
	pkg.Success("Template rendered")

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

