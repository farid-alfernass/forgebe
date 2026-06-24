// Package review compares a git diff against a ForgeBE project policy and
// produces awareness findings. It is AI-agnostic: it never inspects which
// tool produced the changes, only the changes themselves.
package review

import (
	"fmt"

	"github.com/faridtriwicaksono/forgebe/internal/git"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// Severity levels for a Finding.
const (
	SeverityFail = "FAIL"
	SeverityWarn = "WARN"
	SeverityInfo = "INFO"
)

// Finding is a single awareness observation about the diff.
type Finding struct {
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
	Path     string `json:"path,omitempty"`
	Message  string `json:"message"`
}

// Reviewer evaluates a set of file changes against a project policy.
type Reviewer struct {
	profile  *profile.ProjectProfile
	repoPath string
	changes  []git.FileChange
}

// NewReviewer builds a Reviewer. repoPath is the repository root, used by
// rules that need to read file contents (it may be "" for pure-path rules).
func NewReviewer(p *profile.ProjectProfile, repoPath string, changes []git.FileChange) (*Reviewer, error) {
	if p == nil {
		return nil, fmt.Errorf("review: profile is nil")
	}
	return &Reviewer{profile: p, repoPath: repoPath, changes: changes}, nil
}

// Run executes every rule and returns the aggregated findings.
func (r *Reviewer) Run() []Finding {
	var findings []Finding
	findings = append(findings, r.ruleForbiddenPath()...)
	findings = append(findings, r.ruleSensitiveArea()...)
	return findings
}
