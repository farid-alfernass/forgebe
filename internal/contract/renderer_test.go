package contract

import (
	"strings"
	"testing"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestRender_Basic(t *testing.T) {
	p := profile.ProjectProfile{
		Metadata: profile.Metadata{
			ProfileID: "test123",
			RepoPath:  "/home/user/myproject",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Project: profile.Project{
			Name: "myproject",
		},
		Stack: profile.Stack{
			PrimaryLanguage: "go",
			Architecture:    "hexagonal",
			PackageManager:  "go modules",
			TestFramework:   "testing",
		},
		Policy: profile.Policy{
			Testing: profile.TestingPolicy{
				Strategy: "tdd",
				Required: true,
				Coverage: 80,
			},
			Dependencies: profile.DependencyPolicy{
				AllowAddition:   true,
				RequireApproval: true,
			},
			Changes: profile.ChangePolicy{
				PreserveStructure: true,
				AllowRefactor:     true,
			},
		},
	}

	out, err := Render(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checks := []string{
		"Engineering Contract",
		"myproject",
		"/home/user/myproject",
		"go",
		"hexagonal",
		"go modules",
	}

	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("expected %q in output, not found. Output:\n%s", c, out)
		}
	}
}

func TestRender_EmptyProfile(t *testing.T) {
	p := profile.ProjectProfile{}

	out, err := Render(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "Engineering Contract") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}
