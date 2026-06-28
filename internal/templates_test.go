package internal

import "testing"

func TestResolveTemplate(t *testing.T) {
	tests := []struct {
		name string
		give TemplateSelection
		want string
	}{
		{
			name: "gin postgres modular uber fx",
			give: TemplateSelection{
				Framework:    FrameworkGin,
				Database:     DatabasePostgreSQL,
				Architecture: ArchitectureModular,
				DI:           DIUberFx,
			},
			want: "gin-postgres-uber-fx-[modular]",
		},
		{
			name: "gin mysql layered",
			give: TemplateSelection{
				Framework:    FrameworkGin,
				Database:     DatabaseMySQL,
				Architecture: ArchitectureLayered,
			},
			want: "gin-mysql-[layered]",
		},
		{
			name: "fiber postgres full setup",
			give: TemplateSelection{
				Framework:    FrameworkFiber,
				Database:     DatabasePostgreSQL,
				Architecture: ArchitectureFullSetup,
			},
			want: "fiber-full-postgres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveTemplate(tt.give)
			if err != nil {
				t.Fatalf("ResolveTemplate() error = %v", err)
			}
			if got.ID != tt.want {
				t.Fatalf("ResolveTemplate().ID = %q, want %q", got.ID, tt.want)
			}
		})
	}
}

func TestResolveTemplate_DIOnlyAvailableForModular(t *testing.T) {
	_, err := ResolveTemplate(TemplateSelection{
		Framework:    FrameworkGin,
		Database:     DatabasePostgreSQL,
		Architecture: ArchitectureLayered,
		DI:           DIUberFx,
	})
	if err == nil {
		t.Fatal("ResolveTemplate() error = nil, want error")
	}
}

func TestOptionsFor(t *testing.T) {
	t.Run("fiber only supports postgres full setup", func(t *testing.T) {
		dbs := OptionsFor(TemplateSelection{Framework: FrameworkFiber}, TemplateFieldDatabase)
		assertOptionValues(t, dbs, []string{DatabasePostgreSQL})

		architectures := OptionsFor(TemplateSelection{
			Framework: FrameworkFiber,
			Database:  DatabasePostgreSQL,
		}, TemplateFieldArchitecture)
		assertOptionValues(t, architectures, []string{ArchitectureFullSetup})
	})

	t.Run("gin modular supports all DI options", func(t *testing.T) {
		options := OptionsFor(TemplateSelection{
			Framework:    FrameworkGin,
			Database:     DatabasePostgreSQL,
			Architecture: ArchitectureModular,
		}, TemplateFieldDI)
		assertOptionValues(t, options, []string{DINone, DIUberDig, DIUberFx})
	})
}

func assertOptionValues(t *testing.T, got []Option, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(options) = %d, want %d; options = %#v", len(got), len(want), got)
	}

	for i := range want {
		if got[i].Value != want[i] {
			t.Fatalf("options[%d].Value = %q, want %q", i, got[i].Value, want[i])
		}
	}
}
