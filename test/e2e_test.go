package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type Case struct {
	Name string
	Args []string
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
		{"base-postgres", []string{"--preset", "base", "--db", "postgres"}},
		{"base-mysql", []string{"--preset", "base", "--db", "mysql"}},
		{"full-postgres", []string{"--preset", "full", "--db", "postgres"}},
		{"full-mysql", []string{"--preset", "full", "--db", "mysql"}},
		{"gin-postgres-uber-fx", []string{"--framework", "gin", "--db", "postgres", "--architecture", "modular", "--di", "uber-fx"}},
	}

	root := repoRoot(t)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			workspace, err := os.MkdirTemp(root, ".goscaff-e2e-"+tc.Name+"-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(workspace)

			// ✅ unik per subtest (biar gak tabrakan di repo root)
			projectName := "myapp-" + tc.Name
			projectDst := filepath.Join(workspace, projectName)

			// safety: kalau ada sisa folder dari crash sebelumnya
			_ = os.RemoveAll(projectDst)

			// 1) generate dari repo root ke workspace satu drive
			args := append([]string{"run", ".", "new", projectDst}, tc.Args...)
			runCmd(t, root, "go", args...)

			// pastikan kebentuk
			if _, err := os.Stat(projectDst); err != nil {
				t.Fatalf("generated project not found at %s: %v", projectDst, err)
			}

			// 2) build
			runCmd(t, projectDst, "go", "build", "./cmd/api")
		})
	}
}
