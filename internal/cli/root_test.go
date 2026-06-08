package cli

import (
	"bytes"
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd == nil {
		t.Fatal("NewRootCmd returned nil")
	}
	if cmd.Use != "forgebe" {
		t.Errorf("expected Use 'forgebe', got %q", cmd.Use)
	}
	// Verify all subcommands are registered
	subCmds := cmd.Commands()
	if len(subCmds) < 10 {
		t.Errorf("expected at least 10 subcommands, got %d", len(subCmds))
	}
}

func TestNewRootCmd_Help(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("help command failed: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Error("expected help output, got empty string")
	}
}

func TestNewRootCmd_UnknownCommand(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"nonexistent-cmd"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown command, got nil")
	}
}

func TestNewRootCmd_JSONFlag(t *testing.T) {
	cmd := NewRootCmd()
	if err := cmd.ParseFlags([]string{"--json"}); err != nil {
		t.Fatalf("parse --json flag failed: %v", err)
	}
	val, err := cmd.Flags().GetBool("json")
	if err != nil {
		t.Fatalf("get json flag: %v", err)
	}
	if !val {
		t.Error("expected --json to be true")
	}
}

func TestNewRootCmd_VersionSubcommand(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"version"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected version output, got empty")
	}
}
