package internal

import (
	"sort"
)

func ListTemplates() ([]string, error) {
	entries, err := templateFS.ReadDir("templates")
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
		return "Gin Full Postgres"
	case "gin-full-mysql":
		return "Gin Full Mysql"
	case "gin-mysql-modular":
		return "Gin Mysql (Modular)"
	case "gin-mysql-uber-dig-modular":
		return "Gin Mysql Uber Dig (Modular)"
	case "gin-mysql-uber-fx-modular":
		return "Gin Mysql Uber Fx (Modular)"
	case "gin-postgres-modular":
		return "Gin Postgres (Modular)"
	case "gin-postgres-uber-dig-modular":
		return "Gin Postgres Uber Dig (Modular)"
	case "gin-postgres-uber-fx-modular":
		return "Gin Postgres Uber Fx (Modular)"
	default:
		// fallback sederhana
		return id
	}
}
