package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

// newRootCmdWithTempHome creates a root command with a temp home dir so
// storage.NewPaths() / storage.ListProjectIDs() find no projects.
func newRootCmdWithTempHome(t *testing.T, args []string) (*bytes.Buffer, error) {
	tmpHome := t.TempDir()
	// Simulate a minimal ~/.forgebe so NewPaths() / ListProjectIDs can at least
	// traverse the projects directory without error (empty list).
	if err := os.MkdirAll(tmpHome+"/"+storage.DefaultRootDir+"/projects", 0755); err != nil {
		t.Fatal(err)
	}
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	t.Cleanup(func() { os.Setenv("HOME", oldHome) })

	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs(args)
	err := cmd.Execute()
	return &out, err
}

// --- Sync Command Tests ---

func TestSyncCmd_NoProjectFound(t *testing.T) {
	_, err := newRootCmdWithTempHome(t, []string{"sync"})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expectedErr := "sync: no projects found. Run 'forgebe init' first"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestStatusCmd_NoProjectFound(t *testing.T) {
	_, err := newRootCmdWithTempHome(t, []string{"status"})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expectedErr := "status: no projects found. Run 'forgebe init' first"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestWatchCmd_NoProjectFound(t *testing.T) {
	_, err := newRootCmdWithTempHome(t, []string{"watch"})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expectedErr := "watch: no projects found. Run 'forgebe init' first"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

// newRootCmdWithProject sets up a temporary ForgeBE home with a project profile
// and returns a root command runner that can find that project.
func newRootCmdWithProject(t *testing.T, args []string) (*bytes.Buffer, error) {
	tmpHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	t.Cleanup(func() { os.Setenv("HOME", oldHome) })

	// Mimic 'init' by writing a sample profile into the projects dir.
	if err := os.MkdirAll(tmpHome+"/"+storage.DefaultRootDir+"/projects/testproj", 0755); err != nil {
		t.Fatal(err)
	}

	// Also pretend the root of the project exists
	repoRoot := t.TempDir() + "/testproj-repo"
	if err := os.MkdirAll(repoRoot, 0755); err != nil {
		t.Fatal(err)
	}
	// Create .github subdirectory for copilot-instructions.md
	if err := os.MkdirAll(repoRoot+"/.github", 0755); err != nil {
		t.Fatal(err)
	}

	profileContent := `version: "1.0"
metadata:
  profile_id: testproj
  created_at: "2026-06-08T00:00:00Z"
  repo_path: "` + repoRoot + `"
project:
  name: testproj
stack:
  primary_language: go
  framework: ""
areas:
  context:
    claude: true
    cursor: true
    copilot: true
    hermes: true
ai:
  interaction_mode: hybrid
  model_strategy: single
watch:
  recursive: true
  debounce_duration: 1s
`
	if err := os.WriteFile(tmpHome+"/"+storage.DefaultRootDir+"/projects/testproj/project-profile.yaml", []byte(profileContent), 0644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs(args)
	err := cmd.Execute()
	return &out, err
}

func TestWatchCmd_InvalidDebounce(t *testing.T) {
	// First create a real project, then feed it invalid flags
	_, _ = newRootCmdWithProject(t, []string{"version"}) // side effect: set up project infra

	out, err := newRootCmdWithProject(t, []string{"watch", "testproj", "--debounce", "invalid"})
	if err == nil {
		t.Fatal("expected an error for invalid debounce, got nil; output:\n" + out.String())
	}

	if !strings.Contains(err.Error(), "invalid debounce duration") {
		t.Errorf("expected error containing 'invalid debounce duration', got %q", err.Error())
	}
}

func TestWatchCmd_InvalidFullResyncEvery(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"watch", "testproj", "--full-resync-every", "invalid"})
	if err == nil {
		t.Fatal("expected an error for invalid full-resync-every, got nil; output:\n" + out.String())
	}

	if !strings.Contains(err.Error(), "invalid full-resync-every duration") {
		t.Errorf("expected error containing 'invalid full-resync-every duration', got %q", err.Error())
	}
}

func TestSyncCmd_DryRun(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"sync", "testproj", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "Dry run:") {
		t.Errorf("expected dry-run output, got:\n%s", output)
	}
	if !strings.Contains(output, "CLAUDE.md") {
		t.Errorf("expected CLAUDE.md in output, got:\n%s", output)
	}
}

func TestSyncCmd_Force(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"sync", "testproj", "--force"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "Synced") {
		t.Errorf("expected 'Synced' in output, got:\n%s", output)
	}
}

func TestSyncCmd_JSON(t *testing.T) {
	out, err := newRootCmdWithProject(t, []string{"sync", "testproj", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "\"project_id\"") {
		t.Errorf("expected JSON with project_id, got:\n%s", output)
	}
}

func TestStatusCmd_WithProject(t *testing.T) {
	// First sync so there's something to report
	_, _ = newRootCmdWithProject(t, []string{"sync", "testproj"})

	out, err := newRootCmdWithProject(t, []string{"status", "testproj"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "ForgeBE Context Status") {
		t.Errorf("expected status header, got:\n%s", output)
	}
	if !strings.Contains(output, "Summary:") {
		t.Errorf("expected Summary line, got:\n%s", output)
	}
}

func TestStatusCmd_JSON(t *testing.T) {
	// First sync
	_, _ = newRootCmdWithProject(t, []string{"sync", "testproj"})

	out, err := newRootCmdWithProject(t, []string{"status", "testproj", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v; output:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "\"summary\"") {
		t.Errorf("expected JSON with summary, got:\n%s", output)
	}
	if !strings.Contains(output, "\"up_to_date\"") {
		t.Errorf("expected JSON with up_to_date, got:\n%s", output)
	}
}
