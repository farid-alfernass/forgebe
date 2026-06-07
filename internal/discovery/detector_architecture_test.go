package discovery

import "testing"

func TestDetectArchitecture_CleanHexagonal(t *testing.T) {
	dirs := []string{"domain", "application", "infrastructure", "interface"}
	result := DetectArchitecture(".", dirs)
	if result.Style != "clean-hexagonal" {
		t.Fatalf("expected clean-hexagonal, got %+v", result)
	}
}

func TestDetectArchitecture_Layered(t *testing.T) {
	dirs := []string{"controllers", "services", "repositories"}
	result := DetectArchitecture(".", dirs)
	if result.Style != "layered-service-repository" {
		t.Fatalf("expected layered-service-repository, got %+v", result)
	}
}

func TestDetectArchitecture_GoLayout(t *testing.T) {
	dirs := []string{"cmd", "internal"}
	result := DetectArchitecture(".", dirs)
	if result.Style != "go-service-layout" {
		t.Fatalf("expected go-service-layout, got %+v", result)
	}
}
