package profile

import (
	"testing"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestStoreSaveAndLoadProfile(t *testing.T) {
	paths := storage.NewPathsWithRoot(t.TempDir())
	store := NewStore(paths)
	profile := ProjectProfile{
		Version: "1",
		Metadata: Metadata{ProfileID: "demo_123", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Project: Project{Name: "demo", Maturity: "existing", Type: "service", Languages: []string{"go"}},
		Stack: Stack{PrimaryLanguage: "go", Architecture: "layered"},
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
