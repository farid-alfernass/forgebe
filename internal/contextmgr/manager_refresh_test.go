package contextmgr

import (
	"os"
	"testing"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestManager_RefreshProfile(t *testing.T) {
	tmpRoot := t.TempDir()
	paths := storage.NewPathsWithRoot(tmpRoot)
	if err := paths.EnsureDirectories(); err != nil {
		t.Fatal(err)
	}

	projectID := "refresh-test"
	p := profile.NewSampleProfile()
	p.Metadata.ProfileID = projectID
	p.Metadata.RepoPath = t.TempDir()
	p.Project.Name = "original"

	store := profile.NewStore(paths)
	if err := store.SaveProfile(p); err != nil {
		t.Fatal(err)
	}

	mgr := &Manager{
		ProjectID: projectID,
		Profile:   &p,
		Paths:     paths,
		Store:     store,
		Tools:     DefaultTools,
	}

	// Verify initial state
	if mgr.Profile.Project.Name != "original" {
		t.Fatalf("expected 'original', got %q", mgr.Profile.Project.Name)
	}

	// Change profile on disk
	p.Project.Name = "updated"
	if err := store.SaveProfile(p); err != nil {
		t.Fatal(err)
	}

	// RefreshProfile should pick up the change
	if err := mgr.RefreshProfile(); err != nil {
		t.Fatalf("RefreshProfile failed: %v", err)
	}
	if mgr.Profile.Project.Name != "updated" {
		t.Errorf("expected 'updated' after refresh, got %q", mgr.Profile.Project.Name)
	}
}

func TestManager_RefreshProfile_Error(t *testing.T) {
	tmpRoot := t.TempDir()
	paths := storage.NewPathsWithRoot(tmpRoot)
	if err := paths.EnsureDirectories(); err != nil {
		t.Fatal(err)
	}

	// Don't save a profile — RefreshProfile should fail
	p := profile.NewSampleProfile()
	p.Metadata.ProfileID = "nonexistent-refresh"

	mgr := &Manager{
		ProjectID: "nonexistent-refresh",
		Profile:   &p,
		Paths:     paths,
		Store:     profile.NewStore(paths),
		Tools:     DefaultTools,
	}

	err := mgr.RefreshProfile()
	if err == nil {
		t.Fatal("expected error from RefreshProfile when profile doesn't exist")
	}
}

func TestManager_GetWatchConfigWithValues(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Watch = profile.Watch{
		Recursive:        true,
		DebounceDuration: 5 * time.Second,
		FullResyncEvery:  60 * time.Second,
		IgnorePatterns:   []string{"custom_ignore"},
		MatchPatterns:    []string{"*.go"},
	}
	mgr := &Manager{Profile: &p}

	cfg := mgr.GetWatchConfig()
	if cfg.DebounceDuration != 5*time.Second {
		t.Errorf("expected 5s debounce, got %v", cfg.DebounceDuration)
	}
	if cfg.FullResyncEvery != 60*time.Second {
		t.Errorf("expected 60s full resync, got %v", cfg.FullResyncEvery)
	}
	if len(cfg.IgnorePatterns) != 1 || cfg.IgnorePatterns[0] != "custom_ignore" {
		t.Errorf("expected custom_ignore, got %v", cfg.IgnorePatterns)
	}
	if len(cfg.MatchPatterns) != 1 || cfg.MatchPatterns[0] != "*.go" {
		t.Errorf("expected *.go match, got %v", cfg.MatchPatterns)
	}
}

func TestManager_SyncWriteError(t *testing.T) {
	tmpRoot := t.TempDir()
	paths := storage.NewPathsWithRoot(tmpRoot)
	if err := paths.EnsureDirectories(); err != nil {
		t.Fatal(err)
	}

	// Use a non-writable directory as repo path
	readonlyRepo := t.TempDir()
	if err := os.Chmod(readonlyRepo, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readonlyRepo, 0755) // cleanup

	projectID := "write-error-test"
	p := profile.NewSampleProfile()
	p.Metadata.ProfileID = projectID
	p.Metadata.RepoPath = readonlyRepo

	store := profile.NewStore(paths)
	if err := store.SaveProfile(p); err != nil {
		t.Fatal(err)
	}

	mgr := &Manager{
		ProjectID: projectID,
		Profile:   &p,
		Paths:     paths,
		Store:     store,
		Tools:     []string{"claude"},
	}

	_, err := mgr.Sync(true, false)
	if err == nil {
		t.Fatal("expected error from Sync when repo is not writable")
	}
}

func TestIsRelevantEvent_ConfigFiles(t *testing.T) {
	// Test config files that should be relevant
	configCases := []struct {
		path string
		want bool
	}{
		{"docker-compose.yml", true},
		{"docker-compose.yaml", true},
		// .env.* files are hidden (start with .) and not exact ".env" or ".gitignore"
		// so they are filtered out by the hidden file check
		{".env.example", false},
		{".env.local", false},
		{".env.development", false},
		{".env.test", false},
		{".env.production", false},
		{"tsconfig.json", true},
		{"jsconfig.json", true},
		{"webpack.config.js", true},
		{"next.config.js", true},
		{"nuxt.config.js", true},
		{"jest.config.js", true},
		{"prettier.config.js", true},
		{"eslint.config.js", true},
		// Additional manifests
		{"Package.swift", true},
		{"build.gradle.kts", true},
		{"Cargo.lock", true},
		{"package-lock.json", true},
		{"pnpm-lock.yaml", true},
		{"setup.py", true},
		{"Pipfile", true},
		{"Pipfile.lock", true},
		{"composer.json", true},
		{"composer.lock", true},
		{"Gemfile", true},
		{"Gemfile.lock", true},
		{"build.sbt", true},
		{"project.clj", true},
		{"Mixfile", true},
		{"mix.lock", true},
		{"CMakeLists.txt", true},
		{"configure.ac", true},
		{"configure.in", true},
		// Source paths
		{"lib/utils.rb", true},
		{"app/models/user.rb", true},
		{"services/auth/main.go", true},
		// Skip paths
		{"build/output.js", false},
		{"out/index.html", false},
	}

	for _, tc := range configCases {
		t.Run(tc.path, func(t *testing.T) {
			if got := IsRelevantEvent(tc.path); got != tc.want {
				t.Errorf("IsRelevantEvent(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestIsRelevantEventWithPatterns_IgnoreContains(t *testing.T) {
	// Test that ignore pattern with 'contains' works
	got := IsRelevantEventWithPatterns("some/generated-code/file.go", []string{"generated"}, nil)
	if got {
		t.Error("expected path with 'generated' component to be ignored")
	}

	// Manifest should still pass even with ignore patterns
	got = IsRelevantEventWithPatterns("go.mod", []string{"generated"}, nil)
	if !got {
		t.Error("manifest go.mod should always be relevant even with ignore patterns")
	}
}

func TestIsRelevantEventWithPatterns_MatchContainsPath(t *testing.T) {
	// Match pattern checks cleanPath contains
	got := IsRelevantEventWithPatterns("configs/production/app.yaml", nil, []string{"production"})
	if !got {
		t.Error("expected match on path containing 'production'")
	}
}
