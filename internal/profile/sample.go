package profile

import "time"

// NewSampleProfile returns a fully populated ProjectProfile for testing and documentation.
func NewSampleProfile() ProjectProfile {
	return ProjectProfile{
		Version: "1",
		Metadata: Metadata{
			ProfileID:  "sample_project_a1b2c3d4",
			RepoPath:   "/home/user/projects/sample",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			ForgebeVer: "dev",
			InitMethod: "guided",
		},
		Project: Project{
			Name:        "sample",
			Description: "A sample project for testing",
			Maturity:    "existing codebase",
			Type:        "service",
			Languages:   []string{"go"},
		},
		Stack: Stack{
			PrimaryLanguage: "go",
			Framework:       "gin",
			PackageManager:  "go modules",
			Architecture:    "clean/hexagonal",
			TestFramework:   "go test",
			CI:              []string{"github-actions"},
		},
		Policy: Policy{
			Testing:      TestingPolicy{Required: true, Strategy: "tdd", UnitRequired: true},
			Dependencies: DependencyPolicy{AllowAddition: false, RequireApproval: true},
			Changes:      ChangePolicy{PreserveStructure: true, AllowRefactor: false},
			Quality:      QualityPolicy{LintRequired: true, FormatRequired: true, ReviewRequired: true},
			Delivery:     DeliveryPolicy{Mode: "hybrid", Priority: "safety"},
		},
		Areas: Areas{
			SourceRoots:    []string{"cmd/", "internal/"},
			TestRoots:      []string{"internal/"},
			SensitiveAreas: []string{"Authentication", "Database migrations"},
		},
		AI: AI{
			InteractionMode: "hybrid",
			ModelStrategy:   "multi",
		},
	}
}
