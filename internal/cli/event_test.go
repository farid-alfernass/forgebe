package cli

import (
	"path/filepath"
	"testing"
)

func TestIsRelevantEvent(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Manifest files (should be relevant)
		{"go.mod file", "go.mod", true},
		{"package.json file", "package.json", true},
		{"requirements.txt file", "requirements.txt", true},
		{"pom.xml file", "pom.xml", true},
		{"Cargo.toml file", "Cargo.toml", true},

		// Config files (should be relevant)
		{"Dockerfile file", "Dockerfile", true},
		{"docker-compose.yml file", "docker-compose.yml", true},
		{".gitignore file", ".gitignore", true},
		{".env file", ".env", true},
		{"tsconfig.json file", "tsconfig.json", true},
		{"webpack.config.js file", "webpack.config.js", true},
		{"jest.config.js file", "jest.config.js", true},

		// AI context files (not in manifest/config list, so NOT relevant)
		{"CLAUDE.md file", "CLAUDE.md", false},
		{".cursorrules file", ".cursorrules", false},

		// Source files in src/ or internal/ etc. (should be relevant)
		{"source file in src", filepath.Join("src", "main.go"), true},
		{"source file in internal", filepath.Join("internal", "handler.go"), true},
		{"source file in pkg", filepath.Join("pkg", "utils.go"), true},
		{"source file in cmd", filepath.Join("cmd", "app.go"), true},
		{"source file in services", filepath.Join("services", "api.go"), true},

		// Files to ignore (should NOT be relevant)
		{"hidden file (not special)", ".random", false},
		{"node_modules file", filepath.Join("node_modules", "dep", "index.js"), false},
		{"vendor file", filepath.Join("vendor", "dep.go"), false},
		{"dist file", filepath.Join("dist", "bundle.js"), false},
		{"build file", filepath.Join("build", "app.js"), false},
		{"out file", filepath.Join("out", "output.js"), false},
		{".git file", filepath.Join(".git", "config"), false},

		// Random files (should NOT be relevant unless they match patterns above)
		{"readme txt", "README.txt", false},
		{"image png", "image.png", false},
		{"regular go file", "main.go", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isRelevantEvent(tc.path)
			if result != tc.expected {
				t.Errorf("isRelevantEvent(%q) = %v, want %v", tc.path, result, tc.expected)
			}
		})
	}
}
