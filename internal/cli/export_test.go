package cli

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestExportSummaryCmd_NoProject(t *testing.T) {
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
	cmd.SetArgs([]string{"export", "summary"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "export summary: no projects found") {
		t.Errorf("expected error 'export summary: no projects found', got %q", err.Error())
	}
}

func TestExportSummaryCmd_WithProject(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"export", "summary", "testproj"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	if !strings.Contains(out.String(), "Project Summary") {
		t.Errorf("expected project summary, got:\n%s", out.String())
	}
}

func TestExportSummaryCmd_JSON(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"export", "summary", "testproj", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	if !strings.Contains(out.String(), "summary") {
		t.Errorf("expected JSON with summary, got:\n%s", out.String())
	}
}

func TestExportSummaryCmd_ToFile(t *testing.T) {
	outputFile := filepath.Join(t.TempDir(), "summary.md")
	out, err := newRootCmdWithProject(t, []string{"export", "summary", "testproj", "-o", outputFile})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	if !strings.Contains(out.String(), "Written to:") {
		t.Errorf("expected 'Written to:' in output, got:\n%s", out.String())
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("output file not created")
	}
}

func TestExportAdapterCmd_NoProject(t *testing.T) {
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
	cmd.SetArgs([]string{"export", "adapter", "claude"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "profile store") {
		t.Errorf("expected profile error, got %q", err.Error())
	}
}

func TestExportAdapterCmd_WithProject(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"export", "adapter", "claude", "testproj"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	if !strings.Contains(out.String(), "Exported:") {
		t.Errorf("expected 'Exported:' in output, got:\n%s", out.String())
	}

	match := regexp.MustCompile(`Exported: (.+)`).FindStringSubmatch(out.String())
	if len(match) != 2 {
		t.Fatalf("could not extract exported path from: %s", out.String())
	}

	if _, err := os.Stat(match[1]); os.IsNotExist(err) {
		t.Fatalf("exported file not found at %s", match[1])
	}
}

func TestExportBundleCmd_NoProject(t *testing.T) {
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
	cmd.SetArgs([]string{"export", "bundle"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "export bundle: no projects found") {
		t.Errorf("expected error 'export bundle: no projects found', got %q", err.Error())
	}
}

func TestExportBundleCmd_WithProject(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"export", "bundle", "testproj"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	if !strings.Contains(out.String(), "Exported bundle:") {
		t.Errorf("expected 'Exported bundle:' in output, got:\n%s", out.String())
	}

	match := regexp.MustCompile(`Exported bundle: (.+)`).FindStringSubmatch(out.String())
	if len(match) != 2 {
		t.Fatalf("could not extract bundle path from: %s", out.String())
	}

	if _, err := os.Stat(match[1]); os.IsNotExist(err) {
		t.Fatalf("exported bundle not found at %s", match[1])
	}

	r, err := zip.OpenReader(match[1])
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	defer r.Close()

	if len(r.File) == 0 {
		t.Errorf("bundle zip is empty, expected profile files")
	}
}
