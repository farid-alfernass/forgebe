package verify

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// Result represents a single verification check.
type Result struct {
	Check   string `json:"check"`
	Status  string `json:"status"` // PASS, WARN, FAIL
	Message string `json:"message"`
}

// Validator checks a project profile against engineering policies.
type Validator struct {
	profile *profile.ProjectProfile
}

// NewValidator creates a new validator.
func NewValidator(p *profile.ProjectProfile) (*Validator, error) {
	if p == nil {
		return nil, fmt.Errorf("verify: profile is nil")
	}
	return &Validator{profile: p}, nil
}

// Validate runs all checks against the profile.
func (v *Validator) Validate() []Result {
	var results []Result

	results = append(results, v.checkTestingPolicy())
	results = append(results, v.checkDependencyPolicy())
	results = append(results, v.checkChangePolicy())
	results = append(results, v.checkSensitiveAreas())
	results = append(results, v.checkSourceRoots())
	results = append(results, v.checkArchitecture())
	results = append(results, v.checkForbiddenPaths())

	return results
}

// checkTestingPolicy validates testing configuration.
func (v *Validator) checkTestingPolicy() Result {
	p := v.profile.Policy.Testing
	if p.Required && p.Strategy == "" {
		return Result{Check: "testing_policy", Status: "FAIL", Message: "Testing is required but no strategy defined"}
	}
	if p.Coverage < 0 || p.Coverage > 100 {
		return Result{Check: "testing_policy", Status: "WARN", Message: fmt.Sprintf("Coverage target %d%% is outside valid range", p.Coverage)}
	}
	return Result{Check: "testing_policy", Status: "PASS", Message: fmt.Sprintf("Testing is %s, strategy: %s, coverage target: %d%%", boolStr(p.Required), p.Strategy, p.Coverage)}
}

// checkDependencyPolicy validates dependency configuration.
func (v *Validator) checkDependencyPolicy() Result {
	p := v.profile.Policy.Dependencies
	if p.RequireApproval && v.profile.Stack.PackageManager == "" {
		return Result{Check: "dependency_policy", Status: "WARN", Message: "Dependency approval is required but no package manager detected"}
	}
	if len(p.Forbidden) > 0 {
		return Result{Check: "dependency_policy", Status: "PASS", Message: fmt.Sprintf("Forbidden packages defined: %s", strings.Join(p.Forbidden, ", "))}
	}
	return Result{Check: "dependency_policy", Status: "PASS", Message: fmt.Sprintf("Allow addition: %s, require approval: %s", boolStr(p.AllowAddition), boolStr(p.RequireApproval))}
}

// checkChangePolicy validates change configuration.
func (v *Validator) checkChangePolicy() Result {
	p := v.profile.Policy.Changes
	if p.PreserveStructure && !p.AllowCrossFile {
		return Result{Check: "change_policy", Status: "WARN", Message: "Structure preservation with restricted cross-file changes may limit refactoring"}
	}
	return Result{Check: "change_policy", Status: "PASS", Message: fmt.Sprintf("Structure: %s, cross-file: %s", boolStr(p.PreserveStructure), boolStr(p.AllowCrossFile))}
}

// checkSensitiveAreas validates sensitive areas definition.
func (v *Validator) checkSensitiveAreas() Result {
	areas := v.profile.Areas.SensitiveAreas
	if len(areas) == 0 {
		return Result{Check: "sensitive_areas", Status: "WARN", Message: "No sensitive areas defined; AI may operate on critical code without guardrails"}
	}
	return Result{Check: "sensitive_areas", Status: "PASS", Message: fmt.Sprintf("Sensitive areas: %s", strings.Join(areas, ", "))}
}

// checkSourceRoots validates source roots.
func (v *Validator) checkSourceRoots() Result {
	roots := v.profile.Areas.SourceRoots
	if len(roots) == 0 {
		return Result{Check: "source_roots", Status: "WARN", Message: "No source roots defined; AI context may miss key directories"}
	}
	return Result{Check: "source_roots", Status: "PASS", Message: fmt.Sprintf("Source roots: %s", strings.Join(roots, ", "))}
}

// checkArchitecture validates architecture configuration.
func (v *Validator) checkArchitecture() Result {
	arch := v.profile.Stack.Architecture
	if arch == "" {
		return Result{Check: "architecture", Status: "WARN", Message: "No architecture pattern defined; AI may produce inconsistent structure"}
	}
	return Result{Check: "architecture", Status: "PASS", Message: fmt.Sprintf("Architecture: %s", arch)}
}

// checkForbiddenPaths checks forbidden paths from change policy.
func (v *Validator) checkForbiddenPaths() Result {
	paths := v.profile.Policy.Changes.ForbiddenPaths
	if len(paths) == 0 {
		return Result{Check: "forbidden_paths", Status: "PASS", Message: "No forbidden paths defined"}
	}
	return Result{Check: "forbidden_paths", Status: "PASS", Message: fmt.Sprintf("Forbidden paths: %s", strings.Join(paths, ", "))}
}

// ResultsToJSON serializes results to JSON.
func ResultsToJSON(results []Result) string {
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "[]"
	}
	return string(b)
}

// ResultsToText formats results as human-readable text.
func ResultsToText(results []Result) string {
	pass := 0
	warn := 0
	fail := 0

	var buf strings.Builder
	buf.WriteString("ForgeBE Verify Report\n")
	buf.WriteString("=====================\n\n")

	for _, r := range results {
		statusIcon := "✅"
		statusPrefix := "PASS"
		switch r.Status {
		case "WARN":
			statusIcon = "⚠️"
			statusPrefix = "WARN"
			warn++
		case "FAIL":
			statusIcon = "❌"
			statusPrefix = "FAIL"
			fail++
		default:
			pass++
		}

		buf.WriteString(fmt.Sprintf("%s [%s] %s\n", statusIcon, statusPrefix, r.Check))
		buf.WriteString(fmt.Sprintf("   %s\n\n", r.Message))
	}

	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("Summary: %d PASS / %d WARN / %d FAIL\n", pass, warn, fail))
	if fail > 0 {
		buf.WriteString("Result: FAIL")
	} else if warn > 0 {
		buf.WriteString("Result: WARN")
	} else {
		buf.WriteString("Result: PASS")
	}
	buf.WriteString("\n")

	return buf.String()
}

// VerifyAgainstRepo runs checks against a real repo at the given path.
func VerifyAgainstRepo(repoPath string) ([]Result, error) {
	// Check repo path exists
	info, err := os.Stat(repoPath)
	if err != nil {
		return nil, fmt.Errorf("verify: cannot access repo path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("verify: %s is not a directory", repoPath)
	}

	return nil, fmt.Errorf("verify: repo-level verification not yet implemented")
}

func boolStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// Check name regex for validation
var checkNameRe = regexp.MustCompile(`^[a-z_]+$`)
