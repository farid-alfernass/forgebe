package review

import (
	"fmt"
	"os"
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
