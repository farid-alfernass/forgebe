package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteOutputFile_Stdout(t *testing.T) {
	cmd := newExportSummaryCmd()
	outBuf := new(bytes.Buffer)
	cmd.SetOut(outBuf)
	cmd.SetErr(outBuf)

	if err := cmd.ParseFlags([]string{"--output", "-"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	_, err := WriteOutputFile(cmd, "test content", "default.txt")
	if err != nil {
		t.Fatalf("WriteOutputFile failed: %v", err)
	}

	if outBuf.String() != "test content" {
		t.Errorf("expected stdout output, got: %q", outBuf.String())
	}
}

func TestWriteOutputFile_MkdirError(t *testing.T) {
	// Create a temp dir then make it read-only so MkdirAll fails
	tmpDir := t.TempDir()
	readonlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readonlyDir, 0555); err != nil {
		t.Fatalf("mkdir readonly: %v", err)
	}

	cmd := newExportSummaryCmd()
	outPath := filepath.Join(readonlyDir, "subdir", "out.txt")

	if err := cmd.ParseFlags([]string{"--output", outPath}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	_, err := WriteOutputFile(cmd, "test content", "default.txt")
	if err == nil {
		t.Fatal("expected error from MkdirAll, got nil")
	}
	if !strings.Contains(err.Error(), "write output:") {
		t.Errorf("unexpected error: %v", err)
	}

	// Cleanup: remove readonly permission so os.RemoveAll works in TempDir cleanup
	if err := os.Chmod(readonlyDir, 0755); err != nil {
		t.Errorf("chmod cleanup: %v", err)
	}
}

func TestWriteOutputFile_WriteError(t *testing.T) {
	// Point to a directory instead of a file - WriteFile will fail
	tmpDir := t.TempDir()
	isDir := filepath.Join(tmpDir, "isdir")
	if err := os.MkdirAll(isDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cmd := newExportSummaryCmd()
	if err := cmd.ParseFlags([]string{"--output", isDir}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	_, err := WriteOutputFile(cmd, "test content", "default.txt")
	if err == nil {
		t.Fatal("expected error from WriteFile, got nil")
	}
	if !strings.Contains(err.Error(), "write output:") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDir(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "."},
		{"file.txt", "."},
		{"/path/to/file.txt", "/path/to"},
		{"subdir/file.txt", "subdir"},
		{"C:\\path\\file.txt", `C:\path`},
	}
	for _, tc := range tests {
		got := dir(tc.input)
		if got != tc.want {
			t.Errorf("dir(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
