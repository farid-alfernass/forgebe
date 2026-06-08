package profile

import (
	"os"
	"testing"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestStoreSaveAndLoadProfile(t *testing.T) {
	paths := storage.NewPathsWithRoot(t.TempDir())
	store := NewStore(paths)
	profile := ProjectProfile{
		Version:  "1",
		Metadata: Metadata{ProfileID: "demo_123", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Project:  Project{Name: "demo", Maturity: "existing", Type: "service", Languages: []string{"go"}},
		Stack:    Stack{PrimaryLanguage: "go", Architecture: "layered"},
	}
	if err := store.SaveProfile(profile); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := store.LoadProfile("demo_123")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.Project.Name != "demo" {
		t.Fatalf("unexpected profile: %+v", loaded)
	}
}

func TestStoreSaveDiscovery(t *testing.T) {
	paths := storage.NewPathsWithRoot(t.TempDir())
	store := NewStore(paths)

	report := DiscoveryReport{
		Timestamp: time.Now(),
		RepoPath:  "/test",
		Detections: Detections{
			Languages: []LanguageDetection{{Language: "go", Confidence: "high"}},
		},
	}
	if err := store.SaveDiscovery("proj_abc", report); err != nil {
		t.Fatalf("save discovery failed: %v", err)
	}

	// Verify file exists
	path := paths.ProjectDiscoveryPath("proj_abc")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("discovery file not created at %s", path)
	}
}

func TestStoreSaveMetadata(t *testing.T) {
	paths := storage.NewPathsWithRoot(t.TempDir())
	store := NewStore(paths)

	meta := map[string]string{"last_sync": "2026-06-08T00:00:00Z"}
	if err := store.SaveMetadata("proj_xyz", meta); err != nil {
		t.Fatalf("save metadata failed: %v", err)
	}

	// Verify file exists
	path := paths.ProjectMetadataPath("proj_xyz")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("metadata file not created at %s", path)
	}
}

func TestStoreLoadProfile_NotFound(t *testing.T) {
	paths := storage.NewPathsWithRoot(t.TempDir())
	store := NewStore(paths)

	_, err := store.LoadProfile("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
}
