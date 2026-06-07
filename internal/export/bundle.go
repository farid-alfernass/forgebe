package export

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/faridtriwicaksono/forgebe/internal/storage"
)

// ExportBundle creates a portable .forgebe.zip bundle.
func ExportBundle(projectID string, paths *storage.Paths, outputPath string) (err error) {
	projectDir := paths.ProjectDir(projectID)
	info, err := os.Stat(projectDir)
	if err != nil {
		return fmt.Errorf("export bundle: project dir: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("export bundle: %s is not a directory", projectDir)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0700); err != nil {
		return fmt.Errorf("export bundle: create output dir: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("export bundle: create file: %w", err)
	}
	defer func() {
		if cerr := outputFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("export bundle: close file: %w", cerr)
		}
	}()

	zw := zip.NewWriter(outputFile)
	defer func() {
		if cerr := zw.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("export bundle: finalize zip: %w", cerr)
		}
	}()

	// Walk project dir and add files
	walkErr := filepath.WalkDir(projectDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == projectDir {
			return nil
		}
		rel, err := filepath.Rel(projectDir, path)
		if err != nil {
			return err
		}
		zipPath := filepath.Join(projectID, filepath.ToSlash(rel))

		if d.IsDir() {
			_, err := zw.Create(zipPath + "/")
			return err
		}
		return addFileToZip(zw, path, zipPath)
	})
	if walkErr != nil {
		return fmt.Errorf("export bundle: walk: %w", walkErr)
	}

	return nil
}

// ImportBundle extracts a .forgebe.zip bundle into targetDir using symlink-safe paths.
func ImportBundle(bundlePath, targetDir string) error {
	reader, err := zip.OpenReader(bundlePath)
	if err != nil {
		return fmt.Errorf("import bundle: open zip: %w", err)
	}
	defer reader.Close()

	// Resolve target to real path to prevent symlink escape
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("import bundle: resolve target: %w", err)
	}
	cleanTarget := filepath.Clean(absTarget)

	// Ensure target directory exists
	if err := os.MkdirAll(cleanTarget, 0700); err != nil {
		return fmt.Errorf("import bundle: create target: %w", err)
	}

	// Resolve real path after creation to handle intermediary symlinks
	realTarget, err := filepath.EvalSymlinks(cleanTarget)
	if err != nil {
		return fmt.Errorf("import bundle: eval symlinks target: %w", err)
	}

	for _, f := range reader.File {
		fpath := filepath.Join(realTarget, f.Name)

		if !strings.HasPrefix(filepath.Clean(fpath), realTarget+string(os.PathSeparator)) {
			return fmt.Errorf("import bundle: illegal file path %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0700); err != nil {
				return fmt.Errorf("import bundle: mkdir: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0700); err != nil {
			return fmt.Errorf("import bundle: mkdir parent: %w", err)
		}

		// Check if writing through a symlink (symlink traversal)
		parentDir, err := filepath.EvalSymlinks(filepath.Dir(fpath))
		if err != nil {
			return fmt.Errorf("import bundle: eval symlinks: %w", err)
		}
		if !strings.HasPrefix(parentDir, realTarget) {
			return fmt.Errorf("import bundle: symlink traversal blocked for %s", f.Name)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("import bundle: create %s: %w", f.Name, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("import bundle: open zip entry %s: %w", f.Name, err)
		}

		_, err = io.Copy(outFile, rc)

		if cerr := outFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("import bundle: close output file %s: %w", f.Name, cerr)
		}
		if cerr := rc.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("import bundle: close zip entry %s: %w", f.Name, cerr)
		}
		if err != nil {
			return fmt.Errorf("import bundle: write %s: %w", f.Name, err)
		}
	}
	return nil
}

// DefaultBundlePath returns the default export path for a project bundle.
func DefaultBundlePath(projectID string, paths *storage.Paths) string {
	return filepath.Join(paths.ExportsDir(), fmt.Sprintf("%s.forgebe.zip", projectID))
}

func addFileToZip(zw *zip.Writer, srcPath, zipPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := zw.Create(filepath.ToSlash(zipPath))
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err
}
