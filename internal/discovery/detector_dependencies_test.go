package discovery

import (
	"testing"
)

func TestDetectDependencies_Go(t *testing.T) {
	files := []string{"go.mod", "main.go"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "go modules" {
		t.Errorf("expected 'go modules', got %q", result.PackageManager)
	}
	if result.ManifestFile != "go.mod" {
		t.Errorf("expected 'go.mod', got %q", result.ManifestFile)
	}
	if result.Dependencies != 1 {
		t.Errorf("expected 1 dependency, got %d", result.Dependencies)
	}
}

func TestDetectDependencies_NPM(t *testing.T) {
	files := []string{"package.json", "package-lock.json"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "npm" {
		t.Errorf("expected 'npm', got %q", result.PackageManager)
	}
}

func TestDetectDependencies_PNPM(t *testing.T) {
	files := []string{"pnpm-lock.yaml", "package.json"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "pnpm" {
		t.Errorf("expected 'pnpm', got %q", result.PackageManager)
	}
	if result.LockFile == "" {
		t.Error("expected lock file to be set")
	}
}

func TestDetectDependencies_Yarn(t *testing.T) {
	files := []string{"yarn.lock", "package.json"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "yarn" {
		t.Errorf("expected 'yarn', got %q", result.PackageManager)
	}
}

func TestDetectDependencies_Python(t *testing.T) {
	files := []string{"requirements.txt", "main.py"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "python" {
		t.Errorf("expected 'python', got %q", result.PackageManager)
	}
}

func TestDetectDependencies_PythonPyproject(t *testing.T) {
	files := []string{"pyproject.toml", "main.py"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "python" {
		t.Errorf("expected 'python', got %q", result.PackageManager)
	}
}

func TestDetectDependencies_Maven(t *testing.T) {
	files := []string{"pom.xml", "src/Main.java"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "maven" {
		t.Errorf("expected 'maven', got %q", result.PackageManager)
	}
}

func TestDetectDependencies_Gradle(t *testing.T) {
	files := []string{"build.gradle", "src/main/Main.java"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "gradle" {
		t.Errorf("expected 'gradle', got %q", result.PackageManager)
	}
}

func TestDetectDependencies_GradleKTS(t *testing.T) {
	files := []string{"build.gradle.kts", "src/main/Main.kt"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "gradle" {
		t.Errorf("expected 'gradle', got %q", result.PackageManager)
	}
}

func TestDetectDependencies_Cargo(t *testing.T) {
	files := []string{"Cargo.toml", "src/main.rs"}
	result := DetectDependencies("/project", files)

	if result.PackageManager != "cargo" {
		t.Errorf("expected 'cargo', got %q", result.PackageManager)
	}
}

func TestDetectDependencies_Empty(t *testing.T) {
	result := DetectDependencies("/project", []string{})
	if result.PackageManager != "" {
		t.Errorf("expected empty package manager, got %q", result.PackageManager)
	}
}

func TestDetectDependencies_NoManifest(t *testing.T) {
	files := []string{"yarn.lock", "package-lock.json"}
	result := DetectDependencies("/project", files)

	// No manifest means no dependencies detected
	if result.Dependencies != 0 {
		t.Errorf("expected 0 dependencies, got %d", result.Dependencies)
	}
}
