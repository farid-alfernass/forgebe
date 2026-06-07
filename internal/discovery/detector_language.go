package discovery

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// manifestLangMap maps filenames to languages.
var manifestLangMap = map[string]string{
	"go.mod":           "go",
	"package.json":     "javascript",
	"pom.xml":          "java",
	"requirements.txt": "python",
	"pyproject.toml":   "python",
	"Cargo.toml":       "rust",
	"Gemfile":          "ruby",
	"composer.json":    "php",
}

var extLangMap = map[string]string{
	".go":    "go",
	".js":    "javascript",
	".ts":    "typescript",
	".tsx":   "typescript",
	".py":    "python",
	".java":  "java",
	".rs":    "rust",
	".rb":    "ruby",
	".php":   "php",
	".cpp":   "cpp",
	".cs":    "csharp",
	".swift": "swift",
	".kt":    "kotlin",
}

// DetectLanguages detects primary languages with root-priority.
func DetectLanguages(root string, files []string) []profile.LanguageDetection {
	type langInfo struct {
		language   string
		confidence float64 // 0.0-1.0 internal, mapped to string later
		manifest   string
	}
	langs := make(map[string]*langInfo)
	extCount := make(map[string]int)

	for _, file := range files {
		base := filepath.Base(file)
		rel, _ := filepath.Rel(root, file)
		depth := strings.Count(rel, string(filepath.Separator))

		if lang, ok := manifestLangMap[base]; ok {
			if _, exists := langs[lang]; !exists {
				langs[lang] = &langInfo{language: lang, manifest: base}
			}
			// Root (depth 0) gets highest confidence
			candidate := 1.0
			if depth > 0 {
				candidate = 0.7 - (float64(depth) * 0.15)
				if candidate < 0.1 {
					candidate = 0.1
				}
			}
			if candidate > langs[lang].confidence {
				langs[lang].confidence = candidate
			}
		}

		ext := filepath.Ext(file)
		if ext != "" {
			extCount[ext]++
		}
	}

	// Count extensions to compute secondary evidence
	for ext := range extCount {
		if lang, ok := extLangMap[ext]; ok {
			if _, exists := langs[lang]; !exists {
				langs[lang] = &langInfo{language: lang, confidence: 0.3}
			}
		}
	}

	// Sort by confidence descending
	sorted := make([]*langInfo, 0, len(langs))
	for _, info := range langs {
		sorted = append(sorted, info)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].confidence > sorted[j].confidence
	})

	var result []profile.LanguageDetection
	for _, info := range sorted {
		result = append(result, profile.LanguageDetection{
			Language:   info.language,
			Confidence: confidenceLabel(info.confidence),
			Manifest:   info.manifest,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func confidenceLabel(score float64) string {
	if score >= 0.9 {
		return "high"
	} else if score >= 0.5 {
		return "medium"
	}
	return "low"
}
