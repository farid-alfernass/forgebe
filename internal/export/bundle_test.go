package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestBundleExportAndImport_Roundtrip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "forgebe-bundle-test-*")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	exportDir := filepath.Join(tmpDir, "exports")
	importDir := filepath.Join(tmpDir, "imports")
	os.MkdirAll(exportDir, 0700)
	os.MkdirAll(importDir, 0700)

	paths := storage.NewPathsWithRoot(filepath.Join(tmpDir, ".forgebe"))
	if err := paths.EnsureBaseDirs(); err != nil {
		t.Fatalf("ensure dirs: %v", err)
	}

	profileID := "test-project-123"
	projectDir := paths.ProjectDir(profileID)
	os.MkdirAll(projectDir, 0700)

	sample := profile.NewSampleProfile()
	sample.Metadata.ProfileID = profileID
	sample.Metadata.RepoPath = projectDir

	store := profile.NewStore(paths)
	if err := store.SaveProfile(sample); err != nil {
		t.Fatalf("save profile: %v", err)
	}

	contractContent := "# Engineering Contract\n\n## Testing\nStrategy: tdd\n"
	contractPath := paths.ProjectContractPath(profileID)
	if err := os.WriteFile(contractPath, []byte(contractContent), 0600); err != nil {
		t.Fatalf("write contract: %v", err)
	}

	bundlePath := filepath.Join(exportDir, "test-project-123.forgebe.zip")
	if err := ExportBundle(profileID, paths, bundlePath); err != nil {
		t.Fatalf("export bundle: %v", err)
	}

	if _, err := os.Stat(bundlePath); os.IsNotExist(err) {
		t.Fatal("bundle file was not created")
	}

	if err := ImportBundle(bundlePath, importDir); err != nil {
		t.Fatalf("import bundle: %v", err)
	}

	importedProfilePath := filepath.Join(importDir, profileID, "project-profile.yaml")
	if _, err := os.Stat(importedProfilePath); os.IsNotExist(err) {
		t.Fatal("imported profile not found")
	}

	importedContractPath := filepath.Join(importDir, profileID, "engineering-contract.md")
	data, err := os.ReadFile(importedContractPath)
	if err != nil {
		t.Fatalf("read imported contract: %v", err)
	}
	if !strings.Contains(string(data), "Strategy: tdd") {
		t.Fatal("imported contract content mismatch")
	}

	info, err := os.Stat(bundlePath)
	if err != nil {
		t.Fatalf("stat bundle: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("bundle file is empty")
	}
}
