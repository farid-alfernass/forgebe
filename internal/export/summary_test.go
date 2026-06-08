package export

import (
	"strings"
	"testing"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestRenderSummary_Basic(t *testing.T) {
	p := profile.ProjectProfile{
		Metadata: profile.Metadata{
			ProfileID: "test123",
			RepoPath:  "/home/user/project",
			CreatedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC),
		},
		Project: profile.Project{
			Name: "myproject",
		},
		Stack: profile.Stack{
			PrimaryLanguage: "go",
			Framework:       "cobra",
			Architecture:    "hexagonal",
			PackageManager:  "go modules",
			TestFramework:   "testing",
			CI:              []string{"github-actions"},
		},
		Policy: profile.Policy{
			Testing: profile.TestingPolicy{
				Strategy: "tdd",
				Required: true,
				Coverage: 80,
			},
			Dependencies: profile.DependencyPolicy{
				AllowAddition:   true,
				RequireApproval: false,
				Forbidden:       []string{"unsafe-lib"},
			},
			Changes: profile.ChangePolicy{
				PreserveStructure: true,
				AllowRefactor:     true,
				AllowCrossFile:    true,
			},
			Delivery: profile.DeliveryPolicy{
				Mode:     "ci-cd",
				Priority: "quality",
			},
		},
		Areas: profile.Areas{
			SensitiveAreas: []string{"auth", "payments"},
		},
		AI: profile.AI{
			InteractionMode: "hybrid",
			ModelStrategy:   "single",
		},
	}

	out, err := RenderSummary(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checks := []string{
		"ForgeBE Project Summary",
		"myproject",
		"test123",
		"go",
		"cobra",
		"hexagonal",
		"github-actions",
		"tdd",
		"80%",
		"unsafe-lib",
		"auth",
		"payments",
		"hybrid",
		"single",
	}

	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("expected %q in output, not found. Output:\n%s", c, out)
		}
	}
}

func TestRenderSummary_EmptySlices(t *testing.T) {
	p := profile.ProjectProfile{
		Metadata: profile.Metadata{
			ProfileID: "empty",
			CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Project: profile.Project{
			Name: "empty-project",
		},
	}

	out, err := RenderSummary(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "None declared") {
		t.Errorf("expected 'None declared' for empty sensitive areas, got:\n%s", out)
	}
	if !strings.Contains(out, "none") {
		t.Errorf("expected 'none' for empty slices, got:\n%s", out)
	}
}
