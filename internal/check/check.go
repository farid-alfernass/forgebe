package check

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// Result represents the outcome of a health check.
type Result struct {
	Status      string            `json:"status"` // PASS, WARN, FAIL
	Checks      map[string]string `json:"checks"`
	TotalChecks int               `json:"total_checks"`
	Passed      int               `json:"passed"`
	Warned      int               `json:"warned"`
	Failed      int               `json:"failed"`
	Message     string            `json:"message"`
}

// Checker performs health checks on a project profile.
type Checker struct {
	profile *profile.ProjectProfile
}

// NewChecker creates a new checker.
func NewChecker(p *profile.ProjectProfile) (*Checker, error) {
	if p == nil {
		return nil, fmt.Errorf("check: profile is nil")
	}
	return &Checker{profile: p}, nil
}

// Check runs all health checks and aggregates results.
func (c *Checker) Check() *Result {
	checks := make(map[string]string)

	// Profile freshness
	freshness := c.checkProfileFreshness()
	checks["profile_freshness"] = freshness

	// Repository exists
	repoStatus := c.checkRepoExists()
	checks["repo_exists"] = repoStatus

	// Verify snapshot
	verifyStatus := c.checkVerify()
	checks["verify"] = verifyStatus

	// Sync status (simulated from profile metadata)
	syncStatus := c.checkSync()
	checks["sync_status"] = syncStatus

	// Count results
	passed := 0
	warned := 0
	failed := 0
	for _, s := range checks {
		switch s {
		case "PASS":
			passed++
		case "WARN":
			warned++
		case "FAIL":
			failed++
		}
	}

	overallStatus := "PASS"
	if failed > 0 {
		overallStatus = "FAIL"
	} else if warned > 0 {
		overallStatus = "WARN"
	}

	return &Result{
		Status:      overallStatus,
		Checks:      checks,
		TotalChecks: len(checks),
		Passed:      passed,
		Warned:      warned,
		Failed:      failed,
		Message:     fmt.Sprintf("%d passed, %d warnings, %d failures", passed, warned, failed),
	}
}

func (c *Checker) checkProfileFreshness() string {
	updated := c.profile.Metadata.UpdatedAt
	if updated.IsZero() {
		return "WARN"
	}

	age := time.Since(updated)
	switch {
	case age > 7*24*time.Hour:
		return "WARN"
	default:
		return "PASS"
	}
}

func (c *Checker) checkRepoExists() string {
	repoPath := c.profile.Metadata.RepoPath
	if repoPath == "" {
		return "FAIL"
	}
	info, err := os.Stat(repoPath)
	if err != nil || !info.IsDir() {
		return "FAIL"
	}
	return "PASS"
}

func (c *Checker) checkVerify() string {
	// Quick verification: check key fields are populated
	if c.profile.Stack.PrimaryLanguage == "" {
		return "WARN"
	}
	if c.profile.Policy.Testing.Required && c.profile.Policy.Testing.Strategy == "" {
		return "WARN"
	}
	return "PASS"
}

func (c *Checker) checkSync() string {
	// Sync is considered "synced" if profile has been updated recently
	if c.profile.Metadata.UpdatedAt.IsZero() {
		return "WARN"
	}
	if c.profile.Metadata.InitMethod == "" {
		return "WARN"
	}
	return "PASS"
}

// Text formats the result as human-readable text.
func (r *Result) Text() string {
	var buf strings.Builder
	buf.WriteString("ForgeBE Check\n")
	buf.WriteString("=============\n\n")

	// Sort check keys for deterministic output
	keys := make([]string, 0, len(r.Checks))
	for k := range r.Checks {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := r.Checks[k]
		icon := "✅"
		switch v {
		case "WARN":
			icon = "⚠️"
		case "FAIL":
			icon = "❌"
		}
		buf.WriteString(fmt.Sprintf("%s %s: %s\n", icon, k, v))
	}

	buf.WriteString("\n---\n")
	buf.WriteString(fmt.Sprintf("Result: %s (%s)\n", r.Status, r.Message))
	return buf.String()
}

// JSON formats the result as JSON.
func (r *Result) JSON() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
