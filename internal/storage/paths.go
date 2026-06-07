package storage

import (
	"os"
	"path/filepath"
)

const (
	// DefaultRootDir is the default ForgeBE home directory
	DefaultRootDir = ".forgebe"

	// ProjectsDirName is the subdirectory for project profiles
	ProjectsDirName = "projects"

	// ExportsDirName is the subdirectory for exports
	ExportsDirName = "exports"

	// CacheDirName is the subdirectory for cache
	CacheDirName = "cache"

	// TmpDirName is the subdirectory for temporary files
	TmpDirName = "tmp"

	// LogsDirName is the subdirectory for logs
	LogsDirName = "logs"
)

// Paths manages ForgeBE storage paths
type Paths struct {
	root string
}

// NewPaths creates a new Paths instance
func NewPaths() (*Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	root := filepath.Join(home, DefaultRootDir)
	return &Paths{root: root}, nil
}

// NewPathsWithRoot creates a Paths instance with a custom root
func NewPathsWithRoot(root string) *Paths {
	return &Paths{root: root}
}

// Root returns the ForgeBE root directory
func (p *Paths) Root() string {
	return p.root
}

// EnsureBaseDirs creates the ForgeBE root directory and all standard subdirectories.
func (p *Paths) EnsureBaseDirs() error {
	dirs := []string{
		p.root,
		p.ProjectsDir(),
		p.ExportsDir(),
		p.CacheDir(),
		p.TmpDir(),
		p.LogsDir(),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0700); err != nil {
			return err
		}
	}
	return nil
}

// ProjectsDir returns the projects directory
func (p *Paths) ProjectsDir() string {
	return filepath.Join(p.root, ProjectsDirName)
}

// ProjectDir returns the directory for a specific project
func (p *Paths) ProjectDir(projectID string) string {
	return filepath.Join(p.ProjectsDir(), projectID)
}

// ProjectProfilePath returns the path to a project's profile file
func (p *Paths) ProjectProfilePath(projectID string) string {
	return filepath.Join(p.ProjectDir(projectID), "project-profile.yaml")
}

// ProjectDiscoveryPath returns the path to a project's discovery report
func (p *Paths) ProjectDiscoveryPath(projectID string) string {
	return filepath.Join(p.ProjectDir(projectID), "discovery-report.yaml")
}

// ProjectContractPath returns the path to a project's engineering contract
func (p *Paths) ProjectContractPath(projectID string) string {
	return filepath.Join(p.ProjectDir(projectID), "engineering-contract.md")
}

// ProjectSummaryPath returns the path to a project's summary
func (p *Paths) ProjectSummaryPath(projectID string) string {
	return filepath.Join(p.ProjectDir(projectID), "summary.md")
}

// ProjectMetadataPath returns the path to a project's metadata
func (p *Paths) ProjectMetadataPath(projectID string) string {
	return filepath.Join(p.ProjectDir(projectID), "metadata.json")
}

// ExportsDir returns the exports directory
func (p *Paths) ExportsDir() string {
	return filepath.Join(p.root, ExportsDirName)
}

// CacheDir returns the cache directory
func (p *Paths) CacheDir() string {
	return filepath.Join(p.root, CacheDirName)
}

// TmpDir returns the temporary directory
func (p *Paths) TmpDir() string {
	return filepath.Join(p.root, TmpDirName)
}

// LogsDir returns the logs directory
func (p *Paths) LogsDir() string {
	return filepath.Join(p.root, LogsDirName)
}

// EnsureDirectories creates all required directories
func (p *Paths) EnsureDirectories() error {
	dirs := []string{
		p.Root(),
		p.ProjectsDir(),
		p.ExportsDir(),
		p.CacheDir(),
		p.TmpDir(),
		p.LogsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// EnsureProjectDir creates the directory for a specific project
func (p *Paths) EnsureProjectDir(projectID string) error {
	return os.MkdirAll(p.ProjectDir(projectID), 0755)
}
