package scaffold

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func renderDirTo(srcRoot, outRoot, dstBase string, opts Options) error {
	entries, err := templateFS.ReadDir(srcRoot)
	if err != nil {
		return err
	}

	for _, e := range entries {
		srcPath := filepath.ToSlash(filepath.Join(srcRoot, e.Name()))
 
		outPath := filepath.Join(outRoot, dstBase, e.Name())

		if e.IsDir() {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return err
			} 
			if err := renderDirTo(srcPath, outRoot, filepath.Join(dstBase, e.Name()), opts); err != nil {
				return err
			}
			continue
		}

		if strings.HasSuffix(e.Name(), ".tmpl") {
			outPath = strings.TrimSuffix(outPath, ".tmpl")

			b, err := templateFS.ReadFile(srcPath)
			if err != nil {
				return err
			}

			t, err := template.New(e.Name()).Parse(string(b))
			if err != nil {
				return fmt.Errorf("parse template %s: %w", srcPath, err)
			}

			var buf bytes.Buffer
			if err := t.Execute(&buf, map[string]any{
				"PROJECT_NAME": opts.ProjectName,
				"MODULE_PATH":  opts.ModulePath,
				"DB":           opts.DB,
			}); err != nil {
				return fmt.Errorf("execute template %s: %w", srcPath, err)
			}

			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
				return err
			}
			continue
		}

		b, err := templateFS.ReadFile(srcPath)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, b, 0o644); err != nil {
			return err
		}
	}

	return nil
}

func renderDir(srcRoot, outRoot string, opts Options) error {
	return renderDirTo(srcRoot, outRoot, "", opts)
}