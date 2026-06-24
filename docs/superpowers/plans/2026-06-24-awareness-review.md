# Awareness Review (`forgebe review`) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `forgebe review` command that compares the git diff (the AI's output) against the project's ForgeBE policy and prints a prioritized "awareness report".

**Architecture:** A new `internal/git` package extracts changed files from git (working tree / staged / since-ref). A new `internal/review` package maps those changes against `profile.Policy` to produce `Finding`s, rendered as text or JSON. A new `internal/cli/review.go` wires it together, mirroring the existing `verify`/`check` cobra commands. Informational by default (exit 0); `--strict` exits non-zero on any FAIL.

**Tech Stack:** Go 1.23 standard library only (`os/exec`, `encoding/json`, `path/filepath`, `bufio`, `strings`). No new third-party dependencies.

## Global Constraints

- Module path: `github.com/faridtriwicaksono/forgebe` (use this in all imports).
- Go version floor: `go 1.23`.
- **AI-agnostic hard rule:** no file in `internal/review` or `internal/git` may import or reference any AI vendor (anthropic, claude, openai, gpt, copilot, cursor, gemini, hermes). Enforced by a guard test in Task 7.
- **No new dependencies** (CLAUDE.md requires dependency approval; this plan adds none — stdlib only).
- **TDD:** write the failing test first, watch it fail, then implement (CLAUDE.md testing strategy = tdd).
- Package test coverage must stay **≥ 80%** (existing SonarCloud gate).
- `gofmt -l .` clean and `go vet ./...` clean before every commit.
- Follow the existing cobra pattern in `internal/cli/verify.go` (use `resolveProjectIDArg`, `loadProfileByID`, `cmd.OutOrStdout()`, `SilenceUsage: true`).
- Severity strings are exactly `"FAIL"`, `"WARN"`, `"INFO"`.

---

## File Structure

| File | Responsibility |
|------|----------------|
| `internal/git/diff.go` | Run `git diff` and return `[]FileChange` for a `RangeSpec`; detect git repo. |
| `internal/git/diff_test.go` | Tests against a real temp git repo. |
| `internal/review/reviewer.go` | `Finding`, severity consts, `Reviewer`, `NewReviewer`, `Run` (dispatch). |
| `internal/review/rules.go` | Individual diff↔policy rules + path/file matching helpers. |
| `internal/review/report.go` | `Report`/`Summary` aggregation + `Text()`/`JSON()`. |
| `internal/review/reviewer_test.go` | Tests for `NewReviewer` + path-based rules. |
| `internal/review/rules_test.go` | Tests for dependency + missing_test rules. |
| `internal/review/report_test.go` | Tests for summary rule + report rendering. |
| `internal/review/agnostic_test.go` | AI-agnostic guard test. |
| `internal/cli/review.go` | `newReviewCmd()` cobra command wiring. |
| `internal/cli/review_test.go` | CLI-level test against a temp repo + profile. |
| `internal/cli/root.go` | Register `newReviewCmd()`. |
| `docs/cli-reference.md` | Document the `review` command. |
| `README.md` | Mention `forgebe review` in the workflow section. |

---

## Task 1: `internal/git` — diff extraction

**Files:**
- Create: `internal/git/diff.go`
- Test: `internal/git/diff_test.go`

**Interfaces:**
- Consumes: nothing (stdlib only).
- Produces:
  - `type FileChange struct { Path, OldPath, Status string; Added, Deleted int }`
  - `type RangeSpec struct { Staged bool; Since string }`
  - `func ChangedFiles(repoPath string, spec RangeSpec) ([]FileChange, error)`
  - `func IsRepo(repoPath string) bool`

- [ ] **Step 1: Write the failing test**

Create `internal/git/diff_test.go`:

```go
package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmd := exec.Command("git", "-C", dir, "init")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	return dir
}

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

func TestIsRepo(t *testing.T) {
	dir := initRepo(t)
	if !IsRepo(dir) {
		t.Error("expected IsRepo=true for a git repo")
	}
	if IsRepo(t.TempDir()) {
		t.Error("expected IsRepo=false for a non-repo dir")
	}
}

func TestChangedFiles_WorkingTree_TrackedAndUntracked(t *testing.T) {
	dir := initRepo(t)
	write(t, dir, "main.go", "package main\n")
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "-m", "init")

	// modify tracked file
	write(t, dir, "main.go", "package main\n\nfunc main() {}\n")
	// add a brand-new untracked file
	write(t, dir, "internal/new.go", "package internal\n")

	changes, err := ChangedFiles(dir, RangeSpec{})
	if err != nil {
		t.Fatalf("ChangedFiles: %v", err)
	}

	paths := map[string]FileChange{}
	for _, c := range changes {
		paths[c.Path] = c
	}
	if _, ok := paths["main.go"]; !ok {
		t.Errorf("expected main.go in changes, got %+v", changes)
	}
	if nf, ok := paths["internal/new.go"]; !ok {
		t.Errorf("expected untracked internal/new.go in changes, got %+v", changes)
	} else if nf.Status != "A" {
		t.Errorf("expected untracked file status A, got %q", nf.Status)
	}
}

func TestChangedFiles_Staged(t *testing.T) {
	dir := initRepo(t)
	write(t, dir, "a.go", "package a\n")
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "-m", "init")

	write(t, dir, "b.go", "package b\n")
	gitRun(t, dir, "add", "b.go")

	changes, err := ChangedFiles(dir, RangeSpec{Staged: true})
	if err != nil {
		t.Fatalf("ChangedFiles staged: %v", err)
	}
	found := false
	for _, c := range changes {
		if c.Path == "b.go" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected b.go in staged changes, got %+v", changes)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./internal/git/ -run TestChangedFiles -v`
Expected: FAIL — `undefined: ChangedFiles` / `undefined: RangeSpec` (package does not compile yet).

- [ ] **Step 3: Implement `internal/git/diff.go`**

Create `internal/git/diff.go`:

```go
// Package git extracts changed files from a git repository. It is
// deliberately AI-agnostic: it inspects the result of changes (the diff),
// never the agent that produced them.
package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// FileChange describes a single changed file from a git diff.
type FileChange struct {
	Path    string `json:"path"`
	OldPath string `json:"old_path,omitempty"`
	Status  string `json:"status"` // A, M, D
	Added   int    `json:"added"`
	Deleted int    `json:"deleted"`
}

// RangeSpec selects which diff to inspect.
type RangeSpec struct {
	Staged bool   // git diff --cached
	Since  string // git diff <ref>...HEAD; empty => working tree vs HEAD
}

// IsRepo reports whether repoPath is inside a git work tree.
func IsRepo(repoPath string) bool {
	out, err := runGit(repoPath, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// ChangedFiles returns the files changed in repoPath for the given range.
func ChangedFiles(repoPath string, spec RangeSpec) ([]FileChange, error) {
	switch {
	case spec.Staged:
		return diffRange(repoPath, "--cached")
	case spec.Since != "":
		return diffRange(repoPath, spec.Since+"...HEAD")
	default:
		var tracked []FileChange
		if hasHEAD(repoPath) {
			var err error
			tracked, err = diffRange(repoPath, "HEAD")
			if err != nil {
				return nil, err
			}
		}
		untracked, err := untrackedFiles(repoPath)
		if err != nil {
			return nil, err
		}
		return append(tracked, untracked...), nil
	}
}

func hasHEAD(repoPath string) bool {
	_, err := runGit(repoPath, "rev-parse", "--verify", "HEAD")
	return err == nil
}

// diffRange runs name-status + numstat for one diff target and merges them.
func diffRange(repoPath string, target ...string) ([]FileChange, error) {
	nameArgs := append([]string{"diff", "--name-status", "--no-renames"}, target...)
	ns, err := runGit(repoPath, nameArgs...)
	if err != nil {
		return nil, err
	}
	numArgs := append([]string{"diff", "--numstat", "--no-renames"}, target...)
	num, err := runGit(repoPath, numArgs...)
	if err != nil {
		return nil, err
	}
	return mergeDiff(ns, num), nil
}

func mergeDiff(nameStatus, numStat string) []FileChange {
	order := []string{}
	byPath := map[string]*FileChange{}
	get := func(path string) *FileChange {
		if fc, ok := byPath[path]; ok {
			return fc
		}
		fc := &FileChange{Path: path}
		byPath[path] = fc
		order = append(order, path)
		return fc
	}

	sc := bufio.NewScanner(strings.NewReader(nameStatus))
	for sc.Scan() {
		fields := strings.SplitN(strings.TrimRight(sc.Text(), "\n"), "\t", 2)
		if len(fields) != 2 || strings.TrimSpace(fields[1]) == "" {
			continue
		}
		status := strings.TrimSpace(fields[0])
		fc := get(strings.TrimSpace(fields[1]))
		if status != "" {
			fc.Status = status[:1]
		}
	}

	sc = bufio.NewScanner(strings.NewReader(numStat))
	for sc.Scan() {
		fields := strings.SplitN(strings.TrimRight(sc.Text(), "\n"), "\t", 3)
		if len(fields) != 3 || strings.TrimSpace(fields[2]) == "" {
			continue
		}
		fc := get(strings.TrimSpace(fields[2]))
		fc.Added = atoiSafe(fields[0])
		fc.Deleted = atoiSafe(fields[1])
	}

	result := make([]FileChange, 0, len(order))
	for _, p := range order {
		result = append(result, *byPath[p])
	}
	return result
}

func untrackedFiles(repoPath string) ([]FileChange, error) {
	out, err := runGit(repoPath, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	var result []FileChange
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		path := strings.TrimSpace(sc.Text())
		if path == "" {
			continue
		}
		result = append(result, FileChange{Path: path, Status: "A", Added: countLines(repoPath, path)})
	}
	return result, nil
}

func countLines(repoPath, rel string) int {
	data, err := os.ReadFile(filepath.Join(repoPath, rel))
	if err != nil || len(data) == 0 {
		return 0
	}
	return bytes.Count(data, []byte{'\n'})
}

func atoiSafe(s string) int {
	s = strings.TrimSpace(s)
	if s == "-" || s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func runGit(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", repoPath}, args...)...)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(errBuf.String()))
	}
	return out.String(), nil
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./internal/git/ -v`
Expected: PASS (TestIsRepo, TestChangedFiles_WorkingTree_TrackedAndUntracked, TestChangedFiles_Staged).

- [ ] **Step 5: Verify formatting/vet and commit**

Run: `gofmt -l internal/git/ && go vet ./internal/git/`
Expected: no output from gofmt, no vet warnings.

```bash
git add internal/git/diff.go internal/git/diff_test.go
git commit -m "feat: add internal/git diff extraction package"
```

---

## Task 2: `internal/review` core + path-based rules

**Files:**
- Create: `internal/review/reviewer.go`
- Create: `internal/review/rules.go`
- Test: `internal/review/reviewer_test.go`

**Interfaces:**
- Consumes: `git.FileChange` (Task 1); `profile.ProjectProfile`, `profile.Policy`, `profile.Areas` (existing).
- Produces:
  - `type Finding struct { Rule, Severity, Path, Message string }`
  - consts `SeverityFail = "FAIL"`, `SeverityWarn = "WARN"`, `SeverityInfo = "INFO"`
  - `func NewReviewer(p *profile.ProjectProfile, repoPath string, changes []git.FileChange) (*Reviewer, error)`
  - `func (r *Reviewer) Run() []Finding`
  - helpers `pathMatches(p, pattern string) bool`, `looksLikePath(s string) bool`

- [ ] **Step 1: Write the failing test**

Create `internal/review/reviewer_test.go`:

```go
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
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./internal/review/ -run "TestNewReviewer|TestRule" -v`
Expected: FAIL — undefined `NewReviewer`, `Finding`, `SeverityFail` (package does not compile).

- [ ] **Step 3: Implement `internal/review/reviewer.go`**

Create `internal/review/reviewer.go`:

```go
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
```

- [ ] **Step 4: Implement `internal/review/rules.go` (path rules + helpers)**

Create `internal/review/rules.go`:

```go
package review

import (
	"fmt"
	"path/filepath"
	"strings"
)

// pathMatches reports whether file path p matches pattern. It supports glob
// matching (on the full slash path and on the basename) and falls back to a
// case-insensitive substring match so loose keywords also match.
func pathMatches(p, pattern string) bool {
	if pattern == "" {
		return false
	}
	p = filepath.ToSlash(p)
	pattern = filepath.ToSlash(pattern)
	if ok, _ := filepath.Match(pattern, p); ok {
		return true
	}
	if ok, _ := filepath.Match(pattern, filepath.Base(p)); ok {
		return true
	}
	needle := strings.ToLower(strings.Trim(pattern, "*/"))
	return needle != "" && strings.Contains(strings.ToLower(p), needle)
}

// looksLikePath reports whether s is meant as a path/glob rather than prose.
// Descriptive sensitive areas like "Authentication / Authorization" return false.
func looksLikePath(s string) bool {
	return strings.ContainsAny(s, "/*.") && !strings.Contains(s, " ")
}

func (r *Reviewer) ruleForbiddenPath() []Finding {
	var out []Finding
	for _, c := range r.changes {
		for _, pat := range r.profile.Policy.Changes.ForbiddenPaths {
			if pathMatches(c.Path, pat) {
				out = append(out, Finding{
					Rule:     "forbidden_path",
					Severity: SeverityFail,
					Path:     c.Path,
					Message:  fmt.Sprintf("%s matches forbidden path %q", c.Path, pat),
				})
				break
			}
		}
	}
	return out
}

func (r *Reviewer) ruleSensitiveArea() []Finding {
	var patterns []string
	patterns = append(patterns, r.profile.Policy.Changes.RequireApproval...)
	patterns = append(patterns, r.profile.Areas.PublicAPIs...)
	for _, s := range r.profile.Areas.SensitiveAreas {
		if looksLikePath(s) {
			patterns = append(patterns, s)
		}
	}

	var out []Finding
	for _, c := range r.changes {
		for _, pat := range patterns {
			if pathMatches(c.Path, pat) {
				out = append(out, Finding{
					Rule:     "sensitive_area",
					Severity: SeverityWarn,
					Path:     c.Path,
					Message:  fmt.Sprintf("%s touches a sensitive/approval area (%q) — confirm you reviewed this", c.Path, pat),
				})
				break
			}
		}
	}
	return out
}
```

- [ ] **Step 5: Run the tests to verify they pass**

Run: `go test ./internal/review/ -run "TestNewReviewer|TestRule" -v`
Expected: PASS (TestNewReviewer_NilProfile, TestRule_ForbiddenPath, TestRule_SensitiveArea).

- [ ] **Step 6: Verify and commit**

Run: `gofmt -l internal/review/ && go vet ./internal/review/`
Expected: clean.

```bash
git add internal/review/reviewer.go internal/review/rules.go internal/review/reviewer_test.go
git commit -m "feat: add review core with forbidden_path and sensitive_area rules"
```

---

## Task 3: dependency rules

**Files:**
- Modify: `internal/review/reviewer.go` (add rule call to `Run`)
- Modify: `internal/review/rules.go` (add `ruleDependency` + `manifestFiles` + `manifestContains`, and `os` import)
- Test: `internal/review/rules_test.go`

**Interfaces:**
- Consumes: `profile.Policy.Dependencies` (`AllowAddition`, `RequireApproval`, `Forbidden`); reads files under `repoPath`.
- Produces: findings with rule `dependency_added` (FAIL if additions disallowed, else WARN if approval required) and `dependency_forbidden` (FAIL).

- [ ] **Step 1: Write the failing test**

Create `internal/review/rules_test.go`:

```go
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
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./internal/review/ -run TestRule_Dependency -v`
Expected: FAIL — `dependency_added` / `dependency_forbidden` findings are never produced (rules not implemented).

- [ ] **Step 3: Update the `rules.go` import block to add `os`**

In `internal/review/rules.go`, change the import block at the top to:

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)
```

- [ ] **Step 4: Append `ruleDependency` + helpers to `rules.go`**

Append to `internal/review/rules.go`:

```go
// manifestFiles maps a dependency manifest filename to its ecosystem.
var manifestFiles = map[string]string{
	"go.mod":           "go",
	"package.json":     "node",
	"requirements.txt": "python",
	"pyproject.toml":   "python",
	"Cargo.toml":       "rust",
	"pom.xml":          "java",
	"build.gradle":     "java",
}

func (r *Reviewer) ruleDependency() []Finding {
	dep := r.profile.Policy.Dependencies
	var out []Finding
	for _, c := range r.changes {
		base := filepath.Base(filepath.ToSlash(c.Path))
		if _, ok := manifestFiles[base]; !ok {
			continue
		}
		if c.Status != "A" && c.Added == 0 {
			continue
		}

		switch {
		case !dep.AllowAddition:
			out = append(out, Finding{
				Rule:     "dependency_added",
				Severity: SeverityFail,
				Path:     c.Path,
				Message:  fmt.Sprintf("%s changed but dependency additions are not allowed", c.Path),
			})
		case dep.RequireApproval:
			out = append(out, Finding{
				Rule:     "dependency_added",
				Severity: SeverityWarn,
				Path:     c.Path,
				Message:  fmt.Sprintf("%s changed — new dependencies require your approval", c.Path),
			})
		}

		for _, forbidden := range dep.Forbidden {
			if forbidden == "" {
				continue
			}
			if manifestContains(r.repoPath, c.Path, forbidden) {
				out = append(out, Finding{
					Rule:     "dependency_forbidden",
					Severity: SeverityFail,
					Path:     c.Path,
					Message:  fmt.Sprintf("forbidden dependency %q present in %s", forbidden, c.Path),
				})
			}
		}
	}
	return out
}

func manifestContains(repoPath, rel, token string) bool {
	if repoPath == "" {
		return false
	}
	data, err := os.ReadFile(filepath.Join(repoPath, rel))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), token)
}
```

- [ ] **Step 5: Wire `ruleDependency` into `Run`**

In `internal/review/reviewer.go`, update `Run`:

```go
func (r *Reviewer) Run() []Finding {
	var findings []Finding
	findings = append(findings, r.ruleForbiddenPath()...)
	findings = append(findings, r.ruleSensitiveArea()...)
	findings = append(findings, r.ruleDependency()...)
	return findings
}
```

- [ ] **Step 6: Run the tests to verify they pass**

Run: `go test ./internal/review/ -run TestRule_Dependency -v`
Expected: PASS (all four dependency tests).

- [ ] **Step 7: Verify and commit**

Run: `gofmt -l internal/review/ && go vet ./internal/review/`
Expected: clean.

```bash
git add internal/review/rules.go internal/review/reviewer.go internal/review/rules_test.go
git commit -m "feat: add dependency_added and dependency_forbidden review rules"
```

---

## Task 4: `missing_test` rule

**Files:**
- Modify: `internal/review/reviewer.go` (add rule call to `Run`)
- Modify: `internal/review/rules.go` (add `ruleMissingTest` + helpers)
- Test: `internal/review/rules_test.go` (append)

**Interfaces:**
- Consumes: `profile.Policy.Testing` (`Required`, `UnitRequired`), `profile.Areas` (`SourceRoots`, `TestRoots`).
- Produces: findings with rule `missing_test` (WARN) for a changed source file with no test changed alongside it.

- [ ] **Step 1: Write the failing test**

Append to `internal/review/rules_test.go`:

```go
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
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./internal/review/ -run TestRule_MissingTest -v`
Expected: FAIL — `missing_test` findings never produced.

- [ ] **Step 3: Append `ruleMissingTest` + helpers to `rules.go`**

Append to `internal/review/rules.go`:

```go
var codeExtensions = map[string]bool{
	".go": true, ".js": true, ".ts": true, ".jsx": true, ".tsx": true,
	".py": true, ".rb": true, ".java": true, ".rs": true, ".kt": true, ".php": true,
}

func isSourceFile(p string) bool {
	return codeExtensions[strings.ToLower(filepath.Ext(p))]
}

func isTestFile(p string) bool {
	base := strings.ToLower(filepath.Base(p))
	return strings.HasSuffix(base, "_test.go") ||
		strings.HasSuffix(base, "_test.py") ||
		strings.HasPrefix(base, "test_") ||
		strings.Contains(base, ".test.") ||
		strings.Contains(base, ".spec.")
}

func underAnyRoot(p string, roots []string) bool {
	p = filepath.ToSlash(p)
	for _, root := range roots {
		root = strings.Trim(filepath.ToSlash(root), "/")
		if root == "" {
			continue
		}
		if p == root || strings.HasPrefix(p, root+"/") {
			return true
		}
	}
	return false
}

// testKey reduces a path to a key shared by a source file and its sibling test
// (same directory, same base name). MVP heuristic — not call-graph analysis.
func testKey(p string) string {
	p = filepath.ToSlash(p)
	dir := filepath.Dir(p)
	name := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
	name = strings.TrimSuffix(name, "_test")
	name = strings.TrimSuffix(name, ".test")
	name = strings.TrimSuffix(name, ".spec")
	name = strings.TrimPrefix(name, "test_")
	return dir + "/" + name
}

func (r *Reviewer) ruleMissingTest() []Finding {
	t := r.profile.Policy.Testing
	if !t.Required && !t.UnitRequired {
		return nil
	}

	testChanged := map[string]bool{}
	for _, c := range r.changes {
		if isTestFile(c.Path) || underAnyRoot(c.Path, r.profile.Areas.TestRoots) {
			testChanged[testKey(c.Path)] = true
		}
	}

	var out []Finding
	for _, c := range r.changes {
		if c.Status == "D" || !isSourceFile(c.Path) || isTestFile(c.Path) {
			continue
		}
		if len(r.profile.Areas.SourceRoots) > 0 && !underAnyRoot(c.Path, r.profile.Areas.SourceRoots) {
			continue
		}
		if testChanged[testKey(c.Path)] {
			continue
		}
		out = append(out, Finding{
			Rule:     "missing_test",
			Severity: SeverityWarn,
			Path:     c.Path,
			Message:  fmt.Sprintf("%s changed without a corresponding test in the same change", c.Path),
		})
	}
	return out
}
```

- [ ] **Step 4: Wire `ruleMissingTest` into `Run`**

In `internal/review/reviewer.go`, update `Run`:

```go
func (r *Reviewer) Run() []Finding {
	var findings []Finding
	findings = append(findings, r.ruleForbiddenPath()...)
	findings = append(findings, r.ruleSensitiveArea()...)
	findings = append(findings, r.ruleDependency()...)
	findings = append(findings, r.ruleMissingTest()...)
	return findings
}
```

- [ ] **Step 5: Run the tests to verify they pass**

Run: `go test ./internal/review/ -run TestRule_MissingTest -v`
Expected: PASS (all three missing_test tests).

- [ ] **Step 6: Verify and commit**

Run: `gofmt -l internal/review/ && go vet ./internal/review/`
Expected: clean.

```bash
git add internal/review/rules.go internal/review/reviewer.go internal/review/rules_test.go
git commit -m "feat: add missing_test review rule"
```

---

## Task 5: `summary` rule + report rendering

**Files:**
- Modify: `internal/review/reviewer.go` (add `ruleSummary` call to `Run`)
- Modify: `internal/review/rules.go` (add `ruleSummary` + `topDir`, and `sort` import)
- Create: `internal/review/report.go`
- Test: `internal/review/report_test.go`

**Interfaces:**
- Consumes: `[]Finding` (Task 2).
- Produces:
  - `type Summary struct { Fail, Warn, Info int }`
  - `type Report struct { Range string; Findings []Finding; Summary Summary }`
  - `func NewReport(findings []Finding, rangeLabel string) Report`
  - `func (rep Report) Text() string`, `func (rep Report) JSON() string`, `func (rep Report) HasFail() bool`

- [ ] **Step 1: Write the failing test**

Create `internal/review/report_test.go`:

```go
package review

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/git"
	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestRule_Summary(t *testing.T) {
	p := profile.NewSampleProfile()
	changes := []git.FileChange{
		{Path: "internal/a.go", Status: "M", Added: 3, Deleted: 1},
		{Path: "cmd/b.go", Status: "A", Added: 10, Deleted: 0},
	}
	r, _ := NewReviewer(&p, "", changes)
	got := findingsByRule(r.Run(), "summary")
	if len(got) != 1 || got[0].Severity != SeverityInfo {
		t.Fatalf("expected 1 INFO summary, got %+v", got)
	}
	if !strings.Contains(got[0].Message, "2 files changed") {
		t.Errorf("summary message missing count: %q", got[0].Message)
	}
}

func TestNewReport_CountsAndJSON(t *testing.T) {
	findings := []Finding{
		{Rule: "forbidden_path", Severity: SeverityFail, Path: "x", Message: "m"},
		{Rule: "missing_test", Severity: SeverityWarn, Path: "y", Message: "m"},
		{Rule: "summary", Severity: SeverityInfo, Message: "1 files changed (+1 / -0)"},
	}
	rep := NewReport(findings, "working tree (vs HEAD)")
	if rep.Summary.Fail != 1 || rep.Summary.Warn != 1 || rep.Summary.Info != 1 {
		t.Fatalf("unexpected summary: %+v", rep.Summary)
	}
	if !rep.HasFail() {
		t.Error("expected HasFail=true")
	}

	var decoded Report
	if err := json.Unmarshal([]byte(rep.JSON()), &decoded); err != nil {
		t.Fatalf("JSON not valid: %v", err)
	}
	if len(decoded.Findings) != 3 || decoded.Range != "working tree (vs HEAD)" {
		t.Errorf("round-trip mismatch: %+v", decoded)
	}

	text := rep.Text()
	if !strings.Contains(text, "ForgeBE Awareness Review") || !strings.Contains(text, "REVIEW NEEDED") {
		t.Errorf("text missing header/result: %q", text)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./internal/review/ -run "TestRule_Summary|TestNewReport" -v`
Expected: FAIL — `NewReport`, `Report`, `Summary` undefined; summary finding not produced.

- [ ] **Step 3: Update the `rules.go` import block to add `sort`**

In `internal/review/rules.go`, change the import block to:

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)
```

- [ ] **Step 4: Append `ruleSummary` + `topDir` to `rules.go`**

Append to `internal/review/rules.go`:

```go
func topDir(p string) string {
	p = filepath.ToSlash(p)
	if i := strings.Index(p, "/"); i >= 0 {
		return p[:i]
	}
	return "."
}

func (r *Reviewer) ruleSummary() []Finding {
	if len(r.changes) == 0 {
		return []Finding{{Rule: "summary", Severity: SeverityInfo, Message: "no changes detected"}}
	}
	totalAdded, totalDeleted := 0, 0
	groups := map[string]int{}
	for _, c := range r.changes {
		totalAdded += c.Added
		totalDeleted += c.Deleted
		groups[topDir(c.Path)]++
	}

	parts := make([]string, 0, len(groups))
	for dir, n := range groups {
		parts = append(parts, fmt.Sprintf("%s:%d", dir, n))
	}
	sort.Strings(parts)

	msg := fmt.Sprintf("%d files changed (+%d / -%d)  [%s]",
		len(r.changes), totalAdded, totalDeleted, strings.Join(parts, " "))
	return []Finding{{Rule: "summary", Severity: SeverityInfo, Message: msg}}
}
```

- [ ] **Step 5: Wire `ruleSummary` into `Run`**

In `internal/review/reviewer.go`, update `Run` (summary last so it appears at the end):

```go
func (r *Reviewer) Run() []Finding {
	var findings []Finding
	findings = append(findings, r.ruleForbiddenPath()...)
	findings = append(findings, r.ruleSensitiveArea()...)
	findings = append(findings, r.ruleDependency()...)
	findings = append(findings, r.ruleMissingTest()...)
	findings = append(findings, r.ruleSummary()...)
	return findings
}
```

- [ ] **Step 6: Implement `internal/review/report.go`**

Create `internal/review/report.go`:

```go
package review

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Summary counts findings by severity.
type Summary struct {
	Fail int `json:"fail"`
	Warn int `json:"warn"`
	Info int `json:"info"`
}

// Report is the full awareness result for one diff.
type Report struct {
	Range    string    `json:"range"`
	Findings []Finding `json:"findings"`
	Summary  Summary   `json:"summary"`
}

// NewReport aggregates findings into a Report.
func NewReport(findings []Finding, rangeLabel string) Report {
	var s Summary
	for _, f := range findings {
		switch f.Severity {
		case SeverityFail:
			s.Fail++
		case SeverityWarn:
			s.Warn++
		case SeverityInfo:
			s.Info++
		}
	}
	return Report{Range: rangeLabel, Findings: findings, Summary: s}
}

// HasFail reports whether any FAIL finding exists.
func (rep Report) HasFail() bool { return rep.Summary.Fail > 0 }

// JSON renders the report as indented JSON.
func (rep Report) JSON() string {
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}

// Text renders the report as a human-facing awareness report.
func (rep Report) Text() string {
	var b strings.Builder
	b.WriteString("ForgeBE Awareness Review\n")
	b.WriteString("========================\n")
	b.WriteString("Range: " + rep.Range + "\n\n")
	for _, f := range rep.Findings {
		b.WriteString(fmt.Sprintf("%s %-5s %-20s %s\n", icon(f.Severity), f.Severity, f.Rule, f.Message))
	}
	b.WriteString("\n---\n")
	result := "OK"
	if rep.Summary.Fail > 0 {
		result = "REVIEW NEEDED"
	} else if rep.Summary.Warn > 0 {
		result = "REVIEW SUGGESTED"
	}
	b.WriteString(fmt.Sprintf("%d need your awareness  •  Result: %s\n",
		rep.Summary.Fail+rep.Summary.Warn, result))
	return b.String()
}

func icon(severity string) string {
	switch severity {
	case SeverityFail:
		return "❌"
	case SeverityWarn:
		return "⚠️ "
	case SeverityInfo:
		return "ℹ️ "
	default:
		return "• "
	}
}
```

- [ ] **Step 7: Run the tests to verify they pass**

Run: `go test ./internal/review/ -v`
Expected: PASS (all review tests, including summary + report).

- [ ] **Step 8: Verify and commit**

Run: `gofmt -l internal/review/ && go vet ./internal/review/`
Expected: clean.

```bash
git add internal/review/rules.go internal/review/reviewer.go internal/review/report.go internal/review/report_test.go
git commit -m "feat: add summary rule and awareness report rendering"
```

---

## Task 6: CLI `forgebe review` command

**Files:**
- Create: `internal/cli/review.go`
- Modify: `internal/cli/root.go` (register after `newCheckCmd()` at line 34)
- Test: `internal/cli/review_test.go`

**Interfaces:**
- Consumes: `git.IsRepo`, `git.ChangedFiles`, `git.RangeSpec` (Task 1); `review.NewReviewer`, `review.NewReport` (Tasks 2-5); existing `resolveProjectIDArg`, `loadProfileByID`.
- Produces: cobra command `review` registered on the root command.

- [ ] **Step 1: Write the failing test**

Create `internal/cli/review_test.go`:

```go
package cli

import (
	"bytes"
	"testing"
)

func TestReviewCmd_Registered(t *testing.T) {
	root := NewRootCmd()
	found := false
	for _, c := range root.Commands() {
		if c.Name() == "review" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected 'review' command to be registered")
	}
}

func TestReviewCmd_HelpRuns(t *testing.T) {
	root := NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"review", "--help"})
	if err := root.Execute(); err != nil {
		t.Fatalf("review --help failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("awareness")) {
		t.Errorf("help text should mention awareness, got: %s", buf.String())
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./internal/cli/ -run TestReviewCmd -v`
Expected: FAIL — `review` command not registered (`newReviewCmd` undefined).

- [ ] **Step 3: Implement `internal/cli/review.go`**

Create `internal/cli/review.go`:

```go
package cli

import (
	"fmt"
	"os"

	"github.com/faridtriwicaksono/forgebe/internal/git"
	"github.com/faridtriwicaksono/forgebe/internal/review"
	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {
	var (
		jsonOutput bool
		staged     bool
		since      string
		strict     bool
	)

	cmd := &cobra.Command{
		Use:   "review [project-id]",
		Short: "Review code changes against your policy (awareness report)",
		Long: `Compare the git diff (what your AI produced) against this project's
ForgeBE policy and print an awareness report: what changed, and what needs
your conscious review.

Informational by default (exit 0). Use --strict to exit non-zero on any
FAIL finding, for pre-commit hooks and CI. Works with any AI tool — it
inspects the diff, not the agent.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectID, err := resolveProjectIDArg(args)
			if err != nil {
				return fmt.Errorf("cannot resolve project: %w", err)
			}

			proj, err := loadProfileByID(projectID)
			if err != nil {
				return fmt.Errorf("project %q not found: %w", projectID, err)
			}

			repoPath := proj.Metadata.RepoPath
			if repoPath == "" {
				repoPath, _ = os.Getwd()
			}
			if !git.IsRepo(repoPath) {
				return fmt.Errorf("review: %s is not a git repository", repoPath)
			}

			spec := git.RangeSpec{Staged: staged, Since: since}
			changes, err := git.ChangedFiles(repoPath, spec)
			if err != nil {
				return fmt.Errorf("review: %w", err)
			}

			reviewer, err := review.NewReviewer(proj, repoPath, changes)
			if err != nil {
				return err
			}

			report := review.NewReport(reviewer.Run(), rangeLabel(spec))
			if jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), report.JSON())
			} else {
				fmt.Fprint(cmd.OutOrStdout(), report.Text())
			}

			if strict && report.HasFail() {
				os.Exit(1)
			}
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	cmd.Flags().BoolVar(&staged, "staged", false, "Review staged changes (git diff --cached)")
	cmd.Flags().StringVar(&since, "since", "", "Review changes since a git ref (e.g. main)")
	cmd.Flags().BoolVar(&strict, "strict", false, "Exit non-zero if any FAIL finding")

	return cmd
}

func rangeLabel(spec git.RangeSpec) string {
	switch {
	case spec.Staged:
		return "staged (vs HEAD)"
	case spec.Since != "":
		return "since " + spec.Since
	default:
		return "working tree (vs HEAD)"
	}
}
```

- [ ] **Step 4: Register the command in `root.go`**

In `internal/cli/root.go`, add the `newReviewCmd()` line between `newCheckCmd()` and `newPromptCmd()`:

```go
	cmd.AddCommand(newCheckCmd())
	cmd.AddCommand(newReviewCmd())
	cmd.AddCommand(newPromptCmd())
```

- [ ] **Step 5: Run the tests to verify they pass**

Run: `go test ./internal/cli/ -run TestReviewCmd -v`
Expected: PASS (TestReviewCmd_Registered, TestReviewCmd_HelpRuns).

- [ ] **Step 6: Build the binary and smoke-test**

Run: `go build -o bin/forgebe ./cmd/forgebe && ./bin/forgebe review --help`
Expected: build succeeds; help text prints and mentions "awareness".

- [ ] **Step 7: Verify and commit**

Run: `gofmt -l internal/cli/ && go vet ./internal/cli/`
Expected: clean.

```bash
git add internal/cli/review.go internal/cli/root.go internal/cli/review_test.go
git commit -m "feat: add forgebe review CLI command"
```

---

## Task 7: AI-agnostic guard test + docs + final verification

**Files:**
- Create: `internal/review/agnostic_test.go`
- Modify: `docs/cli-reference.md`
- Modify: `README.md`

**Interfaces:**
- Consumes: nothing new.
- Produces: a guard test; documentation.

- [ ] **Step 1: Write the AI-agnostic guard test**

Create `internal/review/agnostic_test.go`:

```go
package review

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoAIVendorReferences enforces ForgeBE's AI-agnostic guarantee: the
// review and git packages must never reference a specific AI vendor.
func TestNoAIVendorReferences(t *testing.T) {
	vendors := []string{"anthropic", "openai", "copilot", "cursor", "gemini", "hermes"}
	dirs := []string{".", "../git"}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read %s: %v", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
				continue
			}
			if strings.HasSuffix(e.Name(), "_test.go") {
				continue // this test names vendors on purpose
			}
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatalf("read %s: %v", e.Name(), err)
			}
			lower := strings.ToLower(string(data))
			for _, v := range vendors {
				if strings.Contains(lower, v) {
					t.Errorf("%s/%s references AI vendor %q — review/git must stay AI-agnostic", dir, e.Name(), v)
				}
			}
		}
	}
}
```

- [ ] **Step 2: Run the guard test**

Run: `go test ./internal/review/ -run TestNoAIVendorReferences -v`
Expected: PASS (no vendor names in non-test source of `internal/review` or `internal/git`).

Note: "claude" is intentionally omitted from the vendor list — the list targets vendor SDK/tool identifiers that must not leak into these two packages. If a future edit imports a vendor, this test fails.

- [ ] **Step 3: Document the command in `docs/cli-reference.md`**

Append a new section to the end of `docs/cli-reference.md`:

````markdown
## `forgebe review`

Compare the current git diff against your project policy and print an
awareness report — what changed and what needs your conscious review.
Works with any AI tool because it inspects the diff, not the agent.

```bash
forgebe review                 # working tree vs HEAD (informational, exit 0)
forgebe review --staged        # staged changes (for pre-commit)
forgebe review --since main    # changes since a ref (for CI/PR)
forgebe review --json          # machine-readable output
forgebe review --strict        # exit 1 if any FAIL finding (gate mode)
```

Rules checked (mapped to your profile policy):

| Rule | Policy source | Severity |
|------|---------------|----------|
| `forbidden_path` | `policy.changes.forbidden_paths` | FAIL |
| `sensitive_area` | `areas.sensitive_areas` / `areas.public_apis` / `policy.changes.require_approval` | WARN |
| `dependency_added` | `policy.dependencies.allow_addition` / `require_approval` | FAIL/WARN |
| `dependency_forbidden` | `policy.dependencies.forbidden` | FAIL |
| `missing_test` | `policy.testing.required` + `areas.test_roots` | WARN |
| `summary` | — | INFO |
````

- [ ] **Step 4: Mention `review` in `README.md`**

In `README.md`, under the "Workflow Harian dengan AI" section, add this block after the TDD-AI Loop description:

````markdown
### Awareness Review

Setelah AI mengubah kode, jaga agar kamu tetap *aware*:

```bash
# Lihat apa yang berubah & apa yang melenceng dari policy-mu
forgebe review

# Mode gate untuk pre-commit / CI
forgebe review --strict
```

`forgebe review` membandingkan git diff dengan policy project — bekerja dengan
AI tool apa pun (Claude, Cursor, Copilot, Hermes), karena yang diperiksa adalah
perubahannya, bukan agennya.
````

- [ ] **Step 5: Full test suite + coverage + vet + build**

Run:
```bash
go test ./... -cover
go vet ./...
gofmt -l .
go build -o bin/forgebe ./cmd/forgebe
```
Expected: all tests PASS; `internal/git` and `internal/review` coverage ≥ 80%; no vet warnings; `gofmt -l` prints nothing; binary builds.

- [ ] **Step 6: End-to-end smoke test**

Run (from the repo root, which is a git repo):
```bash
# Ensure a profile exists for this repo (creates one if missing)
./bin/forgebe init || true
# Make a throwaway change, then review it
printf '\n// scratch\n' >> internal/cli/review.go
./bin/forgebe review
git checkout -- internal/cli/review.go
```
Expected: `forgebe review` prints the "ForgeBE Awareness Review" report listing `internal/cli/review.go` in the summary and exits 0.

- [ ] **Step 7: Commit**

```bash
git add internal/review/agnostic_test.go docs/cli-reference.md README.md
git commit -m "test: add AI-agnostic guard; docs: document forgebe review"
```

---

## Self-Review

**Spec coverage:**
- Goal (git-diff → policy → awareness report): Tasks 1-6. ✓
- Posture (informational default / `--strict` gate): Task 6 (`strict && HasFail` → `os.Exit(1)`). ✓
- AI-agnostic hard rule: enforced by Task 7 guard test; no vendor imports anywhere. ✓
- Rule set (forbidden_path, sensitive_area, dependency_added, dependency_forbidden, missing_test, summary): Tasks 2-5. ✓
- Output text + JSON, flags (`--staged`, `--since`, `--json`, `--strict`), exit codes: Task 6 + report.go. ✓
- Architecture (`internal/git`, `internal/review`, `internal/cli/review.go`, root registration): Tasks 1, 2-5, 6. ✓
- Testing strategy (TDD, temp repo for git, synthetic changes for rules, ≥80% coverage): every task + Task 7 Step 5. ✓
- Non-goals (no arch import analysis, no provenance, no hook/Action, no auto-fix): none implemented — respected. ✓
- Acceptance criteria 1-7: covered by Task 6 smoke + Task 7 Steps 5-6 + guard test. ✓

**Type consistency:** `FileChange`, `RangeSpec`, `ChangedFiles`, `IsRepo` (git); `Finding`, `SeverityFail/Warn/Info`, `NewReviewer(p, repoPath, changes)`, `Run`, `Report`, `Summary`, `NewReport`, `Text`, `JSON`, `HasFail` (review) — used identically across Tasks 1-7. `Run` is extended additively in Tasks 3/4/5 with the exact full body shown each time. ✓

**Placeholder scan:** no TBD/TODO; every code step shows complete code; commands have expected output. ✓
```
