package review

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
		if isTestFile(c.Path) || (underAnyRoot(c.Path, r.profile.Areas.TestRoots) && !isSourceFile(c.Path)) {
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
