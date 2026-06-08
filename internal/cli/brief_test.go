package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestBriefCmd_NoProject(t *testing.T) {
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
	cmd.SetArgs([]string{"brief", "claude"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for no project, got nil")
	}
	// The brief command attempts to use the current directory as project ID
	// The exact error message may vary but it should mention profile
	if !strings.Contains(err.Error(), "profile") {
		t.Errorf("expected profile error, got: %q", err.Error())
	}
}

func TestBriefCmd_WithProject(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"brief", "claude", "testproj"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "Project:") {
		t.Errorf("expected brief with Project info, got:\n%s", output)
	}
	if !strings.Contains(output, "testproj") {
		t.Errorf("expected project name in brief, got:\n%s", output)
	}
}
