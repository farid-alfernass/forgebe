package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewPaths(t *testing.T) {
	p, err := NewPaths()
	if err != nil {
		t.Fatal(err)
	}
	if p.Root() == "" {
		t.Error("expected non-empty root")
	}
	expected := filepath.Join(os.Getenv("HOME"), DefaultRootDir)
	if p.Root() != expected {
		t.Errorf("Root() = %q, want %q", p.Root(), expected)
	}
}

func TestPathsAllMethods(t *testing.T) {
	root := t.TempDir()
	p := NewPathsWithRoot(root)

	paths := map[string]string{
		"Root":                 p.Root(),
		"ProjectsDir":          p.ProjectsDir(),
		"ExportsDir":           p.ExportsDir(),
		"CacheDir":             p.CacheDir(),
		"TmpDir":               p.TmpDir(),
		"LogsDir":              p.LogsDir(),
		"ProjectDir":           p.ProjectDir("test"),
		"ProjectProfilePath":   p.ProjectProfilePath("test"),
		"ProjectDiscoveryPath": p.ProjectDiscoveryPath("test"),
		"ProjectContractPath":  p.ProjectContractPath("test"),
		"ProjectSummaryPath":   p.ProjectSummaryPath("test"),
		"ProjectMetadataPath":  p.ProjectMetadataPath("test"),
	}

	for name, val := range paths {
		if val == "" {
			t.Errorf("%s returned empty string", name)
		}
	}
}

func TestEnsureBaseDirs(t *testing.T) {
	root := t.TempDir()
	p := NewPathsWithRoot(root)

	if err := p.EnsureBaseDirs(); err != nil {
		t.Fatalf("EnsureBaseDirs failed: %v", err)
	}

	baseDirs := []string{
		p.Root(),
		p.ProjectsDir(),
		p.ExportsDir(),
		p.CacheDir(),
		p.TmpDir(),
		p.LogsDir(),
	}

	for _, d := range baseDirs {
		fi, err := os.Stat(d)
		if err != nil {
			t.Errorf("directory %s should exist: %v", d, err)
		} else if !fi.IsDir() {
			t.Errorf("%s is not a directory", d)
		}
	}
}

func TestEnsureProjectDir(t *testing.T) {
	root := t.TempDir()
	p := NewPathsWithRoot(root)
	if err := p.EnsureBaseDirs(); err != nil {
		t.Fatal(err)
	}

	if err := p.EnsureProjectDir("test-proj"); err != nil {
		t.Fatalf("EnsureProjectDir failed: %v", err)
	}

	fi, err := os.Stat(p.ProjectDir("test-proj"))
	if err != nil || !fi.IsDir() {
		t.Errorf("project directory should exist")
	}
}

func TestListProjectIDs(t *testing.T) {
	root := t.TempDir()
	p := NewPathsWithRoot(root)
	if err := p.EnsureBaseDirs(); err != nil {
		t.Fatal(err)
	}

	// Empty list
	ids, err := ListProjectIDs(p)
	if err != nil {
		t.Fatalf("ListProjectIDs failed: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty, got %v", ids)
	}

	// Create two projects
	if err := p.EnsureProjectDir("proj-a"); err != nil {
		t.Fatal(err)
	}
	if err := p.EnsureProjectDir("proj-b"); err != nil {
		t.Fatal(err)
	}

	ids, err = ListProjectIDs(p)
	if err != nil {
		t.Fatalf("ListProjectIDs failed: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 projects, got %d", len(ids))
	}
	// Should be sorted
	if ids[0] != "proj-a" || ids[1] != "proj-b" {
		t.Errorf("expected [proj-a, proj-b], got %v", ids)
	}
}

func TestRemove(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(f, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Remove(f); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Error("file should have been removed")
	}
}

func TestRemove_NonExistent(t *testing.T) {
	// Remove on non-existent path should not error (os.RemoveAll semantics)
	err := Remove("/nonexistent-path-xyz-12345")
	if err != nil {
		t.Errorf("Remove non-existent should not error, got: %v", err)
	}
}

func TestAtomicWrite_NewFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "new-file.txt")

	if err := AtomicWrite(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "hello" {
		t.Errorf("content = %q, want %q", string(content), "hello")
	}
}

func TestAtomicWrite_Overwrite(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "overwrite.txt")

	if err := AtomicWrite(f, []byte("first"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := AtomicWrite(f, []byte("second"), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "second" {
		t.Errorf("content = %q, want %q", string(content), "second")
	}
}

func TestReadFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "read.txt")
	if err := os.WriteFile(f, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	data, err := ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "data" {
		t.Errorf("ReadFile = %q, want %q", string(data), "data")
	}
}

func TestReadFile_NotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent-xyz-12345")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestExists(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "exists.txt")
	if err := os.WriteFile(f, []byte("yes"), 0644); err != nil {
		t.Fatal(err)
	}

	if !Exists(f) {
		t.Error("Exists should return true for existing file")
	}
	if Exists("/nonexistent-xyz-12345") {
		t.Error("Exists should return false for non-existent path")
	}
}

func TestEnsureDirectories(t *testing.T) {
	root := t.TempDir()
	p := NewPathsWithRoot(root)

	if err := p.EnsureDirectories(); err != nil {
		t.Fatal(err)
	}

	dirs := []string{
		p.Root(),
		p.ProjectsDir(),
		p.ExportsDir(),
		p.CacheDir(),
		p.TmpDir(),
		p.LogsDir(),
	}

	for _, d := range dirs {
		fi, err := os.Stat(d)
		if err != nil || !fi.IsDir() {
			t.Errorf("EnsureDirectories: %s should exist as directory", d)
		}
	}
}
