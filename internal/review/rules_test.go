package review

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/git"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestRule_DependencyAdded_NotAllowed(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Dependencies.AllowAddition = false
	changes := []git.FileChange{{Path: "go.mod", Status: "M", Added: 2}}
	r, _ := NewReviewer(&p, "", changes)
	got := findingsByRule(r.Run(), "dependency_added")
	if len(got) != 1 || got[0].Severity != SeverityFail {
		t.Fatalf("expected 1 FAIL dependency_added, got %+v", got)
	}
}

func TestRule_DependencyAdded_RequiresApproval(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Dependencies.AllowAddition = true
	p.Policy.Dependencies.RequireApproval = true
	changes := []git.FileChange{{Path: "package.json", Status: "M", Added: 1}}
	r, _ := NewReviewer(&p, "", changes)
	got := findingsByRule(r.Run(), "dependency_added")
	if len(got) != 1 || got[0].Severity != SeverityWarn {
		t.Fatalf("expected 1 WARN dependency_added, got %+v", got)
	}
}

func TestRule_DependencyForbidden(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("require github.com/evil/pkg v1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := profile.NewSampleProfile()
	p.Policy.Dependencies.AllowAddition = true
	p.Policy.Dependencies.Forbidden = []string{"github.com/evil/pkg"}
	changes := []git.FileChange{{Path: "go.mod", Status: "M", Added: 1}}
	r, _ := NewReviewer(&p, dir, changes)
	got := findingsByRule(r.Run(), "dependency_forbidden")
	if len(got) != 1 || got[0].Severity != SeverityFail {
		t.Fatalf("expected 1 FAIL dependency_forbidden, got %+v", got)
	}
}

func TestRule_DependencyAdded_IgnoresNonManifest(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Dependencies.AllowAddition = false
	changes := []git.FileChange{{Path: "internal/foo.go", Status: "M", Added: 9}}
	r, _ := NewReviewer(&p, "", changes)
	if got := findingsByRule(r.Run(), "dependency_added"); len(got) != 0 {
		t.Fatalf("expected no dependency findings for non-manifest, got %+v", got)
	}
}

func TestRule_MissingTest_Flagged(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Testing.Required = true
	p.Areas.SourceRoots = []string{"internal"}
	p.Areas.TestRoots = []string{"internal"}
	changes := []git.FileChange{
		{Path: "internal/payment/charge.go", Status: "A", Added: 20},
	}
	r, _ := NewReviewer(&p, "", changes)
	got := findingsByRule(r.Run(), "missing_test")
	if len(got) != 1 || got[0].Severity != SeverityWarn {
		t.Fatalf("expected 1 WARN missing_test, got %+v", got)
	}
}

func TestRule_MissingTest_SatisfiedBySiblingTest(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Testing.Required = true
	p.Areas.SourceRoots = []string{"internal"}
	p.Areas.TestRoots = []string{"internal"}
	changes := []git.FileChange{
		{Path: "internal/payment/charge.go", Status: "M", Added: 10},
		{Path: "internal/payment/charge_test.go", Status: "M", Added: 15},
	}
	r, _ := NewReviewer(&p, "", changes)
	if got := findingsByRule(r.Run(), "missing_test"); len(got) != 0 {
		t.Fatalf("expected no missing_test when sibling test changed, got %+v", got)
	}
}

func TestRule_MissingTest_DisabledWhenTestingNotRequired(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Testing.Required = false
	p.Policy.Testing.UnitRequired = false
	changes := []git.FileChange{{Path: "internal/foo.go", Status: "A", Added: 5}}
	r, _ := NewReviewer(&p, "", changes)
	if got := findingsByRule(r.Run(), "missing_test"); len(got) != 0 {
		t.Fatalf("expected no missing_test when testing not required, got %+v", got)
	}
}

func TestRule_MissingTest_TestFileNotFlagged(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Testing.Required = true
	p.Areas.SourceRoots = []string{"internal"}
	p.Areas.TestRoots = []string{"internal"}
	changes := []git.FileChange{
		{Path: "internal/foo/bar_test.go", Status: "A", Added: 10},
	}
	r, _ := NewReviewer(&p, "", changes)
	if got := findingsByRule(r.Run(), "missing_test"); len(got) != 0 {
		t.Fatalf("a *_test.go file must never be flagged as missing a test, got %+v", got)
	}
}

func TestRule_MissingTest_UnitRequiredOnly(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Testing.Required = false
	p.Policy.Testing.UnitRequired = true
	p.Areas.SourceRoots = []string{"internal"}
	p.Areas.TestRoots = []string{"internal"}
	changes := []git.FileChange{
		{Path: "internal/payment/charge.go", Status: "A", Added: 20},
	}
	r, _ := NewReviewer(&p, "", changes)
	got := findingsByRule(r.Run(), "missing_test")
	if len(got) != 1 || got[0].Severity != SeverityWarn {
		t.Fatalf("expected 1 WARN missing_test when only UnitRequired, got %+v", got)
	}
}
