package scaffold

import (
	"fmt"
	"path/filepath"

	"github.com/nbintang/goscaff/pkg"
)

type Printer interface {
	PrintNextSteps()
}

type printerImpl struct {
	opts ScaffoldOptions
}

func NewPrinter(opts ScaffoldOptions) Printer {
	return &printerImpl{opts}
}

func (p *printerImpl) printStep(title string, commands ...string) {
	fmt.Println()
	pkg.ColorStepTitle.Printf("  %s\n", title)
	for _, c := range commands {
		pkg.ColorBullet.Print("    $ ")
		pkg.ColorCmd.Println(c)
	}
}

func (p *printerImpl) PrintNextSteps() {
	fmt.Println()

	pkg.ColorOk.Printf("✓ ")
	pkg.ColorOk.Println("Project generated successfully")

	pkg.ColorTip.Print("⚑ ")
	pkg.ColorTip.Println("Review your environment variables before running")

	fmt.Println()
	pkg.ColorHeader.Println("Next steps")
	pkg.ColorBullet.Println("────────────────────────────────────────")

	projectDir := filepath.Base(filepath.Clean(p.opts.OutDir))
	p.printStep("Go to project directory", fmt.Sprintf("cd %s", projectDir))
	p.printStep("Setup environment", "cp .env.example .env.local")
	pkg.ColorNote.Println("    Configure database and app settings inside .env.local before running the project.")

	if p.opts.Preset == "full" {
		p.printFullPresetSteps()
	} else {
		p.printBasePresetSteps()
	}

	fmt.Println()
	pkg.ColorNote.Println("  • Server: http://localhost:8080")
	pkg.ColorNote.Println("  • Edit .env.local if config changes")
	fmt.Println()
}

func (p *printerImpl) printFullPresetSteps() {
	fmt.Println()
	pkg.ColorStepTitle.Println("  FULL preset detected")

	pkg.ColorNote.Println("    This preset uses Makefile and Air for development.")
	pkg.ColorNote.Println("    Make sure you have the following installed:")

	pkg.ColorBullet.Println("      - make")
	pkg.ColorBullet.Println("      - air (live reload)")
	pkg.ColorBullet.Println()

	p.printStep("Install Air (if not installed)", "go install github.com/air-verse/air@latest")
	p.printStep("Start dependencies", "make docker")
	p.printStep("Run migration", "make migrate")
	p.printStep("Run seed", "make seed")
	p.printStep("Run development server", "make dev")
}

func (p *printerImpl) printBasePresetSteps() {
	fmt.Println()
	pkg.ColorStepTitle.Println("  BASE preset detected")

	p.printStep("Run migration", "go run ./cmd/migrate")
	p.printStep("Run seed", "go run ./cmd/seed")
	p.printStep("Run the app", "go run ./cmd/api")
}
