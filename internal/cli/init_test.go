package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestInitCmd_NonInteractive(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", oldHome)

	// Create a fake project directory with a go.mod so discovery finds something
	projectDir := t.TempDir()
	if err := os.WriteFile(projectDir+"/go.mod", []byte("module testmod\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(projectDir+"/internal", 0755); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"init", projectDir, "--non-interactive"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "ForgeBE initialized successfully") {
		t.Errorf("expected success message, got:\n%s", output)
	}
	if !strings.Contains(output, "Project ID:") {
		t.Errorf("expected Project ID in output, got:\n%s", output)
	}
	if !strings.Contains(output, "Profile:") {
		t.Errorf("expected Profile path in output, got:\n%s", output)
	}

	// Verify profile was saved
	paths, _ := storage.NewPaths()
	entries, _ := storage.ListProjectIDs(paths)
	if len(entries) == 0 {
		t.Fatal("expected at least one project after init")
	}
}

func TestImportCmd_NoFile(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"import", "/nonexistent/path.zip"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "import") {
		t.Errorf("expected import error, got: %q", err.Error())
	}
}

func TestImportCmd_InvalidZip(t *testing.T) {
	tmpFile := t.TempDir() + "/bad.zip"
	if err := os.WriteFile(tmpFile, []byte("not a zip"), 0644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"import", tmpFile})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid zip")
	}
	if !strings.Contains(err.Error(), "import") {
		t.Errorf("expected import error, got: %q", err.Error())
	}
}
