package check

import (
	"testing"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestNewChecker_InvalidProfile(t *testing.T) {
	_, err := NewChecker(nil)
	if err == nil {
		t.Fatal("expected error for nil profile")
	}
}

func TestChecker_HappyPath(t *testing.T) {
	p := profile.NewSampleProfile()
	c, err := NewChecker(&p)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	result := c.Check()
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestChecker_ExitCodeZeroOnAllPass(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Metadata.RepoPath = t.TempDir() // Use real path so repo_exists passes
	c, _ := NewChecker(&p)

	result := c.Check()
	if result.Status == "FAIL" {
		t.Errorf("expected no FAIL with sample profile (real path), got %s", result.Status)
	}
}

func TestChecker_ReportsMetrics(t *testing.T) {
	p := profile.NewSampleProfile()
	c, _ := NewChecker(&p)

	result := c.Check()
	if result.Checks == nil {
		t.Fatal("expected checks map")
	}
	if result.TotalChecks == 0 {
		t.Fatal("expected non-zero check count")
	}
}

func TestChecker_ProfileFreshness(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Metadata.UpdatedAt = time.Now().Add(-72 * time.Hour) // 3 days old
	c, _ := NewChecker(&p)

	result := c.Check()
	freshCheck, ok := result.Checks["profile_freshness"]
	if !ok {
		t.Fatal("expected profile_freshness check")
	}
	if freshCheck != "PASS" && freshCheck != "WARN" {
		t.Errorf("expected PASS or WARN for 3-day-old profile, got %s", freshCheck)
	}
}

func TestChecker_StaleProfile(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Metadata.UpdatedAt = time.Now().Add(-8 * 24 * time.Hour) // 8 days old
	p.Metadata.RepoPath = t.TempDir()                          // real path
	c, _ := NewChecker(&p)

	result := c.Check()
	freshCheck := result.Checks["profile_freshness"]
	if freshCheck != "WARN" {
		t.Errorf("expected WARN for 8-day-old profile freshness, got %s", freshCheck)
	}
}

func TestChecker_VerifyIntegration(t *testing.T) {
	p := profile.NewSampleProfile()
	c, _ := NewChecker(&p)

	result := c.Check()
	if _, ok := result.Checks["verify"]; !ok {
		t.Fatal("expected verify integration check")
	}
}

func TestChecker_RepositoryCheck(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Metadata.RepoPath = "/tmp"
	c, _ := NewChecker(&p)

	result := c.Check()
	repoCheck, ok := result.Checks["repo_exists"]
	if !ok {
		t.Fatal("expected repo_exists check")
	}
	if repoCheck != "PASS" {
		t.Errorf("expected PASS for /tmp repo path, got %s", repoCheck)
	}
}

func TestChecker_MissingRepoPath(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Metadata.RepoPath = "/nonexistent-path-xyz-abc"
	c, _ := NewChecker(&p)

	result := c.Check()
	repoCheck, ok := result.Checks["repo_exists"]
	if !ok {
		t.Fatal("expected repo_exists check")
	}
	if repoCheck != "FAIL" {
		t.Errorf("expected FAIL for non-existent repo path, got %s", repoCheck)
	}
}

func TestChecker_CheckCount(t *testing.T) {
	p := profile.NewSampleProfile()
	c, _ := NewChecker(&p)

	result := c.Check()
	if result.TotalChecks < 4 {
		t.Errorf("expected at least 4 checks, got %d", result.TotalChecks)
	}
}

func TestChecker_TextOutput(t *testing.T) {
	p := profile.NewSampleProfile()
	c, _ := NewChecker(&p)

	result := c.Check()
	text := result.Text()
	if len(text) == 0 {
		t.Fatal("expected non-empty text output")
	}
}

func TestChecker_JSONOutput(t *testing.T) {
	p := profile.NewSampleProfile()
	c, _ := NewChecker(&p)

	result := c.Check()
	jsonStr := result.JSON()
	if len(jsonStr) == 0 {
		t.Fatal("expected non-empty JSON output")
	}
}

func TestChecker_SyncCheck(t *testing.T) {
	p := profile.NewSampleProfile()
	c, _ := NewChecker(&p)

	result := c.Check()
	if _, ok := result.Checks["sync_status"]; !ok {
		t.Fatal("expected sync_status check")
	}
}
