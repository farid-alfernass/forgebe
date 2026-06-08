package contextmgr

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/adapters"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

// SyncedFile tracks a generated adapter file.
type SyncedFile struct {
	Tool     string `json:"tool"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	SyncedAt string `json:"synced_at"`
}

// Report is the result of a sync operation.
type Report struct {
	ProjectID string       `json:"project_id"`
	RepoPath  string       `json:"repo_path"`
	Files     []SyncedFile `json:"files"`
	Duration  string       `json:"duration"`
	Error     string       `json:"error,omitempty"`
}

// StatusReport represents the current state of managed context files.
type StatusReport struct {
	ProjectID    string        `json:"project_id"`
	RepoPath     string        `json:"repo_path"`
	ProfileAge   string        `json:"profile_age"`
	LastSync     string        `json:"last_sync"`
	SyncedFiles  []SyncedFile  `json:"synced_files"`
	Summary      StatusSummary `json:"summary"`
	MissingFiles []FileStatus  `json:"missing_files"`
	ChangedFiles []FileStatus  `json:"changed_files"`
	Outdated     []string      `json:"outdated"`
}

// StatusSummary gives compact counts for status UX and JSON consumers.
type StatusSummary struct {
	Total    int `json:"total"`
	Synced   int `json:"synced"`
	Missing  int `json:"missing"`
	Changed  int `json:"changed"`
	UpToDate int `json:"up_to_date"`
}

// FileStatus describes a managed file that needs attention.
type FileStatus struct {
	Tool         string `json:"tool"`
	Filename     string `json:"filename"`
	Reason       string `json:"reason"`
	ExpectedHash string `json:"expected_hash,omitempty"`
	ActualHash   string `json:"actual_hash,omitempty"`
}

// DefaultTools is the list of tools to manage by default.
var DefaultTools = []string{"claude", "cursor", "copilot", "hermes", "generic"}

// Manager handles syncing and watching AI context files for a project.
type Manager struct {
	ProjectID string
	Profile   *profile.ProjectProfile
	Paths     *storage.Paths
	Store     *profile.Store
	Tools     []string
}

// NewManager creates a Manager for the given projectID.
func NewManager(projectID string) (*Manager, error) {
	paths, err := storage.NewPaths()
	if err != nil {
		return nil, fmt.Errorf("contextmgr: paths: %w", err)
	}
	return NewManagerWithPaths(projectID, paths)
}

// NewManagerWithPaths creates a Manager with a custom storage paths (for testing).
func NewManagerWithPaths(projectID string, paths *storage.Paths) (*Manager, error) {
	store := profile.NewStore(paths)
	p, err := store.LoadProfile(projectID)
	if err != nil {
		return nil, fmt.Errorf("contextmgr: load profile: %w", err)
	}
	return &Manager{
		ProjectID: projectID,
		Profile:   p,
		Paths:     paths,
		Store:     store,
		Tools:     DefaultTools,
	}, nil
}

// Sync performs a one-shot update of all AI context files in the project root.
func (m *Manager) Sync(force bool, dryRun bool) (*Report, error) {
	start := time.Now()

	repoPath := m.Profile.Metadata.RepoPath
	if repoPath == "" {
		return nil, fmt.Errorf("contextmgr: repo_path not set in profile; run 'forgebe init' first")
	}

	var files []SyncedFile

	for _, tool := range m.Tools {
		content, filename, err := adapters.Render(tool, *m.Profile)
		if err != nil {
			// Skip unsupported tools silently
			continue
		}

		destPath := filepath.Join(repoPath, filename)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

		if dryRun {
			files = append(files, SyncedFile{
				Tool:     tool,
				Filename: filename,
				Hash:     hash,
				SyncedAt: "(dry-run)",
			})
			continue
		}

		// Skip if content hasn't changed and not forced
		if !force {
			if _, err := os.Stat(destPath); err == nil {
				existing, err := os.ReadFile(destPath)
				if err == nil {
					existingHash := fmt.Sprintf("%x", sha256.Sum256(existing))
					if existingHash == hash {
						files = append(files, SyncedFile{
							Tool:     tool,
							Filename: filename,
							Hash:     hash,
							SyncedAt: start.Format(time.RFC3339),
						})
						continue
					}
				}
			}
		}

		if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("contextmgr: write %s: %w", filename, err)
		}

		files = append(files, SyncedFile{
			Tool:     tool,
			Filename: filename,
			Hash:     hash,
			SyncedAt: start.Format(time.RFC3339),
		})
	}

	if dryRun {
		return &Report{
			ProjectID: m.ProjectID,
			RepoPath:  repoPath,
			Files:     files,
			Duration:  time.Since(start).Round(time.Millisecond).String(),
		}, nil
	}

	// Update profile timestamp
	m.Profile.Metadata.UpdatedAt = start
	if err := m.Store.SaveProfile(*m.Profile); err != nil {
		return nil, fmt.Errorf("contextmgr: save profile: %w", err)
	}

	// Save sync metadata for tracking
	syncMeta := map[string]interface{}{
		"last_sync": start.Format(time.RFC3339),
		"files":     files,
	}
	if err := m.Store.SaveMetadata(m.ProjectID, syncMeta); err != nil {
		return nil, fmt.Errorf("contextmgr: save metadata: %w", err)
	}

	return &Report{
		ProjectID: m.ProjectID,
		RepoPath:  repoPath,
		Files:     files,
		Duration:  time.Since(start).Round(time.Millisecond).String(),
	}, nil
}

// Status returns the current state of managed context files.
func (m *Manager) Status() (*StatusReport, error) {
	repoPath := m.Profile.Metadata.RepoPath
	profileAge := time.Since(m.Profile.Metadata.CreatedAt).Round(time.Second).String()

	// Try to read last sync from metadata
	var lastSync string
	metaPath := m.Paths.ProjectMetadataPath(m.ProjectID)
	if data, err := os.ReadFile(metaPath); err == nil {
		// Simple parse — just extract last_sync if available
		var meta map[string]interface{}
		if err := json.Unmarshal(data, &meta); err == nil {
			if ls, ok := meta["last_sync"].(string); ok {
				lastSync = ls
			}
		}
	}

	var syncedFiles []SyncedFile
	var missingFiles []FileStatus
	var changedFiles []FileStatus
	var outdated []string
	summary := StatusSummary{}

	for _, tool := range m.Tools {
		content, filename, err := adapters.Render(tool, *m.Profile)
		if err != nil {
			continue
		}
		summary.Total++

		expectedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
		destPath := filepath.Join(repoPath, filename)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			summary.Missing++
			missingFiles = append(missingFiles, FileStatus{
				Tool:         tool,
				Filename:     filename,
				Reason:       "missing",
				ExpectedHash: expectedHash,
			})
			outdated = append(outdated, filename+" (missing)")
			continue
		}

		existing, err := os.ReadFile(destPath)
		if err != nil {
			summary.Changed++
			changedFiles = append(changedFiles, FileStatus{
				Tool:         tool,
				Filename:     filename,
				Reason:       "unreadable",
				ExpectedHash: expectedHash,
			})
			outdated = append(outdated, filename+" (unreadable)")
			continue
		}

		existingHash := fmt.Sprintf("%x", sha256.Sum256(existing))
		syncedFiles = append(syncedFiles, SyncedFile{
			Tool:     tool,
			Filename: filename,
			Hash:     existingHash,
		})
		if existingHash != expectedHash {
			summary.Changed++
			changedFiles = append(changedFiles, FileStatus{
				Tool:         tool,
				Filename:     filename,
				Reason:       "changed",
				ExpectedHash: expectedHash,
				ActualHash:   existingHash,
			})
			outdated = append(outdated, filename+" (changed)")
			continue
		}
		summary.UpToDate++
	}
	summary.Synced = len(syncedFiles)

	return &StatusReport{
		ProjectID:    m.ProjectID,
		RepoPath:     repoPath,
		ProfileAge:   profileAge,
		LastSync:     lastSync,
		SyncedFiles:  syncedFiles,
		Summary:      summary,
		MissingFiles: missingFiles,
		ChangedFiles: changedFiles,
		Outdated:     outdated,
	}, nil
}

// RefreshProfile reloads the profile from disk, picking up external changes.
func (m *Manager) RefreshProfile() error {
	p, err := m.Store.LoadProfile(m.ProjectID)
	if err != nil {
		return err
	}
	m.Profile = p
	return nil
}

// IsRelevantEvent determines if a file change is relevant for triggering a sync.
func IsRelevantEvent(path string) bool {
	// Normalize path
	cleanPath := filepath.Clean(path)
	base := filepath.Base(cleanPath)

	// Ignore hidden files and directories (starting with .) except for specific config files
	if strings.HasPrefix(base, ".") && base != ".gitignore" && base != ".env" {
		return false
	}

	// Manifest files that affect stack detection
	manifestFiles := map[string]bool{
		"go.mod":            true,
		"go.sum":            true,
		"Package.swift":     true,
		"pom.xml":           true,
		"build.gradle":      true,
		"build.gradle.kts":  true,
		"Cargo.toml":        true,
		"Cargo.lock":        true,
		"package.json":      true,
		"package-lock.json": true,
		"yarn.lock":         true,
		"pnpm-lock.yaml":    true,
		"requirements.txt":  true,
		"setup.py":          true,
		"pyproject.toml":    true,
		"Pipfile":           true,
		"Pipfile.lock":      true,
		"composer.json":     true,
		"composer.lock":     true,
		"Gemfile":           true,
		"Gemfile.lock":      true,
		"build.sbt":         true,
		"project.clj":       true,
		"Mixfile":           true,
		"mix.lock":          true,
		"Makefile":          true,
		"makefile":          true,
		"CMakeLists.txt":    true,
		"configure.ac":      true,
		"configure.in":      true,
	}

	if manifestFiles[base] {
		return true
	}

	// Check parts of the path for source directories
	// Skip generated/build directories first
	skipDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		"dist":         true,
		"build":        true,
		"out":          true,
		".git":         true,
	}
	parts := strings.Split(cleanPath, string(os.PathSeparator))
	for _, part := range parts {
		if skipDirs[part] {
			return false
		}
	}
	sourceDirs := map[string]bool{
		"src":      true,
		"lib":      true,
		"app":      true,
		"cmd":      true,
		"internal": true,
		"pkg":      true,
		"api":      true,
		"services": true,
	}
	for _, part := range parts {
		if sourceDirs[part] {
			return true
		}
	}

	// Watch for changes in common config files
	configFiles := map[string]bool{
		"Dockerfile":          true,
		"docker-compose.yml":  true,
		"docker-compose.yaml": true,
		".gitignore":          true,
		".env":                true,
		".env.example":        true,
		".env.local":          true,
		".env.development":    true,
		".env.test":           true,
		".env.production":     true,
		"tsconfig.json":       true,
		"jsconfig.json":       true,
		"webpack.config.js":   true,
		"vite.config.js":      true,
		"next.config.js":      true,
		"nuxt.config.js":      true,
		"jest.config.js":      true,
	}
	if configFiles[base] {
		return true
	}
	// Check for eslint/prettier with extensions
	if strings.HasPrefix(base, "eslint.") || strings.HasPrefix(base, "prettier.") {
		return true
	}

	return false
}

// IsRelevantEventWithPatterns checks if a path is relevant, using configurable
// ignore and match patterns. If matchPatterns is non-empty, only files whose
// base name (or directory component) matches one of the patterns are considered
// relevant. If ignorePatterns is non-empty, directories whose base name matches
// any pattern cause the path to be skipped.
func IsRelevantEventWithPatterns(path string, ignorePatterns, matchPatterns []string) bool {
	cleanPath := filepath.Clean(path)
	base := filepath.Base(cleanPath)

	// First check the hardcoded manifest/config files (always relevant)
	manifestCheck := map[string]bool{
		"go.mod": true, "go.sum": true, "package.json": true, "Cargo.toml": true,
		"pyproject.toml": true, "Dockerfile": true, ".gitignore": true, ".env": true,
		"Makefile": true, "makefile": true,
	}
	if manifestCheck[base] {
		return true
	}

	// Check ignore patterns on path components
	if len(ignorePatterns) > 0 {
		parts := strings.Split(cleanPath, string(os.PathSeparator))
		for _, part := range parts {
			for _, ign := range ignorePatterns {
				if strings.EqualFold(part, ign) || strings.Contains(part, ign) {
					return false
				}
			}
		}
	}

	// If match patterns are specified, the file must match at least one
	if len(matchPatterns) > 0 {
		for _, pat := range matchPatterns {
			if strings.EqualFold(base, pat) || strings.Contains(cleanPath, pat) {
				return true
			}
		}
		// Also allow source directory patterns
		sourceDirs := map[string]bool{
			"src": true, "lib": true, "app": true, "cmd": true,
			"internal": true, "pkg": true, "api": true, "services": true,
		}
		parts := strings.Split(cleanPath, string(os.PathSeparator))
		for _, part := range parts {
			if sourceDirs[part] {
				return true
			}
		}
		return false
	}

	// No match patterns: fall through to standard IsRelevantEvent logic
	return IsRelevantEvent(path)
}

// ShouldWatchDir checks if a directory name should be watched, based on
// the configured ignore patterns.
func ShouldWatchDir(dirName string, ignorePatterns []string) bool {
	// Always ignore .git
	if dirName == ".git" {
		return false
	}

	// Skip hidden directories
	if strings.HasPrefix(dirName, ".") && dirName != "." {
		return false
	}

	for _, ign := range ignorePatterns {
		if strings.EqualFold(dirName, ign) || strings.Contains(dirName, ign) {
			return false
		}
	}
	return true
}

// GetWatchConfig returns the watch configuration from the profile.
// If no config is set, it returns sensible defaults.
func (m *Manager) GetWatchConfig() profile.Watch {
	cfg := m.Profile.Watch
	if cfg.DebounceDuration <= 0 {
		cfg.DebounceDuration = 2 * time.Second
	}
	if cfg.FullResyncEvery <= 0 {
		cfg.FullResyncEvery = 30 * time.Second
	}
	if len(cfg.IgnorePatterns) == 0 {
		cfg.IgnorePatterns = []string{"node_modules", "vendor", ".git", "dist", "build", "out"}
	}
	// Default to recursive watching for existing profiles that do not yet have
	// watch config persisted. This is the safer AI-context-manager default.
	if !cfg.Recursive && isZeroWatchConfig(m.Profile.Watch) {
		cfg.Recursive = true
	}
	return cfg
}

func isZeroWatchConfig(cfg profile.Watch) bool {
	return !cfg.Recursive &&
		cfg.DebounceDuration == 0 &&
		cfg.FullResyncEvery == 0 &&
		len(cfg.IgnorePatterns) == 0 &&
		len(cfg.MatchPatterns) == 0
}
