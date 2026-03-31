package internal

import (
	"fmt"
	"io/fs"
	"sort"
)

func ListTemplates() ([]string, error) {
	fsys, err := templateFS()
	if err != nil {
		return nil, fmt.Errorf("open templates directory: %w", err)
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		// skip folder util kalau ada
		if name == "utils" || name == "_shared" {
			continue
		}
		out = append(out, name)
	}

	sort.Strings(out)
	return out, nil
}

func PrettyTemplateLabel(id string) string {
	switch id {
	case "gin-full-postgres":
		return "Gin & Postgres Full Setup (Modular)"
	case "gin-full-mysql":
		return "Gin & Mysql Full Setup (Modular)"
	case "gin-mysql-modular":
		return "Gin & Mysql (Modular)"
	case "gin-mysql-uber-dig-modular":
		return "Gin & Mysql & Uber Dig (Modular)"
	case "gin-mysql-uber-fx-modular":
		return "Gin & Mysql & Uber Fx (Modular)"
	case "gin-postgres-modular":
		return "Gin & Postgres (Modular)"
	case "gin-postgres-uber-dig-modular":
		return "Gin & Postgres & Uber Dig (Modular)"
	case "gin-postgres-uber-fx-modular":
		return "Gin & Postgres & Uber Fx (Modular)"
	default:
		// fallback sederhana
		return id
	}
}
