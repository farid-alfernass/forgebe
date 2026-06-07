package discovery

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func expectedConfidence(t *testing.T, report *profile.DiscoveryReport) {
	t.Helper()
	if report.Confidence["language"] != "high" {
		t.Fatalf("expected high language confidence, got %s", report.Confidence["language"])
	}
}

func expectedGoLanguage(t *testing.T, report *profile.DiscoveryReport) {
	t.Helper()
	foundGo := false
	for _, l := range report.Detections.Languages {
		if l.Language == "go" {
			foundGo = true
			break
		}
	}
	if !foundGo {
		t.Fatal("expected go to be detected")
	}
}

func expectedLayeredArchitecture(t *testing.T, report *profile.DiscoveryReport) {
	t.Helper()
	if report.Detections.Architecture.Style != "layered-service-repository" {
		t.Fatalf("expected layered architecture, got %s", report.Detections.Architecture.Style)
	}
}

func TestScannerGoLayered(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "../../testdata/repos/go-layered")
	scanner := NewScanner()
	report, err := scanner.Scan(root)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	expectedLayeredArchitecture(t, report)
	expectedGoLanguage(t, report)
	expectedConfidence(t, report)
}
