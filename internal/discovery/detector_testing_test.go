package discovery

import "testing"

func TestDetectTesting_Go(t *testing.T) {
	files := []string{"internal/user/user_test.go", "cmd/app/main.go"}
	result := DetectTesting(".", files, nil)
	if result.Framework != "go test" {
		t.Fatalf("expected go test, got %+v", result)
	}
}

func TestDetectTesting_Pytest(t *testing.T) {
	files := []string{"pytest.ini", "tests/test_app.py"}
	result := DetectTesting(".", files, nil)
	if result.Framework != "pytest" {
		t.Fatalf("expected pytest, got %+v", result)
	}
}

func TestDetectDependencies_GoModules(t *testing.T) {
	files := []string{"go.mod", "go.sum"}
	result := DetectDependencies(".", files)
	if result.PackageManager != "go modules" {
		t.Fatalf("expected go modules, got %+v", result)
	}
}

func TestDetectCI_GitHubActions(t *testing.T) {
	files := []string{".github/workflows/ci.yml"}
	result := DetectCI(".", files)
	if len(result.Providers) == 0 || result.Providers[0] != "github-actions" {
		t.Fatalf("expected github-actions, got %+v", result)
	}
}
