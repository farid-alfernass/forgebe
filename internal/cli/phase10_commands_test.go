package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// === SCAFFOLD ===

func TestScaffoldCmd_Help(t *testing.T) {
	cmd := newScaffoldCmd()
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("scaffold --help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "scaffold") {
		t.Error("expected 'scaffold' in help output")
	}
}

func TestScaffoldCmd_ListTemplates(t *testing.T) {
	cmd := newScaffoldCmd()
	cmd.SetArgs([]string{"--list-templates"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("scaffold --list-templates failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "go-service") {
		t.Error("expected 'go-service' in template list")
	}
	if !strings.Contains(output, "node-express") {
		t.Error("expected 'node-express' in template list")
	}
	if !strings.Contains(output, "python-fastapi") {
		t.Error("expected 'python-fastapi' in template list")
	}
}

func TestScaffoldCmd_FlagsRegistered(t *testing.T) {
	cmd := newScaffoldCmd()

	flags := []string{"output", "template", "force", "dry-run", "list-templates"}
	for _, name := range flags {
		f := cmd.Flags().Lookup(name)
		if f == nil {
			t.Errorf("expected --%s flag to be registered", name)
		}
	}
}

// === ADOPT ===

func TestAdoptCmd_Help(t *testing.T) {
	cmd := newAdoptCmd()
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("adopt --help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "adopt") {
		t.Error("expected 'adopt' in help output")
	}
}

func TestAdoptCmd_FlagsRegistered(t *testing.T) {
	cmd := newAdoptCmd()

	flags := []string{"force", "dry-run"}
	for _, name := range flags {
		f := cmd.Flags().Lookup(name)
		if f == nil {
			t.Errorf("expected --%s flag to be registered", name)
		}
	}
}

func TestAdoptCmd_NonExistentPath(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)

	cmd := newAdoptCmd()
	cmd.SetArgs([]string{"/nonexistent/path/xyz123"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}

func TestAdoptCmd_EmptyPath(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)

	cmd := newAdoptCmd()
	cmd.SetArgs([]string{""})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for empty path")
	}
}

// === VERIFY ===

func TestVerifyCmd_Help(t *testing.T) {
	cmd := newVerifyCmd()
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("verify --help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "verify") {
		t.Error("expected 'verify' in help output")
	}
}

func TestVerifyCmd_FlagsRegistered(t *testing.T) {
	cmd := newVerifyCmd()

	f := cmd.Flags().Lookup("json")
	if f == nil {
		t.Error("expected --json flag to be registered")
	}
}

func TestVerifyCmd_NonExistentProject(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)

	cmd := newVerifyCmd()
	cmd.SetArgs([]string{"nonexistent-project"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent project")
	}
}

// === CHECK ===

func TestCheckCmd_Help(t *testing.T) {
	cmd := newCheckCmd()
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("check --help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "check") {
		t.Error("expected 'check' in help output")
	}
}

func TestCheckCmd_FlagsRegistered(t *testing.T) {
	cmd := newCheckCmd()

	f := cmd.Flags().Lookup("json")
	if f == nil {
		t.Error("expected --json flag to be registered")
	}
}

func TestCheckCmd_NonExistentProject(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)

	cmd := newCheckCmd()
	cmd.SetArgs([]string{"nonexistent-project"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent project")
	}
}

// === PROMPT ===

func TestPromptCmd_Help(t *testing.T) {
	cmd := newPromptCmd()
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("prompt --help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "prompt") {
		t.Error("expected 'prompt' in help output")
	}
	if !strings.Contains(output, "implement") {
		t.Error("expected 'implement' mode in help")
	}
}

func TestPromptCmd_NoArgs(t *testing.T) {
	cmd := newPromptCmd()

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error when no mode provided")
	}
}

func TestPromptCmd_ValidModes(t *testing.T) {
	modes := []string{"implement", "review", "debug", "plan"}
	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			cmd := newPromptCmd()
			err := cmd.Args(cmd, []string{mode})
			if err != nil {
				t.Errorf("expected %s to be a valid mode, got error: %v", mode, err)
			}
		})
	}
}

func TestPromptCmd_TaskArg(t *testing.T) {
	cmd := newPromptCmd()
	err := cmd.Args(cmd, []string{"implement", "some task with spaces"})
	if err != nil {
		t.Errorf("expected args to be valid, got error: %v", err)
	}
}

func TestPromptCmd_TooManyArgs(t *testing.T) {
	cmd := newPromptCmd()
	err := cmd.Args(cmd, []string{"implement", "task", "extra"})
	if err == nil {
		t.Error("expected error for too many args")
	}
}

func TestPromptCmd_NonExistentProject(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)

	cmd := newPromptCmd()
	cmd.SetArgs([]string{"implement", "my task"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent project")
	}
}

// === ADOPT DRY-RUN with Go project ===
func TestAdoptCmd_DryRunWithProject(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("FORGEBE_HOME", homeDir)
	tempDir := t.TempDir()

	os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module test\n\ngo 1.22\n"), 0644)
	os.WriteFile(filepath.Join(tempDir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644)

	cmd := newAdoptCmd()
	cmd.SetArgs([]string{tempDir, "--dry-run"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("adopt --dry-run failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Dry-run") {
		t.Error("expected 'Dry-run' in output")
	}
}
