package test

import (
	"bytes" 
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found from %s upward", wd)
		}
		dir = parent
	}
}


func run(t *testing.T, dir string, name string, args ...string) string {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	// isolasi cache biar gak “ketularan” cache mesin
	cacheRoot := t.TempDir()
	cmd.Env = append(os.Environ(),
		"GOCACHE="+filepath.Join(cacheRoot, "gocache"),
		"GOMODCACHE="+filepath.Join(cacheRoot, "gomodcache"),
		"GIT_TERMINAL_PROMPT=0",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %s %v\nDir: %s\nOutput:\n%s\nErr: %v",
			name, args, dir, out.String(), err,
		)
	}

	return out.String()
}

type Case struct {
	Name   string
	Preset string
	DB     string
}

func GenerateAndBuild(t *testing.T, tc Case) {
	t.Helper()

	root := repoRoot(t)

	workspace := t.TempDir()
	project := filepath.Join(workspace, "myapp")

	// 1) go run . new myapp ...
	args := []string{"run", ".", "new", project}
	if tc.Preset != "" {
		args = append(args, "--preset", tc.Preset)
	}
	if tc.DB != "" {
		args = append(args, "--db", tc.DB)
	}

	run(t, root, "go", args...)

	// 2) cd myapp (project dir)
	// 3) go build ./cmd/api
	run(t, project, "go", "build", "./cmd/api")
}
