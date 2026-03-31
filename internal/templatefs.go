package internal

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

func templateFS() (fs.FS, error) {
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

	if _, file, _, ok := runtime.Caller(0); ok {
		roots = append(roots, filepath.Join(filepath.Dir(file), "templates"))
	}

	if wd, err := os.Getwd(); err == nil {
		roots = append(roots, filepath.Join(wd, "internal", "templates"))
	}

	if exe, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Join(filepath.Dir(exe), "internal", "templates"))
	}

	return roots
}
