package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmd := exec.Command("git", "-C", dir, "init")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	return dir
}

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

func TestIsRepo(t *testing.T) {
	dir := initRepo(t)
	if !IsRepo(dir) {
		t.Error("expected IsRepo=true for a git repo")
	}
	if IsRepo(t.TempDir()) {
		t.Error("expected IsRepo=false for a non-repo dir")
	}
}

func TestChangedFiles_WorkingTree_TrackedAndUntracked(t *testing.T) {
	dir := initRepo(t)
	write(t, dir, "main.go", "package main\n")
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "-m", "init")

	// modify tracked file
	write(t, dir, "main.go", "package main\n\nfunc main() {}\n")
	// add a brand-new untracked file
	write(t, dir, "internal/new.go", "package internal\n")

	changes, err := ChangedFiles(dir, RangeSpec{})
	if err != nil {
		t.Fatalf("ChangedFiles: %v", err)
	}

	paths := map[string]FileChange{}
	for _, c := range changes {
		paths[c.Path] = c
	}
	if _, ok := paths["main.go"]; !ok {
		t.Errorf("expected main.go in changes, got %+v", changes)
	}
	if nf, ok := paths["internal/new.go"]; !ok {
		t.Errorf("expected untracked internal/new.go in changes, got %+v", changes)
	} else if nf.Status != "A" {
		t.Errorf("expected untracked file status A, got %q", nf.Status)
	}
}

func TestChangedFiles_Staged(t *testing.T) {
	dir := initRepo(t)
	write(t, dir, "a.go", "package a\n")
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "-m", "init")

	write(t, dir, "b.go", "package b\n")
	gitRun(t, dir, "add", "b.go")

	changes, err := ChangedFiles(dir, RangeSpec{Staged: true})
	if err != nil {
		t.Fatalf("ChangedFiles staged: %v", err)
	}
	found := false
	for _, c := range changes {
		if c.Path == "b.go" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected b.go in staged changes, got %+v", changes)
	}
}
