package export

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

func TestExportBundle_ProjectDirNotFound(t *testing.T) {
	paths := storage.NewPathsWithRoot(t.TempDir())
	err := ExportBundle("nonexistent", paths, filepath.Join(t.TempDir(), "out.zip"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "project dir") {
		t.Errorf("expected 'project dir' error, got %q", err.Error())
	}
}

func TestExportBundle_ProjectDirIsFile(t *testing.T) {
	tmpDir := t.TempDir()
	paths := storage.NewPathsWithRoot(tmpDir)
	projDir := paths.ProjectDir("testproj")
	if err := os.MkdirAll(filepath.Dir(projDir), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(projDir, []byte("not a dir"), 0644); err != nil {
		t.Fatal(err)
	}
	err := ExportBundle("testproj", paths, filepath.Join(tmpDir, "out.zip"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "is not a directory") {
		t.Errorf("expected 'not a directory' error, got %q", err.Error())
	}
}

func TestExportBundle_Success(t *testing.T) {
	tmpDir := t.TempDir()
	paths := storage.NewPathsWithRoot(tmpDir)

	// Create a real project dir with files
	projDir := paths.ProjectDir("testproj")
	if err := os.MkdirAll(projDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, "CLAUDE.md"), []byte("# Claude"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, "project-profile.yaml"), []byte("version: '1.0'"), 0644); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(tmpDir, "out.zip")
	err := ExportBundle("testproj", paths, outPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify zip content
	r, err := zip.OpenReader(outPath)
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	defer r.Close()

	var files []string
	for _, f := range r.File {
		files = append(files, f.Name)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files in zip, got %d: %v", len(files), files)
	}
}

func TestDefaultBundlePath(t *testing.T) {
	paths := storage.NewPathsWithRoot("/tmp/test")
	path := DefaultBundlePath("myproj", paths)
	if !strings.HasSuffix(path, "myproj.forgebe.zip") {
		t.Errorf("expected path ending with myproj.forgebe.zip, got %s", path)
	}
}

func TestImportBundle_InvalidZip(t *testing.T) {
	tmpDir := t.TempDir()
	badZip := filepath.Join(tmpDir, "bad.zip")
	if err := os.WriteFile(badZip, []byte("not a zip"), 0644); err != nil {
		t.Fatal(err)
	}

	err := ImportBundle(badZip, tmpDir)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "open zip") {
		t.Errorf("expected 'open zip' error, got %q", err.Error())
	}
}

func TestImportBundle_SymlinkTraversalBlocked(t *testing.T) {
	tmpDir := t.TempDir()

	bundlePath := filepath.Join(tmpDir, "malicious.forgebe.zip")
	zf, err := os.Create(bundlePath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(zf)

	header := &zip.FileHeader{
		Name:   "../../evil.txt",
		Method: zip.Store,
	}
	w, err := zw.CreateHeader(header)
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte("evil"))
	zw.Close()
	zf.Close()

	err = ImportBundle(bundlePath, tmpDir)
	if err == nil {
		t.Fatal("expected error for path traversal, got nil")
	}
	if !strings.Contains(err.Error(), "illegal file path") {
		t.Errorf("expected 'illegal file path' error, got %q", err.Error())
	}
}

func TestImportBundle_Success(t *testing.T) {
	tmpDir := t.TempDir()

	bundlePath := filepath.Join(tmpDir, "testproj.forgebe.zip")
	zf, err := os.Create(bundlePath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(zf)

	fw, err := zw.Create("testproj/CLAUDE.md")
	if err != nil {
		t.Fatal(err)
	}
	fw.Write([]byte("# Claude config"))
	zw.Close()
	zf.Close()

	targetDir := filepath.Join(tmpDir, "imported")
	err = ImportBundle(bundlePath, targetDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(targetDir, "testproj/CLAUDE.md"))
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
	if string(content) != "# Claude config" {
		t.Errorf("file content mismatch, got: %s", string(content))
	}
}
