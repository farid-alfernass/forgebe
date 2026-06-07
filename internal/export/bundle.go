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

// ExportBundle creates a portable .forgebe.zip bundle containing profile, contract, and metadata.
func ExportBundle(projectID string, paths *storage.Paths, outputPath string) error {
	projectDir := paths.ProjectDir(projectID)
	info, err := os.Stat(projectDir)
	if err != nil {
		return fmt.Errorf("export bundle: project dir: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("export bundle: %s is not a directory", projectDir)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("export bundle: create file: %w", err)
	}
	defer outputFile.Close()

	zw := zip.NewWriter(outputFile)
	defer zw.Close()

	err = filepath.WalkDir(projectDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == projectDir {
			return nil
		}

		rel, err := filepath.Rel(projectDir, path)
		if err != nil {
			return err
		}

		zipPath := filepath.Join(projectID, rel)

		if d.IsDir() {
			_, err := zw.Create(zipPath + "/")
			return err
		}

		if err := addFileToZip(zw, path, zipPath); err != nil {
			return fmt.Errorf("add %s to zip: %w", rel, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("export bundle: walk project dir: %w", err)
	}

	return nil
}

// ImportBundle extracts a .forgebe.zip bundle into targetDir.
func ImportBundle(bundlePath, targetDir string) error {
	reader, err := zip.OpenReader(bundlePath)
	if err != nil {
		return fmt.Errorf("import bundle: open zip: %w", err)
	}
	defer reader.Close()

	for _, f := range reader.File {
		fpath := filepath.Join(targetDir, f.Name)

		// ZipSlip protection
		cleanTarget := filepath.Clean(targetDir) + string(os.PathSeparator)
		if !strings.HasPrefix(filepath.Clean(fpath), cleanTarget) {
			return fmt.Errorf("import bundle: illegal file path %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0700)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0700); err != nil {
			return fmt.Errorf("import bundle: mkdir: %w", err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("import bundle: create file: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("import bundle: open zip entry: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return fmt.Errorf("import bundle: write file: %w", err)
		}
	}

	return nil
}

// DefaultBundlePath returns the default export path for a project bundle.
func DefaultBundlePath(projectID string, paths *storage.Paths) string {
	return filepath.Join(paths.ExportsDir(), fmt.Sprintf("%s.forgebe.zip", projectID))
}

func addFileToZip(zw *zip.Writer, srcPath, zipPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	w, err := zw.Create(filepath.ToSlash(zipPath))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
