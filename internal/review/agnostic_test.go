package review

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoAIVendorReferences enforces ForgeBE's AI-agnostic guarantee: the
// review and git packages must never reference a specific AI vendor.
func TestNoAIVendorReferences(t *testing.T) {
	vendors := []string{"anthropic", "openai", "copilot", "cursor", "gemini", "hermes"}
	dirs := []string{".", "../git"}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read %s: %v", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
				continue
			}
			if strings.HasSuffix(e.Name(), "_test.go") {
				continue // this test names vendors on purpose
			}
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatalf("read %s: %v", e.Name(), err)
			}
			lower := strings.ToLower(string(data))
			for _, v := range vendors {
				if strings.Contains(lower, v) {
					t.Errorf("%s/%s references AI vendor %q — review/git must stay AI-agnostic", dir, e.Name(), v)
				}
			}
		}
	}
}
