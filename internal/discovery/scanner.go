package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// Scanner scans a repository filesystem and produces a DiscoveryReport.
// It walks the directory tree up to a configurable depth, skipping common
// vendored, generated, and build artifact directories (e.g. .git, node_modules, vendor).
// Each detection category (language, architecture, testing, CI, conventions) is
// delegated to its own detector function for modularity and testability.
type Scanner struct {
	maxDepth int
}

// NewScanner creates a default repository scanner with a sane max directory depth.
func NewScanner() *Scanner {
	return &Scanner{maxDepth: 6}
}

// Scan walks the given repoPath, runs all detectors, and returns a structured
// DiscoveryReport with confidence levels and engineering assumptions.
// Returns an error if path is inaccessible or not a directory.
func (s *Scanner) Scan(repoPath string) (*profile.DiscoveryReport, error) {
	root, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("discovery: resolve repo path: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("discovery: stat repo path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("discovery: repo path is not a directory: %s", root)
	}

	files, dirs, err := walkProject(root, s.maxDepth)
	if err != nil {
		return nil, err
	}

	language := DetectLanguages(root, files)
	architecture := DetectArchitecture(root, dirs)
	testing := DetectTesting(root, files, dirs)
	dependencies := DetectDependencies(root, files)
	ci := DetectCI(root, files)
	conventions := DetectConventions(root, files)

	confidence := map[string]string{
		"language":     aggregateLanguageConfidence(language),
		"architecture": architecture.Confidence,
		"testing":      confidenceFromValue(testing.Framework),
		"dependencies": confidenceFromValue(dependencies.ManifestFile),
		"ci":           confidenceFromSlice(ci.Providers),
	}

	assumptions := buildAssumptions(language, architecture, testing, dependencies)
	suggestions := buildSuggestions(language, architecture, testing, dependencies)

	return &profile.DiscoveryReport{
		Timestamp: time.Now(),
		RepoPath:  root,
		Detections: profile.Detections{
			Languages:    language,
			Architecture: architecture,
			Testing:      testing,
			Dependencies: dependencies,
			CI:           ci,
			Conventions:  conventions,
		},
		Confidence:  confidence,
		Assumptions: assumptions,
		Suggestions: suggestions,
	}, nil
}

// walkProject recursively walks a root directory, collecting file and directory paths
// while skipping common vendored and build artifact directories.
func walkProject(root string, maxDepth int) ([]string, []string, error) {
	var files []string
	var dirs []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("discovery: walk %s: %w", path, err)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		if shouldSkipPath(rel, d.IsDir()) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if depth(rel) > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			dirs = append(dirs, rel)
		} else {
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return files, dirs, nil
}

// shouldSkipPath returns true if the path contains a directory that should be ignored.
func shouldSkipPath(rel string, isDir bool) bool {
	parts := splitPath(rel)
	if len(parts) == 0 {
		return false
	}
	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, "dist": true,
		"build": true, "target": true, ".next": true, ".nuxt": true,
		".venv": true, "venv": true, "__pycache__": true, ".idea": true, ".vscode": true,
		"testdata": true, "fixtures": true, "examples": true,
	}
	for _, part := range parts {
		if skipDirs[part] {
			return true
		}
	}
	return false
}

// splitPath splits a file path into its components.
func splitPath(path string) []string {
	path = filepath.ToSlash(path)
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// depth returns the number of directory levels in a relative path.
func depth(rel string) int {
	if rel == "" || rel == "." {
		return 0
	}
	count := 1
	for _, r := range filepath.ToSlash(rel) {
		if r == '/' {
			count++
		}
	}
	return count
}

func confidenceFromValue(v string) string {
	if v == "" {
		return "low"
	}
	return "high"
}

func confidenceFromSlice(v []string) string {
	if len(v) == 0 {
		return "low"
	}
	return "high"
}

func aggregateLanguageConfidence(langs []profile.LanguageDetection) string {
	if len(langs) == 0 {
		return "low"
	}
	return langs[0].Confidence
}

func buildAssumptions(languages []profile.LanguageDetection, arch profile.ArchitectureDetection, testing profile.TestingDetection, deps profile.DependencyDetection) []string {
	assumptions := []string{
		"ForgeBE will preserve the existing repository structure unless explicitly configured otherwise.",
		"ForgeBE will not add dependencies without explicit approval.",
		"ForgeBE stores project context locally and does not write tracked files into the repository by default.",
	}
	if len(languages) > 0 {
		assumptions = append(assumptions, "Primary language is inferred from manifest files and file counts.")
	}
	if arch.Style != "" {
		assumptions = append(assumptions, "Architecture is inferred from folder naming patterns and should be confirmed during guided onboarding.")
	}
	if testing.Framework != "" {
		assumptions = append(assumptions, "New business logic should include tests using the detected test framework.")
	}
	return assumptions
}

func buildSuggestions(languages []profile.LanguageDetection, arch profile.ArchitectureDetection, testing profile.TestingDetection, deps profile.DependencyDetection) []string {
	var suggestions []string
	if len(languages) == 0 {
		suggestions = append(suggestions, "Run guided onboarding to define the primary language manually.")
	}
	if arch.Confidence == "low" {
		suggestions = append(suggestions, "Confirm architecture style manually because folder structure is not conclusive.")
	}
	if testing.Framework == "" {
		suggestions = append(suggestions, "Define a test policy manually; no clear test framework was detected.")
	}
	return suggestions
}
