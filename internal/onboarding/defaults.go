package onboarding

import (
	"strings"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// ApplyDiscoveryDefaults pre-fills onboarding answers from discovery results.
func ApplyDiscoveryDefaults(a *Answers, report *profile.DiscoveryReport) {
	if a == nil || report == nil {
		return
	}
	if a.PrimaryLanguage == "" && len(report.Detections.Languages) > 0 {
		a.PrimaryLanguage = report.Detections.Languages[0].Language
	}
	if a.PackageManager == "" {
		a.PackageManager = report.Detections.Dependencies.PackageManager
	}
	if a.Architecture == "" {
		a.Architecture = normalizeArchitecture(report.Detections.Architecture.Style)
	}
	if a.Framework == "" {
		a.Framework = report.Detections.Testing.Framework
	}
}

func normalizeArchitecture(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "clean-hexagonal":
		return "clean/hexagonal"
	case "layered-service-repository":
		return "layered (service-repository)"
	case "mvc":
		return "mvc"
	default:
		return "other / not sure"
	}
}

// BuildProfile builds canonical ProjectProfile from onboarding answers.
func BuildProfile(projectID, repoPath string, answers *Answers, report *profile.DiscoveryReport) profile.ProjectProfile {
	now := time.Now()
	p := profile.ProjectProfile{
		Version: "1",
		Metadata: profile.Metadata{
			ProfileID:  projectID,
			RepoPath:   repoPath,
			CreatedAt:  now,
			UpdatedAt:  now,
			ForgebeVer: "dev",
			InitMethod: "guided",
		},
	}
	if answers == nil {
		return p
	}

	p.Project.Name = answers.ProjectName
	p.Project.Maturity = answers.ProjectMaturity
	p.Project.Type = answers.ProjectType
	if answers.PrimaryLanguage != "" {
		p.Project.Languages = []string{answers.PrimaryLanguage}
		p.Stack.PrimaryLanguage = answers.PrimaryLanguage
	}
	p.Stack.Framework = answers.Framework
	p.Stack.PackageManager = answers.PackageManager
	p.Stack.Architecture = answers.Architecture
	if report != nil {
		p.Stack.TestFramework = report.Detections.Testing.Framework
		p.Stack.Linter = report.Detections.Conventions.Linter
		p.Stack.Formatter = report.Detections.Conventions.Formatter
		p.Stack.CI = report.Detections.CI.Providers
		p.Areas.TestRoots = report.Detections.Testing.TestDirs
	}

	p.Policy.Testing = mapTestingPolicy(answers.TestingPolicy)
	p.Policy.Dependencies = mapDependencyPolicy(answers.DependencyPolicy)
	p.Policy.Changes = mapChangePolicy(answers.ChangePolicy)
	p.Policy.Quality = profile.QualityPolicy{LintRequired: true, FormatRequired: true, ReviewRequired: true}
	p.Policy.Delivery = profile.DeliveryPolicy{Mode: mapDeliveryMode(answers.DeliveryMode), Priority: mapQualityPriority(answers.QualityPriority)}
	p.AI = profile.AI{InteractionMode: mapDeliveryMode(answers.DeliveryMode), ModelStrategy: mapModelStrategy(answers.AIModelStrategy)}

	for _, area := range answers.SensitiveAreas {
		if area == "Other (specify)" {
			continue
		}
		p.Areas.SensitiveAreas = append(p.Areas.SensitiveAreas, area)
	}
	if strings.TrimSpace(answers.OtherSensitive) != "" {
		p.Areas.SensitiveAreas = append(p.Areas.SensitiveAreas, strings.TrimSpace(answers.OtherSensitive))
	}
	return p
}

func mapTestingPolicy(v string) profile.TestingPolicy {
	switch v {
	case "TDD (tests required before implementation)":
		return profile.TestingPolicy{Required: true, Strategy: "tdd", UnitRequired: true, IntegrationReqd: false}
	case "Tests only for critical paths":
		return profile.TestingPolicy{Required: true, Strategy: "risk-based", UnitRequired: true, IntegrationReqd: true}
	case "Test after implementation":
		return profile.TestingPolicy{Required: true, Strategy: "test-after", UnitRequired: true, IntegrationReqd: false}
	default:
		return profile.TestingPolicy{Required: false, Strategy: "mixed", UnitRequired: false, IntegrationReqd: false}
	}
}

func mapDependencyPolicy(v string) profile.DependencyPolicy {
	switch v {
	case "No, never add dependencies":
		return profile.DependencyPolicy{AllowAddition: false, RequireApproval: true}
	case "Yes, freely":
		return profile.DependencyPolicy{AllowAddition: true, RequireApproval: false}
	case "Yes, for small utility packages":
		return profile.DependencyPolicy{AllowAddition: true, RequireApproval: true}
	default:
		return profile.DependencyPolicy{AllowAddition: false, RequireApproval: true}
	}
}

func mapChangePolicy(v string) profile.ChangePolicy {
	switch v {
	case "Allow broad refactor with approval":
		return profile.ChangePolicy{PreserveStructure: false, AllowRefactor: true, AllowCrossFile: true, RequireApproval: []string{"broad-refactor"}}
	case "Allow cross-file refactoring":
		return profile.ChangePolicy{PreserveStructure: false, AllowRefactor: true, AllowCrossFile: true}
	case "Allow refactoring within same layer":
		return profile.ChangePolicy{PreserveStructure: true, AllowRefactor: true, AllowCrossFile: false}
	default:
		return profile.ChangePolicy{PreserveStructure: true, AllowRefactor: false, AllowCrossFile: false}
	}
}

func mapQualityPriority(v string) string {
	switch v {
	case "Delivery speed":
		return "speed"
	case "Performance":
		return "performance"
	case "Code quality / maintainability":
		return "maintainability"
	default:
		return "safety"
	}
}

func mapDeliveryMode(v string) string {
	switch v {
	case "Direct execution (implement immediately if confident)":
		return "direct"
	case "Hybrid (plan for complex, direct for simple)":
		return "hybrid"
	default:
		return "plan-first"
	}
}

func mapModelStrategy(v string) string {
	switch v {
	case "Single model (e.g., only Claude)":
		return "single"
	case "Round-robin (rotate models)":
		return "round-robin"
	case "Best available model":
		return "best-available"
	default:
		return "multi"
	}
}

// DefaultAnswersFromDiscovery creates non-interactive defaults from discovery.
func DefaultAnswersFromDiscovery(repoName string, report *profile.DiscoveryReport) *Answers {
	a := &Answers{
		ProjectName:      repoName,
		ProjectMaturity:  "existing codebase",
		ProjectType:      "service",
		TestingPolicy:    "TDD (tests required before implementation)",
		DependencyPolicy: "Only with explicit approval",
		ChangePolicy:     "Preserve structure, minimal changes",
		QualityPriority:  "Safety / minimal regression",
		DeliveryMode:     "Hybrid (plan for complex, direct for simple)",
		AIModelStrategy:  "Multi-model (different models for different tasks)",
		SensitiveAreas: []string{
			"Authentication / Authorization",
			"Database migrations",
			"Public API contracts",
		},
	}
	ApplyDiscoveryDefaults(a, report)
	return a
}
