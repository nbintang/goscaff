package internal

import (
	"fmt"
	"os"

	"github.com/nbintang/goscaff/pkg"
)

type Scaffold interface {
	Generate() error
}
type ScaffoldOptions struct {
	ProjectName  string
	ModulePath   string
	Template     string
	OutDir       string
	Framework    string
	Database     string
	Architecture string
	DI           string
}

type scaffoldImpl struct {
	opts     ScaffoldOptions
	renderer Renderer
}

func NewScaffold(opts ScaffoldOptions, renderer Renderer) Scaffold {
	return &scaffoldImpl{
		opts:     opts,
		renderer: renderer,
	}
}

func (s *scaffoldImpl) Generate() error {
	templateRoot := s.opts.Template

	fmt.Println()
	pkg.Header("Goscaff • Project Generator")
	pkg.Info("Folder       : " + s.opts.OutDir)
	pkg.Info("Project Name : " + s.opts.ProjectName)
	pkg.Info("Module Path  : " + s.opts.ModulePath)
	if s.opts.Framework != "" {
		pkg.Info("Framework    : " + PrettyChoiceLabel(s.opts.Framework))
	}
	if s.opts.Database != "" {
		pkg.Info("Database     : " + PrettyChoiceLabel(s.opts.Database))
	}
	if s.opts.Architecture != "" {
		pkg.Info("Architecture : " + PrettyChoiceLabel(s.opts.Architecture))
	}
	if s.opts.DI != "" {
		pkg.Info("DI           : " + PrettyChoiceLabel(s.opts.DI))
	}
	fmt.Println()

	pkg.Action("Creating project directory")
	if err := os.MkdirAll(s.opts.OutDir, 0o755); err != nil {
		return err
	}
	pkg.Success("Directory created")

	pkg.Action("Copying project template")
	if err := s.renderer.RenderDir(templateRoot, s.opts); err != nil {
		return err
	}
	pkg.Success("Project template copied")

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
