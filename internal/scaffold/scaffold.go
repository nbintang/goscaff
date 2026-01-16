package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed all:templates/**
var templateFS embed.FS

type Options struct {
	ProjectName string
	ModulePath  string
	DB          string
	Preset      string // "base" | "full"
}

func info(format string, args ...any) {
	fmt.Printf("• "+format+"\n", args...)
}

func Generate(outDir string, opts Options) error {
	presetRoot := "templates/base"
	if opts.Preset == "full" {
		presetRoot = "templates/full"
	}

	dbRoot := "templates/db/postgres"
	if opts.DB == "mysql" {
		dbRoot = "templates/db/mysql"
	}

	dstBase := filepath.Join("internal", "infra", "database")

	fmt.Println()
	header("Goscaff • Project Generator")
	infoLine("Folder : " + outDir)
	infoLine("Preset : " + opts.Preset)
	infoLine("DB     : " + opts.DB)
	fmt.Println()

	action("Creating project directory")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	ok("Directory created")

	action("Rendering preset (" + opts.Preset + ")")
	if err := renderDir(presetRoot, outDir, opts); err != nil {
		return err
	}
	ok("Preset rendered")

	// kalau kamu masih mau overlay selalu jalan, ya biarin.
	// Tapi kalau mau base bersih, taruh if opts.Preset == "full"
	action("Rendering database driver (" + opts.DB + ")")
	if err := renderDirTo(dbRoot, outDir, dstBase, opts); err != nil {
		return err
	}
	ok("Database applied")

	action("Running: go mod tidy")
	if err := runVerbose(outDir, "go", "mod", "tidy"); err != nil {
		return err
	}
	ok("Dependencies installed")

	action("Initializing git repository")
	_ = runQuiet(outDir, "git", "init")
	ok("Git initialized")

	printNextSteps(outDir, opts)
	return nil
}
