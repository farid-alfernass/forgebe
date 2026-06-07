package onboarding

import (
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestApplyDiscoveryDefaults_NilSafe(t *testing.T) {
	// Should not panic on nil args
	ApplyDiscoveryDefaults(nil, nil)
	ApplyDiscoveryDefaults(&Answers{}, nil)
	ApplyDiscoveryDefaults(nil, &profile.DiscoveryReport{})
}

func TestApplyDiscoveryDefaults_FillsLanguage(t *testing.T) {
	a := &Answers{}
	report := &profile.DiscoveryReport{
		Detections: profile.Detections{
			Languages: []profile.LanguageDetection{
				{Language: "go", Confidence: "high"},
			},
			Dependencies: profile.DependencyDetection{PackageManager: "go modules"},
			Architecture: profile.ArchitectureDetection{Style: "layered-service-repository"},
			Testing:      profile.TestingDetection{Framework: "go test"},
		},
	}

	ApplyDiscoveryDefaults(a, report)

	if a.PrimaryLanguage != "go" {
		t.Errorf("PrimaryLanguage = %q, want %q", a.PrimaryLanguage, "go")
	}
	if a.PackageManager != "go modules" {
		t.Errorf("PackageManager = %q, want %q", a.PackageManager, "go modules")
	}
	if a.Architecture != "layered (service-repository)" {
		t.Errorf("Architecture = %q, want %q", a.Architecture, "layered (service-repository)")
	}
	if a.Framework != "go test" {
		t.Errorf("Framework = %q, want %q", a.Framework, "go test")
	}
}

func TestApplyDiscoveryDefaults_DoesNotOverwrite(t *testing.T) {
	a := &Answers{
		PrimaryLanguage: "python",
		PackageManager:  "pip",
		Architecture:    "mvc",
		Framework:       "pytest",
	}
	report := &profile.DiscoveryReport{
		Detections: profile.Detections{
			Languages: []profile.LanguageDetection{
				{Language: "go", Confidence: "high"},
			},
			Dependencies: profile.DependencyDetection{PackageManager: "go modules"},
			Architecture: profile.ArchitectureDetection{Style: "clean-hexagonal"},
			Testing:      profile.TestingDetection{Framework: "go test"},
		},
	}

	ApplyDiscoveryDefaults(a, report)

	if a.PrimaryLanguage != "python" {
		t.Errorf("should not overwrite PrimaryLanguage, got %q", a.PrimaryLanguage)
	}
	if a.PackageManager != "pip" {
		t.Errorf("should not overwrite PackageManager, got %q", a.PackageManager)
	}
	if a.Architecture != "mvc" {
		t.Errorf("should not overwrite Architecture, got %q", a.Architecture)
	}
	if a.Framework != "pytest" {
		t.Errorf("should not overwrite Framework, got %q", a.Framework)
	}
}

func TestNormalizeArchitecture(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"clean-hexagonal", "clean/hexagonal"},
		{"layered-service-repository", "layered (service-repository)"},
		{"mvc", "mvc"},
		{"unknown", "other / not sure"},
		{"", "other / not sure"},
	}

	for _, tc := range tests {
		got := normalizeArchitecture(tc.input)
		if got != tc.expected {
			t.Errorf("normalizeArchitecture(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestBuildProfile_NilAnswers(t *testing.T) {
	p := BuildProfile("test-id", "/tmp/test", nil, nil)
	if p.Metadata.ProfileID != "test-id" {
		t.Errorf("ProfileID = %q, want %q", p.Metadata.ProfileID, "test-id")
	}
	if p.Metadata.RepoPath != "/tmp/test" {
		t.Errorf("RepoPath = %q, want %q", p.Metadata.RepoPath, "/tmp/test")
	}
	if p.Version != "1" {
		t.Errorf("Version = %q, want %q", p.Version, "1")
	}
}

func TestBuildProfile_FullAnswers(t *testing.T) {
	answers := &Answers{
		ProjectName:      "myproject",
		ProjectMaturity:  "existing codebase",
		ProjectType:      "service",
		PrimaryLanguage:  "go",
		Framework:        "go test",
		PackageManager:   "go modules",
		Architecture:     "clean/hexagonal",
		TestingPolicy:    "TDD (tests required before implementation)",
		DependencyPolicy: "No, never add dependencies",
		ChangePolicy:     "Preserve structure, minimal changes",
		QualityPriority:  "Safety / minimal regression",
		DeliveryMode:     "Hybrid (plan for complex, direct for simple)",
		AIModelStrategy:  "Multi-model (different models for different tasks)",
		SensitiveAreas:   []string{"Authentication / Authorization", "Database migrations"},
	}

	report := &profile.DiscoveryReport{
		Detections: profile.Detections{
			Testing: profile.TestingDetection{Framework: "go test", TestDirs: []string{"internal/"}},
			CI:      profile.CIDetection{Providers: []string{"github-actions"}},
		},
	}

	p := BuildProfile("proj-123", "/tmp/myproj", answers, report)

	if p.Project.Name != "myproject" {
		t.Errorf("Name = %q, want %q", p.Project.Name, "myproject")
	}
	if p.Stack.PrimaryLanguage != "go" {
		t.Errorf("PrimaryLanguage = %q, want %q", p.Stack.PrimaryLanguage, "go")
	}
	if p.Policy.Testing.Strategy != "tdd" {
		t.Errorf("Testing.Strategy = %q, want %q", p.Policy.Testing.Strategy, "tdd")
	}
	if p.Policy.Testing.Required != true {
		t.Error("Testing.Required should be true")
	}
	if p.Policy.Dependencies.AllowAddition != false {
		t.Error("Dependencies.AllowAddition should be false")
	}
	if p.Policy.Dependencies.RequireApproval != true {
		t.Error("Dependencies.RequireApproval should be true")
	}
	if p.Policy.Changes.PreserveStructure != true {
		t.Error("Changes.PreserveStructure should be true")
	}
	if p.Policy.Delivery.Priority != "safety" {
		t.Errorf("Delivery.Priority = %q, want %q", p.Policy.Delivery.Priority, "safety")
	}
	if p.AI.InteractionMode != "hybrid" {
		t.Errorf("AI.InteractionMode = %q, want %q", p.AI.InteractionMode, "hybrid")
	}
	if p.AI.ModelStrategy != "multi" {
		t.Errorf("AI.ModelStrategy = %q, want %q", p.AI.ModelStrategy, "multi")
	}
	if len(p.Areas.SensitiveAreas) != 2 {
		t.Errorf("SensitiveAreas len = %d, want 2", len(p.Areas.SensitiveAreas))
	}
	if len(p.Stack.CI) != 1 || p.Stack.CI[0] != "github-actions" {
		t.Errorf("CI = %v, want [github-actions]", p.Stack.CI)
	}
}

func TestDefaultAnswersFromDiscovery(t *testing.T) {
	report := &profile.DiscoveryReport{
		Detections: profile.Detections{
			Languages: []profile.LanguageDetection{
				{Language: "python", Confidence: "high"},
			},
			Dependencies: profile.DependencyDetection{PackageManager: "pip"},
			Architecture: profile.ArchitectureDetection{Style: "mvc"},
			Testing:      profile.TestingDetection{Framework: "pytest"},
		},
	}

	a := DefaultAnswersFromDiscovery("myapp", report)

	if a.ProjectName != "myapp" {
		t.Errorf("ProjectName = %q, want %q", a.ProjectName, "myapp")
	}
	if a.PrimaryLanguage != "python" {
		t.Errorf("PrimaryLanguage = %q, want %q", a.PrimaryLanguage, "python")
	}
	if a.PackageManager != "pip" {
		t.Errorf("PackageManager = %q, want %q", a.PackageManager, "pip")
	}
	if a.Architecture != "mvc" {
		t.Errorf("Architecture = %q, want %q", a.Architecture, "mvc")
	}
	if a.Framework != "pytest" {
		t.Errorf("Framework = %q, want %q", a.Framework, "pytest")
	}
	if len(a.SensitiveAreas) != 3 {
		t.Errorf("SensitiveAreas len = %d, want 3", len(a.SensitiveAreas))
	}
}

func TestMapTestingPolicy(t *testing.T) {
	tests := []struct {
		input         string
		expectedStrat string
		expectedReqd  bool
	}{
		{"TDD (tests required before implementation)", "tdd", true},
		{"Tests only for critical paths", "risk-based", true},
		{"Test after implementation", "test-after", true},
		{"unknown", "mixed", false},
	}

	for _, tc := range tests {
		got := mapTestingPolicy(tc.input)
		if got.Strategy != tc.expectedStrat {
			t.Errorf("mapTestingPolicy(%q).Strategy = %q, want %q", tc.input, got.Strategy, tc.expectedStrat)
		}
		if got.Required != tc.expectedReqd {
			t.Errorf("mapTestingPolicy(%q).Required = %v, want %v", tc.input, got.Required, tc.expectedReqd)
		}
	}
}

func TestMapDependencyPolicy(t *testing.T) {
	tests := []struct {
		input       string
		allowAdd    bool
		reqApproval bool
	}{
		{"No, never add dependencies", false, true},
		{"Yes, freely", true, false},
		{"Yes, for small utility packages", true, true},
		{"unknown", false, true},
	}

	for _, tc := range tests {
		got := mapDependencyPolicy(tc.input)
		if got.AllowAddition != tc.allowAdd {
			t.Errorf("mapDependencyPolicy(%q).AllowAddition = %v, want %v", tc.input, got.AllowAddition, tc.allowAdd)
		}
		if got.RequireApproval != tc.reqApproval {
			t.Errorf("mapDependencyPolicy(%q).RequireApproval = %v, want %v", tc.input, got.RequireApproval, tc.reqApproval)
		}
	}
}

func TestMapChangePolicy(t *testing.T) {
	tests := []struct {
		input         string
		preserve      bool
		allowRefactor bool
		allowCross    bool
	}{
		{"Allow broad refactor with approval", false, true, true},
		{"Allow cross-file refactoring", false, true, true},
		{"Allow refactoring within same layer", true, true, false},
		{"Preserve structure, minimal changes", true, false, false},
	}

	for _, tc := range tests {
		got := mapChangePolicy(tc.input)
		if got.PreserveStructure != tc.preserve {
			t.Errorf("mapChangePolicy(%q).PreserveStructure = %v, want %v", tc.input, got.PreserveStructure, tc.preserve)
		}
		if got.AllowRefactor != tc.allowRefactor {
			t.Errorf("mapChangePolicy(%q).AllowRefactor = %v, want %v", tc.input, got.AllowRefactor, tc.allowRefactor)
		}
		if got.AllowCrossFile != tc.allowCross {
			t.Errorf("mapChangePolicy(%q).AllowCrossFile = %v, want %v", tc.input, got.AllowCrossFile, tc.allowCross)
		}
	}
}

func TestMapQualityPriority(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Delivery speed", "speed"},
		{"Performance", "performance"},
		{"Code quality / maintainability", "maintainability"},
		{"Safety / minimal regression", "safety"},
		{"unknown", "safety"},
	}

	for _, tc := range tests {
		got := mapQualityPriority(tc.input)
		if got != tc.expected {
			t.Errorf("mapQualityPriority(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestMapDeliveryMode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Direct execution (implement immediately if confident)", "direct"},
		{"Hybrid (plan for complex, direct for simple)", "hybrid"},
		{"unknown", "plan-first"},
	}

	for _, tc := range tests {
		got := mapDeliveryMode(tc.input)
		if got != tc.expected {
			t.Errorf("mapDeliveryMode(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestMapModelStrategy(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Single model (e.g., only Claude)", "single"},
		{"Round-robin (rotate models)", "round-robin"},
		{"Best available model", "best-available"},
		{"Multi-model (different models for different tasks)", "multi"},
		{"unknown", "multi"},
	}

	for _, tc := range tests {
		got := mapModelStrategy(tc.input)
		if got != tc.expected {
			t.Errorf("mapModelStrategy(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
