package discovery

import (
	"testing"
)

func TestDetectCI_GitHubActionsMultipleWorkflows(t *testing.T) {
	files := []string{".github/workflows/ci.yml", ".github/workflows/release.yml"}
	result := DetectCI("/repo", files)

	if len(result.Providers) != 1 || result.Providers[0] != "github-actions" {
		t.Errorf("expected [github-actions], got %v", result.Providers)
	}
	if len(result.ConfigFile) != 2 {
		t.Errorf("expected 2 config files, got %v", result.ConfigFile)
	}
}

func TestDetectCI_GitLabCI(t *testing.T) {
	files := []string{".gitlab-ci.yml"}
	result := DetectCI("/repo", files)

	if len(result.Providers) != 1 || result.Providers[0] != "gitlab-ci" {
		t.Errorf("expected [gitlab-ci], got %v", result.Providers)
	}
}

func TestDetectCI_Jenkins(t *testing.T) {
	files := []string{"Jenkinsfile"}
	result := DetectCI("/repo", files)

	if len(result.Providers) != 1 || result.Providers[0] != "jenkins" {
		t.Errorf("expected [jenkins], got %v", result.Providers)
	}
}

func TestDetectCI_Multiple(t *testing.T) {
	files := []string{".github/workflows/ci.yml", ".gitlab-ci.yml", "Jenkinsfile"}
	result := DetectCI("/repo", files)

	if len(result.Providers) != 3 {
		t.Errorf("expected 3 providers, got %v", result.Providers)
	}
}

func TestDetectCI_Empty(t *testing.T) {
	result := DetectCI("/repo", []string{})
	if len(result.Providers) != 0 {
		t.Errorf("expected 0 providers, got %v", result.Providers)
	}
}

func TestDetectCI_DuplicateNotAdded(t *testing.T) {
	files := []string{".github/workflows/ci.yml", ".github/workflows/ci.yml"}
	result := DetectCI("/repo", files)

	if len(result.Providers) != 1 {
		t.Errorf("expected 1 provider (no duplicates), got %v", result.Providers)
	}
}

func TestAppendIfMissing(t *testing.T) {
	items := []string{"a", "b"}
	result := appendIfMissing(items, "a")
	if len(result) != 2 {
		t.Errorf("expected no append for existing item, got %v", result)
	}

	result = appendIfMissing(items, "c")
	if len(result) != 3 || result[2] != "c" {
		t.Errorf("expected append for new item, got %v", result)
	}
}

func TestDetectConventions_GolangCI(t *testing.T) {
	files := []string{".golangci.yml"}
	result := DetectConventions("/repo", files)
	if result.Linter != "golangci-lint" {
		t.Errorf("expected golangci-lint, got %q", result.Linter)
	}
}

func TestDetectConventions_GolangCIYAML(t *testing.T) {
	files := []string{".golangci.yaml"}
	result := DetectConventions("/repo", files)
	if result.Linter != "golangci-lint" {
		t.Errorf("expected golangci-lint, got %q", result.Linter)
	}
}

func TestDetectConventions_Prettier(t *testing.T) {
	files := []string{".prettierrc"}
	result := DetectConventions("/repo", files)
	if result.Formatter != "prettier" {
		t.Errorf("expected prettier, got %q", result.Formatter)
	}
}

func TestDetectConventions_PrettierJSON(t *testing.T) {
	files := []string{".prettierrc.json"}
	result := DetectConventions("/repo", files)
	if result.Formatter != "prettier" {
		t.Errorf("expected prettier, got %q", result.Formatter)
	}
}

func TestDetectConventions_Biome(t *testing.T) {
	files := []string{"biome.json"}
	result := DetectConventions("/repo", files)
	if result.Linter != "biome" {
		t.Errorf("expected biome linter, got %q", result.Linter)
	}
	if result.Formatter != "biome" {
		t.Errorf("expected biome formatter, got %q", result.Formatter)
	}
}

func TestDetectConventions_Empty(t *testing.T) {
	result := DetectConventions("/repo", []string{})
	if result.Linter != "" || result.Formatter != "" {
		t.Errorf("expected empty conventions, got linter=%q formatter=%q", result.Linter, result.Formatter)
	}
}

func TestDetectConventions_Unrecognized(t *testing.T) {
	files := []string{"random.config", "somefile.txt"}
	result := DetectConventions("/repo", files)
	if result.Linter != "" || result.Formatter != "" {
		t.Errorf("expected empty conventions for unrecognized files, got linter=%q formatter=%q", result.Linter, result.Formatter)
	}
}
