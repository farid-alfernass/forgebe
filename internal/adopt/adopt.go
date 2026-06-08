package adopt

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/discovery"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

// Result represents the outcome of an adopt operation.
type Result struct {
	ProjectID   string
	ProfilePath string
	RepoPath    string
	Language    string
	Framework   string
	Actions     []string
	DryRun      bool
}

// Options configures the adopt behavior.
type Options struct {
	RepoPath string
	DryRun   bool
	Force    bool
}

// Run performs the adopt operation on a given repository path.
func Run(opts Options) (*Result, error) {
	if opts.RepoPath == "" {
		return nil, fmt.Errorf("adopt: repo path is empty")
	}

	// Verify path exists and is a directory
	info, err := os.Stat(opts.RepoPath)
	if err != nil {
		return nil, fmt.Errorf("adopt: cannot access repo path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("adopt: %s is not a directory", opts.RepoPath)
	}

	absPath, err := filepath.Abs(opts.RepoPath)
	if err != nil {
		return nil, fmt.Errorf("adopt: resolve path: %w", err)
	}

	// Generate project ID
	projectID := profile.ProjectID(absPath, "")

	// Check if profile already exists
	paths, err := storage.NewPaths()
	if err != nil {
		return nil, fmt.Errorf("adopt: storage paths: %w", err)
	}

	store := profile.NewStore(paths)
	existing, _ := store.LoadProfile(projectID)
	if existing != nil && !opts.Force && !opts.DryRun {
		return nil, fmt.Errorf("adopt: project %q already has a profile (use --force to overwrite)", projectID)
	}

	// Run discovery scan
	scanner := discovery.NewScanner()
	report, err := scanner.Scan(absPath)
	if err != nil {
		return nil, fmt.Errorf("adopt: scan failed: %w", err)
	}

	// Build profile from discovery report
	proj := buildProfileFromReport(absPath, projectID, report)

	actions := []string{
		fmt.Sprintf("Create profile: %s", projectID),
		fmt.Sprintf("Language detected: %s", proj.Stack.PrimaryLanguage),
	}
	if proj.Stack.Architecture != "" {
		actions = append(actions, fmt.Sprintf("Architecture: %s", proj.Stack.Architecture))
	}
	if proj.Stack.TestFramework != "" {
		actions = append(actions, fmt.Sprintf("Test framework: %s", proj.Stack.TestFramework))
	}

	result := &Result{
		ProjectID:   projectID,
		ProfilePath: paths.ProjectProfilePath(projectID),
		RepoPath:    absPath,
		Language:    proj.Stack.PrimaryLanguage,
		Framework:   proj.Stack.Framework,
		Actions:     actions,
		DryRun:      opts.DryRun,
	}

	if opts.DryRun {
		return result, nil
	}

	// Save profile
	if err := paths.EnsureProjectDir(projectID); err != nil {
		return nil, fmt.Errorf("adopt: ensure project dir: %w", err)
	}

	if err := store.SaveProfile(proj); err != nil {
		return nil, fmt.Errorf("adopt: save profile: %w", err)
	}

	// Save discovery report
	if err := store.SaveDiscovery(projectID, *report); err != nil {
		return nil, fmt.Errorf("adopt: save discovery: %w", err)
	}

	return result, nil
}

func buildProfileFromReport(repoPath, projectID string, report *profile.DiscoveryReport) profile.ProjectProfile {
	projectName := filepath.Base(repoPath)

	// Extract primary language
	primaryLang := ""
	languages := make([]string, 0)
	if len(report.Detections.Languages) > 0 {
		primaryLang = report.Detections.Languages[0].Language
		for _, l := range report.Detections.Languages {
			languages = append(languages, l.Language)
		}
	}

	p := profile.ProjectProfile{
		Version: "1.0",
		Metadata: profile.Metadata{
			ProfileID:  projectID,
			RepoPath:   repoPath,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			ForgebeVer: "0.6.0",
			InitMethod: "adopt",
		},
		Project: profile.Project{
			Name:      projectName,
			Maturity:  "existing",
			Type:      "service",
			Languages: languages,
		},
		Stack: profile.Stack{
			PrimaryLanguage: primaryLang,
			Architecture:    report.Detections.Architecture.Style,
			TestFramework:   report.Detections.Testing.Framework,
			PackageManager:  report.Detections.Dependencies.PackageManager,
		},
		Policy: profile.Policy{
			Testing: profile.TestingPolicy{
				Required:     true,
				Coverage:     80,
				Strategy:     "tdd",
				UnitRequired: true,
			},
			Dependencies: profile.DependencyPolicy{
				AllowAddition:   true,
				RequireApproval: true,
			},
			Changes: profile.ChangePolicy{
				PreserveStructure: true,
				AllowRefactor:     true,
				AllowCrossFile:    true,
			},
			Quality: profile.QualityPolicy{
				LintRequired:   true,
				FormatRequired: true,
				ReviewRequired: true,
			},
			Delivery: profile.DeliveryPolicy{
				Mode:     "hybrid",
				Priority: "safety",
			},
		},
		AI: profile.AI{
			InteractionMode: "hybrid",
			ModelStrategy:   "multi",
		},
	}

	return p
}
