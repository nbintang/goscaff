package internal

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type Renderer interface {
	RenderDirTo(srcRoot, dstBase string, opts ScaffoldOptions) error
	RenderDir(srcRoot string, opts ScaffoldOptions) error
	RenderFileTo(srcFile, dstFile string, opts ScaffoldOptions) error
}

type rendererImpl struct {
	templateFS fs.FS
}

var scaffoldTemplateMarker = regexp.MustCompile(`\{\{\s*\.(PROJECT_NAME|MODULE_PATH|TEMPLATE|FRAMEWORK|DATABASE|ARCHITECTURE|DI)\b`)

func NewRenderer() Renderer {
	fsys, err := templateFS()
	if err != nil {
		return &rendererImpl{}
	}

	return &rendererImpl{templateFS: fsys}
}

func (r *rendererImpl) RenderDirTo(srcRoot, dstBase string, opts ScaffoldOptions) error {
	if r.templateFS == nil {
		return fmt.Errorf("template filesystem is not available")
	}

	entries, err := fs.ReadDir(r.templateFS, srcRoot)
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

		b, err := fs.ReadFile(r.templateFS, srcPath)
		if err != nil {
			return err
		}

		if shouldRenderTemplate(srcPath, b) {
			outPath = strings.TrimSuffix(outPath, ".tmpl")

			rendered, err := renderTemplate(srcPath, b, opts)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(outPath, rendered, 0o644); err != nil {
				return err
			}
			continue
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

	if r.templateFS == nil {
		return fmt.Errorf("template filesystem is not available")
	}

	b, err := fs.ReadFile(r.templateFS, srcFile)
	if err != nil {
		return err
	}

	outPath := filepath.Join(opts.OutDir, dstFile)

	if shouldRenderTemplate(srcFile, b) {
		rendered, err := renderTemplate(srcFile, b, opts)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		return os.WriteFile(outPath, rendered, 0o644)
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

func shouldRenderTemplate(path string, content []byte) bool {
	return strings.HasSuffix(path, ".tmpl") || scaffoldTemplateMarker.Match(content)
}

func renderTemplate(path string, content []byte, opts ScaffoldOptions) ([]byte, error) {
	t, err := template.New(filepath.Base(path)).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", path, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]any{
		"PROJECT_NAME": opts.ProjectName,
		"MODULE_PATH":  opts.ModulePath,
		"TEMPLATE":     opts.Template,
		"FRAMEWORK":    opts.Framework,
		"DATABASE":     opts.Database,
		"ARCHITECTURE": opts.Architecture,
		"DI":           opts.DI,
	}); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", path, err)
	}

	return buf.Bytes(), nil
}
