package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWriteWritesFileWithPermission(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile.yaml")
	if err := AtomicWrite(path, []byte("hello"), 0600); err != nil {
		t.Fatalf("AtomicWrite failed: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected content: %q", string(data))
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
