package discovery

import "testing"

func TestDetectLanguages_GoProject(t *testing.T) {
	files := []string{"go.mod", "cmd/app/main.go", "internal/service/user.go", "internal/service/user_test.go"}
	langs := DetectLanguages(".", files)
	if len(langs) == 0 || langs[0].Language != "go" {
		t.Fatalf("expected go as primary language, got %+v", langs)
	}
	if langs[0].Confidence != "high" {
		t.Fatalf("expected high confidence, got %+v", langs[0])
	}
}

func TestDetectLanguages_NodeProject(t *testing.T) {
	files := []string{"package.json", "src/index.ts", "src/app.ts", "vitest.config.ts"}
	langs := DetectLanguages(".", files)
	if len(langs) == 0 {
		t.Fatal("expected language detection result")
	}
}
