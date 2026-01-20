package scaffold

func buildDbOverlayFiles(dbRoot string, db DBType) (core []Overlay, full []Overlay) {
	core = []Overlay{
		{Src: dbRoot + "/entity.go.tmpl", Dst: "internal/user/entity.go"},
		{Src: dbRoot + "/repository.go.tmpl", Dst: "internal/user/repository.go"},
		{Src: dbRoot + "/standalone.go.tmpl", Dst: "internal/infra/database/standalone.go"},
		{Src: dbRoot + "/migrate.go.tmpl", Dst: "cmd/migrate/init.go"},
		{Src: dbRoot + "/seed.go.tmpl", Dst: "cmd/seed/init.go"},
		{Src: dbRoot + "/seed.go.tmpl", Dst: "cmd/seed/init.go"},

	}

	full = []Overlay{}

	if db == DBTypePostgres {
		core = append(core, Overlay{
			Src: dbRoot + "/create_enums.go.tmpl",
			Dst: "cmd/migrate/create_enums.go",
		})
	}

	return core, full
}

func applyOverlayFiles(files []Overlay, renderer Renderer, opts ScaffoldOptions) error {
	for _, f := range files {
		if err := renderer.RenderFileTo(f.Src, f.Dst, opts); err != nil {
			return err
		}
	}
	return nil
}
