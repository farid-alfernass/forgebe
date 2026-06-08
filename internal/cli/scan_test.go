package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanCmd(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a dummy go.mod to trigger go detection
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newScanCmd()
	outBuf := new(bytes.Buffer)
	cmd.SetOut(outBuf)
	cmd.SetArgs([]string{tmpDir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("scan cmd failed: %v", err)
	}

	output := outBuf.String()
	if !strings.Contains(output, "language: go") {
		t.Errorf("scan output missing go detection: %s", output)
	}
}

func TestScanCmd_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newScanCmd()
	outBuf := new(bytes.Buffer)
	cmd.SetOut(outBuf)
	// We need to set the global json flag on a root command or just ensure OutputJSON works
	// Root command has the flag, but we can also set it manually for the test
	cmd.Flags().Bool("json", true, "")
	cmd.SetArgs([]string{"--json", tmpDir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("scan cmd --json failed: %v", err)
	}

	output := outBuf.String()
	if !strings.Contains(output, "\"language\": \"go\"") {
		t.Errorf("scan JSON output missing go detection: %s", output)
	}
}

func TestScanCmd_NonExistent(t *testing.T) {
	cmd := newScanCmd()
	cmd.SetArgs([]string{"/non/existent/path/for/forgebe/test"})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error for non-existent path")
	}
}
