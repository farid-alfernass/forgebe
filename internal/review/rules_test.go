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
