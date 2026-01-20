package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type Case struct {
	Name   string
	Preset string
	DB     string
}

func repoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found (repo root not detected)")
		}
		dir = parent
	}
}

func runCmd(t *testing.T, dir, name string, args ...string) string {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
		"GOSCAFF_NON_INTERACTIVE=1", // kalau nanti kamu implement
	)

	if err := cmd.Run(); err != nil {
		t.Fatalf(
			"\nCOMMAND FAILED\nDir: %s\nCmd: %s %v\nOutput:\n%s\nErr: %v\n",
			dir, name, args, out.String(), err,
		)
	}

	return out.String()
}

func Test_Scaffold_E2E(t *testing.T) {
	cases := []Case{
		{"base-postgres", "base", "postgres"},
		{"base-mysql", "base", "mysql"},
		{"full-postgres", "full", "postgres"},
		{"full-mysql", "full", "mysql"},
	}

	root := repoRoot(t)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			workspace := t.TempDir()

			// ✅ unik per subtest (biar gak tabrakan di repo root)
			projectName := "myapp-" + tc.Name
			projectSrc := filepath.Join(root, projectName)
			projectDst := filepath.Join(workspace, projectName)

			// safety: kalau ada sisa folder dari crash sebelumnya
			_ = os.RemoveAll(projectSrc)

			// 1) generate di repo root pakai RELATIVE name
			runCmd(t, root, "go",
				"run", ".", "new", projectName,
				"--preset", tc.Preset,
				"--db", tc.DB,
			)

			// pastikan kebentuk
			if _, err := os.Stat(projectSrc); err != nil {
				t.Fatalf("generated project not found at %s: %v", projectSrc, err)
			}

			// 2) pindahkan ke workspace (unik juga)
			if err := os.Rename(projectSrc, projectDst); err != nil {
				t.Fatalf("failed to move project to temp dir: %v", err)
			}

			// 3) build
			runCmd(t, projectDst, "go", "build", "./cmd/api")
		})
	}
}
