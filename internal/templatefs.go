package internal

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:templates
var embeddedTemplates embed.FS

func templateFS() (fs.FS, error) {
	// Use embedded templates (compiled into the binary)
	sub, err := fs.Sub(embeddedTemplates, "templates")
	if err == nil {
		return sub, nil
	}

	// Fallback: try filesystem paths (useful during development)
	for _, root := range templateRoots() {
		info, err := os.Stat(root)
		if err != nil {
			continue
		}
		if info.IsDir() {
			return os.DirFS(root), nil
		}
	}

	return nil, os.ErrNotExist
}

func templateRoots() []string {
	var roots []string

	if wd, err := os.Getwd(); err == nil {
		roots = append(roots, filepath.Join(wd, "internal", "templates"))
	}

	if exe, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Join(filepath.Dir(exe), "internal", "templates"))
	}

	return roots
}
