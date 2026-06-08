package discovery

import (
	"path/filepath"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func DetectDependencies(root string, files []string) profile.DependencyDetection {
	result := profile.DependencyDetection{}
	for _, f := range files {
		base := filepath.Base(f)
		switch base {
		case "go.mod":
			result.PackageManager = "go modules"
			result.ManifestFile = f
		case "pnpm-lock.yaml":
			result.LockFile = f
			if result.PackageManager == "" {
				result.PackageManager = "pnpm"
			}
		case "yarn.lock":
			result.LockFile = f
			if result.PackageManager == "" {
				result.PackageManager = "yarn"
			}
		case "package.json":
			if result.PackageManager == "" {
				result.PackageManager = "npm"
			}
			result.ManifestFile = f
		case "requirements.txt", "pyproject.toml":
			if result.ManifestFile == "" {
				result.ManifestFile = f
				result.PackageManager = "python"
			}
		case "pom.xml":
			result.PackageManager = "maven"
			result.ManifestFile = f
		case "build.gradle", "build.gradle.kts":
			if result.ManifestFile == "" {
				result.PackageManager = "gradle"
				result.ManifestFile = f
			}
		case "Cargo.toml":
			result.PackageManager = "cargo"
			result.ManifestFile = f
		}
	}
	if result.ManifestFile != "" {
		result.Dependencies = 1
	}
	return result
}
