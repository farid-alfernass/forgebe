package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// AtomicWrite writes data to a temporary file first, then renames it to the target path.
// This ensures that the write is atomic and doesn't leave a partial or corrupt file.
func AtomicWrite(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)

	// Create a temporary file in the same directory as the target file
	tmpFile, err := os.CreateTemp(dir, "forgebe-atomic-*.tmp")
	if err != nil {
		return fmt.Errorf("storage: failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Ensure temp file is cleaned up if anything fails
	defer func() {
		_ = tmpFile.Close()
		if _, err := os.Stat(tmpPath); err == nil {
			_ = os.Remove(tmpPath)
		}
	}()

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("storage: failed to write to temp file: %w", err)
	}

	// Ensure file is synced to disk
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("storage: failed to sync temp file: %w", err)
	}

	// Set permissions
	if err := tmpFile.Chmod(perm); err != nil {
		return fmt.Errorf("storage: failed to set temp file permissions: %w", err)
	}

	// Close temp file before renaming
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("storage: failed to close temp file: %w", err)
	}

	// Rename temp file to target path (atomic on most Unix-like filesystems)
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("storage: failed to rename temp file to target: %w", err)
	}

	return nil
}

// ReadFile reads the entire file from the given path.
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("storage: failed to read file: %w", err)
	}
	return data, nil
}

// Exists checks if a file or directory exists at the given path.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Remove removes the file or directory at the given path.
func Remove(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("storage: failed to remove path: %w", err)
	}
	return nil
}
