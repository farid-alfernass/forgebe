package discovery

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

var languageManifestMap = map[string]string{
	"go.mod":           "go",
	"package.json":     "javascript",
	"tsconfig.json":    "typescript",
	"pyproject.toml":   "python",
	"requirements.txt": "python",
	"pom.xml":          "java",
	"build.gradle":     "java",
	"Cargo.toml":       "rust",
}

var extensionLanguageMap = map[string]string{
	".go":   "go",
	".js":   "javascript",
	".ts":   "typescript",
	".tsx":  "typescript",
	".py":   "python",
	".java": "java",
	".kt":   "kotlin",
	".rs":   "rust",
}

func DetectLanguages(root string, files []string) []profile.LanguageDetection {
	manifestByLang := map[string]string{}
	counts := map[string]int{}
	exts := map[string]map[string]struct{}{}

	for _, f := range files {
		base := filepath.Base(f)
		if lang, ok := languageManifestMap[base]; ok {
			manifestByLang[lang] = f
		}
		ext := strings.ToLower(filepath.Ext(f))
		if lang, ok := extensionLanguageMap[ext]; ok {
			counts[lang]++
			if exts[lang] == nil {
				exts[lang] = map[string]struct{}{}
			}
			exts[lang][ext] = struct{}{}
		}
	}

	langSet := map[string]struct{}{}
	for lang := range manifestByLang {
		langSet[lang] = struct{}{}
	}
	for lang := range counts {
		langSet[lang] = struct{}{}
	}

	var results []profile.LanguageDetection
	for lang := range langSet {
		conf := "medium"
		if manifestByLang[lang] != "" {
			conf = "high"
		} else if counts[lang] <= 1 {
			conf = "low"
		}
		var extList []string
		for ext := range exts[lang] {
			extList = append(extList, ext)
		}
		sort.Strings(extList)
		results = append(results, profile.LanguageDetection{
			Language:   lang,
			Confidence: conf,
			Files:      counts[lang],
			Extensions: extList,
			Manifest:   manifestByLang[lang],
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Confidence == results[j].Confidence {
			if results[i].Files == results[j].Files {
				return results[i].Language < results[j].Language
			}
			return results[i].Files > results[j].Files
		}
		return confidenceRank(results[i].Confidence) > confidenceRank(results[j].Confidence)
	})
	return results
}

func confidenceRank(v string) int {
	switch v {
	case "high":
		return 3
	case "medium":
		return 2
	default:
		return 1
	}
}
