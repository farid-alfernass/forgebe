package discovery

import (
	"path/filepath"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func DetectCI(root string, files []string) profile.CIDetection {
	result := profile.CIDetection{}
	for _, f := range files {
		base := filepath.Base(f)
		switch {
		case f == ".github/workflows/ci.yml" || filepath.Dir(f) == ".github/workflows":
			result.Providers = appendIfMissing(result.Providers, "github-actions")
			result.ConfigFile = appendIfMissing(result.ConfigFile, f)
		case base == ".gitlab-ci.yml":
			result.Providers = appendIfMissing(result.Providers, "gitlab-ci")
			result.ConfigFile = appendIfMissing(result.ConfigFile, f)
		case base == "Jenkinsfile":
			result.Providers = appendIfMissing(result.Providers, "jenkins")
			result.ConfigFile = appendIfMissing(result.ConfigFile, f)
		}
	}
	return result
}

func appendIfMissing(items []string, v string) []string {
	for _, item := range items {
		if item == v {
			return items
		}
	}
	return append(items, v)
}

func DetectConventions(root string, files []string) profile.ConventionDetection {
	result := profile.ConventionDetection{}
	for _, f := range files {
		base := filepath.Base(f)
		switch base {
		case ".golangci.yml", ".golangci.yaml":
			result.Linter = "golangci-lint"
		case ".prettierrc", ".prettierrc.json":
			result.Formatter = "prettier"
		case "biome.json":
			result.Linter = "biome"
			result.Formatter = "biome"
		}
	}
	return result
}
