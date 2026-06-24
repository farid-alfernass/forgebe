package review

import (
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/git"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func findingsByRule(fs []Finding, rule string) []Finding {
	var out []Finding
	for _, f := range fs {
		if f.Rule == rule {
			out = append(out, f)
		}
	}
	return out
}

func TestNewReviewer_NilProfile(t *testing.T) {
	if _, err := NewReviewer(nil, "", nil); err == nil {
		t.Fatal("expected error for nil profile")
	}
}

func TestRule_ForbiddenPath(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Changes.ForbiddenPaths = []string{"internal/auth/*"}
	changes := []git.FileChange{
		{Path: "internal/auth/token.go", Status: "M", Added: 3},
		{Path: "internal/util/str.go", Status: "M", Added: 1},
	}
	r, err := NewReviewer(&p, "", changes)
	if err != nil {
		t.Fatal(err)
	}
	got := findingsByRule(r.Run(), "forbidden_path")
	if len(got) != 1 {
		t.Fatalf("expected 1 forbidden_path finding, got %d: %+v", len(got), got)
	}
	if got[0].Severity != SeverityFail || got[0].Path != "internal/auth/token.go" {
		t.Errorf("unexpected finding: %+v", got[0])
	}
}

func TestRule_SensitiveArea(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Changes.RequireApproval = []string{"internal/payment/*"}
	p.Areas.SensitiveAreas = []string{"Authentication / Authorization"} // descriptive, must be ignored
	changes := []git.FileChange{
		{Path: "internal/payment/charge.go", Status: "M", Added: 5},
	}
	r, _ := NewReviewer(&p, "", changes)
	got := findingsByRule(r.Run(), "sensitive_area")
	if len(got) != 1 {
		t.Fatalf("expected 1 sensitive_area finding, got %d: %+v", len(got), got)
	}
	if got[0].Severity != SeverityWarn {
		t.Errorf("expected WARN, got %q", got[0].Severity)
	}
}
