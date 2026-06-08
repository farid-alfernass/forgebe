package adopt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_EmptyRepoPath(t *testing.T) {
	_, err := Run(Options{RepoPath: ""})
	if err == nil {
		t.Fatal("expected error for empty repo path")
	}
}

func TestRun_NonExistentPath(t *testing.T) {
	_, err := Run(Options{RepoPath: "/nonexistent/path/xyz"})
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}

func TestRun_NotADirectory(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(tmpFile, []byte("hello"), 0644)

	_, err := Run(Options{RepoPath: tmpFile})
	if err == nil {
		t.Fatal("expected error for file instead of directory")
	}
}

func TestRun_DryRun(t *testing.T) {
	// Create a minimal Go project structure
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example\n\ngo 1.22\n"), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "internal"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n"), 0644)

	result, err := Run(Options{RepoPath: tmpDir, DryRun: true})
	if err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !result.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if result.ProjectID == "" {
		t.Error("expected non-empty project ID")
	}
	if result.Language != "go" {
		t.Errorf("expected language=go, got %s", result.Language)
	}
	if len(result.Actions) == 0 {
		t.Error("expected at least one action")
	}
}

func TestRun_ActualAdopt(t *testing.T) {
	// Create a minimal Node.js project structure
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name":"test"}`), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "src"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "src/index.ts"), []byte("console.log('hi')"), 0644)

	// Override FORGEBE_HOME so we don't pollute real storage
	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)

	result, err := Run(Options{RepoPath: tmpDir, DryRun: false})
	if err != nil {
		t.Fatalf("adopt failed: %v", err)
	}

	if result.DryRun {
		t.Error("expected DryRun=false")
	}
	if result.ProjectID == "" {
		t.Error("expected non-empty project ID")
	}

	// Verify profile was saved
	if _, err := os.Stat(result.ProfilePath); os.IsNotExist(err) {
		t.Error("expected profile file to exist after adopt")
	}
}

func TestRun_AlreadyExistsWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n\ngo 1.22\n"), 0644)

	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)

	// First adopt
	_, err := Run(Options{RepoPath: tmpDir, DryRun: false})
	if err != nil {
		t.Fatalf("first adopt failed: %v", err)
	}

	// Second adopt without force should fail
	_, err = Run(Options{RepoPath: tmpDir, DryRun: false, Force: false})
	if err == nil {
		t.Fatal("expected error for duplicate adopt without force")
	}
}

func TestRun_AlreadyExistsWithForce(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n\ngo 1.22\n"), 0644)

	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)

	// First adopt
	_, err := Run(Options{RepoPath: tmpDir, DryRun: false})
	if err != nil {
		t.Fatalf("first adopt failed: %v", err)
	}

	// Second adopt with force should succeed
	result, err := Run(Options{RepoPath: tmpDir, DryRun: false, Force: true})
	if err != nil {
		t.Fatalf("force adopt failed: %v", err)
	}
	if result.ProjectID == "" {
		t.Error("expected non-empty project ID after force adopt")
	}
}
