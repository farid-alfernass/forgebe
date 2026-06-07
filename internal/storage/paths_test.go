package storage

import (
	"path/filepath"
	"testing"
)

func TestPathsBuildExpectedDirectories(t *testing.T) {
	root := t.TempDir()
	p := NewPathsWithRoot(root)
	if p.ProjectsDir() != filepath.Join(root, ProjectsDirName) {
		t.Fatalf("unexpected projects dir")
	}
	if p.ProjectProfilePath("abc") != filepath.Join(root, ProjectsDirName, "abc", "project-profile.yaml") {
		t.Fatalf("unexpected profile path")
	}
}

func TestEnsureDirectoriesCreatesLocalStore(t *testing.T) {
	root := filepath.Join(t.TempDir(), "forgebe")
	p := NewPathsWithRoot(root)
	if err := p.EnsureDirectories(); err != nil {
		t.Fatalf("EnsureDirectories failed: %v", err)
	}
	for _, dir := range []string{p.Root(), p.ProjectsDir(), p.ExportsDir(), p.CacheDir(), p.TmpDir(), p.LogsDir()} {
		if !Exists(dir) {
			t.Fatalf("expected directory to exist: %s", dir)
		}
	}
}
