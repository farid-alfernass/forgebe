package profile

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
	"gopkg.in/yaml.v3"
)

type Store struct {
	paths *storage.Paths
}

func NewStore(paths *storage.Paths) *Store { return &Store{paths: paths} }

func (s *Store) SaveProfile(profile ProjectProfile) error {
	if err := s.paths.EnsureProjectDir(profile.Metadata.ProfileID); err != nil {
		return fmt.Errorf("profile store: ensure project dir: %w", err)
	}
	b, err := yaml.Marshal(profile)
	if err != nil {
		return fmt.Errorf("profile store: marshal profile: %w", err)
	}
	if err := storage.AtomicWrite(s.paths.ProjectProfilePath(profile.Metadata.ProfileID), b, 0600); err != nil {
		return fmt.Errorf("profile store: write profile: %w", err)
	}
	return nil
}

func (s *Store) LoadProfile(projectID string) (*ProjectProfile, error) {
	b, err := os.ReadFile(s.paths.ProjectProfilePath(projectID))
	if err != nil {
		return nil, fmt.Errorf("profile store: read profile: %w", err)
	}
	var p ProjectProfile
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("profile store: unmarshal profile: %w", err)
	}
	return &p, nil
}

func (s *Store) SaveDiscovery(projectID string, report DiscoveryReport) error {
	if err := s.paths.EnsureProjectDir(projectID); err != nil {
		return fmt.Errorf("profile store: ensure project dir: %w", err)
	}
	b, err := yaml.Marshal(report)
	if err != nil {
		return fmt.Errorf("profile store: marshal discovery: %w", err)
	}
	if err := storage.AtomicWrite(s.paths.ProjectDiscoveryPath(projectID), b, 0600); err != nil {
		return fmt.Errorf("profile store: write discovery: %w", err)
	}
	return nil
}

func (s *Store) SaveMetadata(projectID string, meta any) error {
	if err := s.paths.EnsureProjectDir(projectID); err != nil {
		return fmt.Errorf("profile store: ensure project dir: %w", err)
	}
	b, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("profile store: marshal metadata: %w", err)
	}
	if err := storage.AtomicWrite(s.paths.ProjectMetadataPath(projectID), b, 0600); err != nil {
		return fmt.Errorf("profile store: write metadata: %w", err)
	}
	return nil
}
