package contextmgr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestIsRelevantEvent(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		// Manifests
		{"go.mod", true},
		{"package.json", true},
		{"Cargo.toml", true},
		{"build.gradle", true},
		{"requirements.txt", true},
		// Source directories (checked by path segments)
		{"src/main.go", true},
		{"internal/logic.go", true},
		{"cmd/forgebe/main.go", true},
		{"pkg/util/helper.go", true},
		{"api/handler.go", true},
		// Hidden files (except allowed)
		{".git/config", false},
		{".DS_Store", false},
		{".env", true},
		{".gitignore", true},
		// Config files
		{"Dockerfile", true},
		{"tsconfig.json", true},
		{"vite.config.js", true},
		// eslint/prettier
		{"eslint.json", true},
		{"prettier.config.js", true},
		// Non-relevant
		{"README.md", false},
		{"docs/guide.md", false},
		{"node_modules/pkg/index.js", false},
		{"dist/bundle.js", false},
		{"vendor/github.com/foo/bar.go", false},
		{"tmp/debug.log", false},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			if got := IsRelevantEvent(tc.path); got != tc.want {
				t.Errorf("IsRelevantEvent(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestManager_NewManagerErrors(t *testing.T) {
	t.Run("non-existent profile", func(t *testing.T) {
		_, err := NewManager("no-such-project-12345")
		if err == nil {
			t.Fatal("expected error for non-existent projectID, got nil")
		}
		if !strings.Contains(err.Error(), "no such file") {
			t.Errorf("expected 'no such file' error, got: %v", err)
		}
	})
}

func TestManager_SyncAndStatus(t *testing.T) {
	tmpRoot, err := os.MkdirTemp("", "forgebe-mgr-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRoot)

	paths := storage.NewPathsWithRoot(tmpRoot)
	if err := paths.EnsureDirectories(); err != nil {
		t.Fatal(err)
	}

	tmpRepo, err := os.MkdirTemp("", "project-repo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRepo)

	projectID := "test-mgr"
	p := profile.NewSampleProfile()
	p.Metadata.ProfileID = projectID
	p.Metadata.RepoPath = tmpRepo

	store := profile.NewStore(paths)
	if err := store.SaveProfile(p); err != nil {
		t.Fatal(err)
	}

	mgr := &Manager{
		ProjectID: projectID,
		Profile:   &p,
		Paths:     paths,
		Store:     store,
		Tools:     []string{"claude", "hermes"},
	}

	t.Run("sync basic", func(t *testing.T) {
		report, err := mgr.Sync(false, false)
		if err != nil {
			t.Fatalf("Sync failed: %v", err)
		}
		if len(report.Files) != 2 {
			t.Fatalf("expected 2 files, got %d", len(report.Files))
		}

		// Verify files
		claudePath := filepath.Join(tmpRepo, "CLAUDE.md")
		if _, err := os.Stat(claudePath); err != nil {
			t.Errorf("CLAUDE.md missing: %v", err)
		}
		hermesPath := filepath.Join(tmpRepo, "AGENTS.md")
		if _, err := os.Stat(hermesPath); err != nil {
			t.Errorf("AGENTS.md missing: %v", err)
		}
	})

	t.Run("sync idempotent", func(t *testing.T) {
		report, err := mgr.Sync(false, false)
		if err != nil {
			t.Fatal(err)
		}
		if report.ProjectID != projectID {
			t.Error("project ID mismatch")
		}
		if len(report.Files) != 2 {
			t.Errorf("expected 2 files, got %d", len(report.Files))
		}
	})

	t.Run("sync force", func(t *testing.T) {
		report, err := mgr.Sync(true, false)
		if err != nil {
			t.Fatal(err)
		}
		if len(report.Files) != 2 {
			t.Errorf("expected 2 files in force sync, got %d", len(report.Files))
		}
	})

	t.Run("dry run does not write files", func(t *testing.T) {
		// Remove written file
		claudePath := filepath.Join(tmpRepo, "CLAUDE.md")
		os.Remove(claudePath)

		// Dry run should report but not write
		report, err := mgr.Sync(false, true)
		if err != nil {
			t.Fatal(err)
		}
		if len(report.Files) != 2 {
			t.Errorf("expected 2 files in dry run report, got %d", len(report.Files))
		}
		if report.Files[0].SyncedAt != "(dry-run)" {
			t.Errorf("expected synced_at '(dry-run)', got %s", report.Files[0].SyncedAt)
		}
		if _, err := os.Stat(claudePath); !os.IsNotExist(err) {
			t.Error("CLAUDE.md should not exist after dry run")
		}

		// Now actually sync again
		mgr.Sync(false, false)
	})

	t.Run("Status reports correctly", func(t *testing.T) {
		status, err := mgr.Status()
		if err != nil {
			t.Fatalf("Status failed: %v", err)
		}

		if len(status.SyncedFiles) != 2 {
			t.Errorf("expected 2 synced files, got %d", len(status.SyncedFiles))
		}
		if len(status.Outdated) != 0 {
			t.Errorf("expected 0 outdated, got %d: %v", len(status.Outdated), status.Outdated)
		}
		if status.LastSync == "" {
			t.Error("expected last_sync to be set after sync")
		}
	})

	t.Run("Status detects stale file", func(t *testing.T) {
		claudePath := filepath.Join(tmpRepo, "CLAUDE.md")
		if err := os.WriteFile(claudePath, []byte("stale content"), 0644); err != nil {
			t.Fatal(err)
		}

		status, err := mgr.Status()
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, o := range status.Outdated {
			if o == "CLAUDE.md (stale)" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected CLAUDE.md (stale), got outdated: %v", status.Outdated)
		}
	})

	t.Run("Status detects missing files", func(t *testing.T) {
		// Create a manager that expects a file we won't write
		mgrMissing := &Manager{
			ProjectID: projectID,
			Profile:   &p,
			Paths:     paths,
			Store:     store,
			Tools:     []string{"copilot"}, // copilot-instructions.md won't exist
		}

		status, err := mgrMissing.Status()
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, o := range status.Outdated {
			if o == ".github/copilot-instructions.md (missing)" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected copilot-instructions.md (missing), got outdated: %v", status.Outdated)
		}
	})

	t.Run("Status with no sync metadata", func(t *testing.T) {
		// Create a separate project with no sync metadata
		noSyncRoot, err := os.MkdirTemp("", "forgebe-no-sync-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(noSyncRoot)

		noSyncPaths := storage.NewPathsWithRoot(noSyncRoot)
		if err := noSyncPaths.EnsureDirectories(); err != nil {
			t.Fatal(err)
		}

		noSyncRepo, err := os.MkdirTemp("", "no-sync-repo-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(noSyncRepo)

		noSyncID := "no-sync-meta"
		noSyncProfile := profile.NewSampleProfile()
		noSyncProfile.Metadata.ProfileID = noSyncID
		noSyncProfile.Metadata.RepoPath = noSyncRepo

		noSyncStore := profile.NewStore(noSyncPaths)
		if err := noSyncStore.SaveProfile(noSyncProfile); err != nil {
			t.Fatal(err)
		}

		mgrNoSync := &Manager{
			ProjectID: noSyncID,
			Profile:   &noSyncProfile,
			Paths:     noSyncPaths,
			Store:     noSyncStore,
			Tools:     []string{"claude"},
		}

		status, err := mgrNoSync.Status()
		if err != nil {
			t.Fatal(err)
		}
		if status.LastSync != "" {
			t.Errorf("expected LastSync empty, got %s", status.LastSync)
		}
		if len(status.SyncedFiles) != 0 {
			t.Errorf("expected 0 synced files, got %d", len(status.SyncedFiles))
		}
		if len(status.Outdated) != 1 || status.Outdated[0] != "CLAUDE.md (missing)" {
			t.Errorf("expected 1 missing file, got %v", status.Outdated)
		}
	})
}

func TestManager_SyncEmptyRepoPath(t *testing.T) {
	tmpRoot, err := os.MkdirTemp("", "forgebe-nopath-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRoot)

	paths := storage.NewPathsWithRoot(tmpRoot)
	if err := paths.EnsureDirectories(); err != nil {
		t.Fatal(err)
	}

	projectID := "no-repo-path"
	p := profile.NewSampleProfile()
	p.Metadata.ProfileID = projectID
	p.Metadata.RepoPath = "" // empty!

	store := profile.NewStore(paths)
	if err := store.SaveProfile(p); err != nil {
		t.Fatal(err)
	}

	mgr, err := NewManagerWithPaths(projectID, paths)
	if err != nil {
		t.Fatal(err)
	}

	_, err = mgr.Sync(false, false)
	if err == nil {
		t.Fatal("expected error for empty repo_path, got nil")
	}
	if !strings.Contains(err.Error(), "repo_path not set") {
		t.Errorf("expected 'repo_path not set' error, got: %v", err)
	}
}
