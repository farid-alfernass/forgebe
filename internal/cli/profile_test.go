package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestProfileShowCmd_NoProject(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", oldHome)

	if err := os.MkdirAll(tmpHome+"/"+storage.DefaultRootDir+"/projects", 0755); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"profile", "show", "nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent project")
	}
	if !strings.Contains(err.Error(), "profile not found") {
		t.Errorf("expected 'profile not found' error, got: %q", err.Error())
	}
}

func TestProfileShowCmd_WithProject(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"profile", "show", "testproj"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "Project:") {
		t.Errorf("expected 'Project:' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "testproj") {
		t.Errorf("expected 'testproj' in output, got:\n%s", output)
	}
}

func TestProfileShowCmd_JSON(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"profile", "show", "testproj", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "\"metadata\"") {
		t.Errorf("expected JSON with metadata key, got:\n%s", output)
	}
	if !strings.Contains(output, "\"stack\"") {
		t.Errorf("expected JSON with stack key, got:\n%s", output)
	}
}
