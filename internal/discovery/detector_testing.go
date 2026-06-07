package discovery

import (
	"path/filepath"
	"strings"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func DetectTesting(root string, files, dirs []string) profile.TestingDetection {
	result := profile.TestingDetection{}
	seenDirs := map[string]struct{}{}

	for _, f := range files {
		base := filepath.Base(f)
		switch {
		case strings.HasSuffix(base, "_test.go"):
			result.Framework = "go test"
			result.Coverage = true
		case base == "jest.config.js" || base == "jest.config.ts":
			result.Framework = "jest"
			result.ConfigFile = f
		case base == "pytest.ini" || base == "conftest.py":
			result.Framework = "pytest"
			result.ConfigFile = f
		case base == "vitest.config.ts" || base == "vitest.config.js":
			result.Framework = "vitest"
			result.ConfigFile = f
		}

		dir := filepath.Dir(f)
		if dir != "." && (strings.Contains(strings.ToLower(dir), "test") || strings.Contains(strings.ToLower(dir), "spec")) {
			seenDirs[dir] = struct{}{}
		}
	}

	for dir := range seenDirs {
		result.TestDirs = append(result.TestDirs, dir)
	}
	return result
}
