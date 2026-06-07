package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestSanitizeToolName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"claude", "claude"},
		{"cursor", "cursor"},
		{"copilot", "copilot"},
		{"generic", "generic"},
		{"tool-name", "tool-name"},
		{"tool_name", "tool_name"},
		{"tool123", "tool123"},
		{"tool@#name", "toolname"},
		{"", "tool"},
	}

	for _, tc := range tests {
		got := sanitizeToolName(tc.input)
		if got != tc.expected {
			t.Errorf("sanitizeToolName(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestStatusIcon(t *testing.T) {
	tests := []struct {
		status   string
		expected string
	}{
		{"ok", "✅"},
		{"error", "❌"},
		{"warning", "⚠️"},
		{"unknown", "•"},
	}

	for _, tc := range tests {
		got := statusIcon(tc.status)
		if got != tc.expected {
			t.Errorf("statusIcon(%q) = %q, want %q", tc.status, got, tc.expected)
		}
	}
}

func TestOutputJSON(t *testing.T) {
	cmd := NewRootCmd()
	if err := cmd.ParseFlags([]string{"--json"}); err != nil {
		t.Fatal(err)
	}
	if !OutputJSON(cmd) {
		t.Error("expected OutputJSON to be true")
	}
}

func TestOutputPath(t *testing.T) {
	cmd := newExportSummaryCmd()
	if err := cmd.ParseFlags([]string{"--output", "/tmp/test.md"}); err != nil {
		t.Fatal(err)
	}
	got := OutputPath(cmd)
	if got != "/tmp/test.md" {
		t.Errorf("OutputPath() = %q, want %q", got, "/tmp/test.md")
	}
}

func TestWriteOutput(t *testing.T) {
	rootCmd := NewRootCmd()
	outBuf := new(bytes.Buffer)
	rootCmd.SetOut(outBuf)

	// Test text output (no --json flag)
	if err := WriteOutput(rootCmd, "hello", nil); err != nil {
		t.Fatal(err)
	}
	if outBuf.String() != "hello" {
		t.Errorf("text WriteOutput = %q, want %q", outBuf.String(), "hello")
	}
	outBuf.Reset()

	// Test JSON output
	jsonRoot := NewRootCmd()
	jsonBuf := new(bytes.Buffer)
	jsonRoot.SetOut(jsonBuf)
	if err := jsonRoot.ParseFlags([]string{"--json"}); err != nil {
		t.Fatal(err)
	}

	type testJSON struct {
		Key string `json:"key"`
	}
	if err := WriteOutput(jsonRoot, "", testJSON{Key: "value"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(jsonBuf.String(), `"key": "value"`) {
		t.Errorf("JSON WriteOutput = %q, want JSON containing key value", jsonBuf.String())
	}
}

func TestWriteOutputFile(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := newExportSummaryCmd()
	outBuf := new(bytes.Buffer)
	cmd.SetOut(outBuf)

	if err := cmd.ParseFlags([]string{"--output", filepath.Join(tmpDir, "out.md")}); err != nil {
		t.Fatal(err)
	}

	wrotePath, err := WriteOutputFile(cmd, "test content", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(wrotePath, "out.md") {
		t.Errorf("WriteOutputFile path = %q, want suffix out.md", wrotePath)
	}
	content, _ := os.ReadFile(wrotePath)
	if string(content) != "test content" {
		t.Errorf("WriteOutputFile content = %q, want %q", string(content), "test content")
	}
}

func TestDirExists(t *testing.T) {
	if !dirExists("/") {
		t.Error("/ should exist")
	}
	if dirExists("/nonexistent-dir-12345") {
		t.Error("nonexistent dir should not exist")
	}
}

func TestResolveProjectIDArg_WithID(t *testing.T) {
	got, err := resolveProjectIDArg([]string{"test-proj"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "test-proj" {
		t.Errorf("resolveProjectIDArg = %q, want %q", got, "test-proj")
	}
}

func TestLoadProfileByID_NotFound(t *testing.T) {
	_, err := loadProfileByID("nonexistent_proj")
	if err == nil {
		t.Error("expected error when profile not found")
	}
}

func TestListProjectIDs_NoProjects(t *testing.T) {
	tmpDir := t.TempDir()
	p := storage.NewPathsWithRoot(tmpDir)

	// Need to ensure the projects directory exists first
	if err := os.MkdirAll(p.ProjectsDir(), 0700); err != nil {
		t.Fatal(err)
	}

	ids, err := storage.ListProjectIDs(p)
	if err != nil {
		t.Errorf("ListProjectIDs returned error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty list, got %v", ids)
	}
}
