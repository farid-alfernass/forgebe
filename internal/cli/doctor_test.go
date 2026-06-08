package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestDoctorCmd_NoProject(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", oldHome)

	// Create minimal forgebe dir structure
	if err := os.MkdirAll(tmpHome+"/"+storage.DefaultRootDir+"/projects", 0755); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"doctor"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "ForgeBE Doctor") {
		t.Errorf("expected doctor banner, got:\n%s", output)
	}
	if !strings.Contains(output, "No project profile found") {
		t.Errorf("expected 'No project profile found', got:\n%s", output)
	}
}

func TestDoctorCmd_WithProject(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"doctor", "testproj"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "ForgeBE Doctor") {
		t.Errorf("expected doctor banner, got:\n%s", output)
	}
	if !strings.Contains(output, "testproj") {
		t.Errorf("expected project name in output, got:\n%s", output)
	}
}

func TestDoctorCmd_JSON(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"doctor", "testproj", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "\"valid\"") {
		t.Errorf("expected JSON with 'valid' key, got:\n%s", output)
	}
	if !strings.Contains(output, "\"project\"") {
		t.Errorf("expected JSON with 'project' key, got:\n%s", output)
	}
}
