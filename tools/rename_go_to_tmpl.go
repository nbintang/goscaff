package tools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// RenameAllGoToTmpl walks through templateDir and renames:
//   something.go      -> something.go.tmpl
// It skips:
//   - already *.go.tmpl
//   - non-.go files
//
// go run . ./internal/scaffold/templates
//
//
//

func RenameAllGoToTmpl(templateDir string, dryRun bool) error {
	templateDir = filepath.Clean(templateDir)

	return filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		name := d.Name()

		// skip already tmpl
		if strings.HasSuffix(name, ".go.tmpl") {
			return nil
		}

		// only .go
		if !strings.HasSuffix(name, ".go") {
			return nil
		}

		oldPath := path
		newPath := path + ".tmpl" // => file.go.tmpl

		// prevent overwrite if target exists
		if _, statErr := os.Stat(newPath); statErr == nil {
			return fmt.Errorf("target already exists, refusing to overwrite: %s", newPath)
		}

		if dryRun {
			fmt.Printf("[dry-run] %s -> %s\n", oldPath, newPath)
			return nil
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("rename %s -> %s: %w", oldPath, newPath, err)
		}

		fmt.Printf("[renamed] %s -> %s\n", oldPath, newPath)
		return nil
	})
}
