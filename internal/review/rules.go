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
