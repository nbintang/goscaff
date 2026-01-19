package scaffold

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Renderer interface {
	RenderDirTo(srcRoot, dstBase string, opts ScaffoldOptions) error
	RenderDir(srcRoot string, opts ScaffoldOptions) error
	RenderFileTo(srcFile, dstFile string, opts ScaffoldOptions) error
}

type rendererImpl struct {
}

func NewRenderer() Renderer {
	return &rendererImpl{}
}

func (r *rendererImpl) RenderDirTo(srcRoot, dstBase string, opts ScaffoldOptions) error {
	entries, err := templateFS.ReadDir(srcRoot)
	if err != nil {
		return err
	}

	for _, e := range entries {
		srcPath := filepath.ToSlash(filepath.Join(srcRoot, e.Name()))

		outPath := filepath.Join(opts.OutDir, dstBase, e.Name())

		if e.IsDir() {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return err
			}
			if err := r.RenderDirTo(srcPath, filepath.Join(dstBase, e.Name()), opts); err != nil {
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

func (r *rendererImpl) RenderFileTo(srcFile, dstFile string, opts ScaffoldOptions) error {
 

	srcFile = filepath.ToSlash(srcFile)
 
	b, err := templateFS.ReadFile(srcFile)
	if err != nil {
		return err
	}

	outPath := filepath.Join(opts.OutDir, dstFile)

	// kalau .tmpl → execute template
	if strings.HasSuffix(srcFile, ".tmpl") {
		t, err := template.New(filepath.Base(srcFile)).Parse(string(b))
		if err != nil {
			return fmt.Errorf("parse template %s: %w", srcFile, err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, map[string]any{
			"PROJECT_NAME": opts.ProjectName,
			"MODULE_PATH":  opts.ModulePath,
			"DB":           opts.DB,
			"PRESET":       opts.Preset,
		}); err != nil {
			return fmt.Errorf("execute template %s: %w", srcFile, err)
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		return os.WriteFile(outPath, buf.Bytes(), 0o644)
	}

	// non-template: copy raw
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, b, 0o644)
}

func (r *rendererImpl) RenderDir(srcRoot string, opts ScaffoldOptions) error {
	return r.RenderDirTo(srcRoot, "", opts)
}
