package profile

import (
	"fmt"
	"os"
	"strings"
)

// ValidationIssue represents a single validation finding.
type ValidationIssue struct {
	Severity string `yaml:"severity" json:"severity"` // error, warning, info
	Field    string `yaml:"field" json:"field"`
	Message  string `yaml:"message" json:"message"`
}

// ValidationResult contains all validation findings.
type ValidationResult struct {
	Valid   bool              `yaml:"valid" json:"valid"`
	Issues  []ValidationIssue `yaml:"issues,omitempty" json:"issues,omitempty"`
	Summary string            `yaml:"summary" json:"summary"`
}

// ValidateProjectProfile validates a ProjectProfile and returns findings.
func ValidateProjectProfile(p *ProjectProfile) ValidationResult {
	result := ValidationResult{Valid: true}

	if p == nil {
		result.Valid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Severity: "error", Field: "root", Message: "Profile is nil",
		})
		result.Summary = "Profile is missing"
		return result
	}

	// Version
	if p.Version == "" {
		result.addIssue("error", "version", "Version is required")
	}

	// Metadata
	if p.Metadata.ProfileID == "" {
		result.addIssue("error", "metadata.profile_id", "Profile ID is required")
	}
	if p.Metadata.RepoPath == "" {
		result.addIssue("warning", "metadata.repo_path", "Repo path is not set")
	}

	// Project
	if p.Project.Name == "" {
		result.addIssue("warning", "project.name", "Project name is not set")
	}
	if len(p.Project.Languages) == 0 {
		result.addIssue("warning", "project.languages", "No languages declared")
	}

	// Stack
	if p.Stack.PrimaryLanguage == "" {
		result.addIssue("warning", "stack.primary_language", "Primary language not set")
	}
	if p.Stack.Architecture == "" {
		result.addIssue("warning", "stack.architecture", "Architecture not set")
	}

	// Policy
	if p.Policy.Testing.Strategy == "" {
		result.addIssue("warning", "policy.testing.strategy", "Testing strategy not set")
	}

	// Repo accessibility check
	if p.Metadata.RepoPath != "" {
		info, err := os.Stat(p.Metadata.RepoPath)
		if err != nil {
			result.addIssue("warning", "metadata.repo_path", fmt.Sprintf("Repo path not accessible: %s", err))
		} else if !info.IsDir() {
			result.addIssue("warning", "metadata.repo_path", "Repo path is not a directory")
		}
	}

	// Summary
	if len(result.Issues) == 0 {
		result.Summary = "All checks passed"
	} else {
		var errs, warns, infos int
		for _, issue := range result.Issues {
			switch issue.Severity {
			case "error":
				errs++
			case "warning":
				warns++
			case "info":
				infos++
			}
		}
		parts := []string{}
		if errs > 0 {
			parts = append(parts, fmt.Sprintf("%d errors", errs))
		}
		if warns > 0 {
			parts = append(parts, fmt.Sprintf("%d warnings", warns))
		}
		if infos > 0 {
			parts = append(parts, fmt.Sprintf("%d info", infos))
		}
		result.Summary = strings.Join(parts, ", ")
		if errs > 0 {
			result.Valid = false
			result.Summary = "Failed: " + result.Summary
		} else {
			result.Summary = "Passed with " + result.Summary
		}
	}

	return result
}

func (r *ValidationResult) addIssue(severity, field, message string) {
	r.Issues = append(r.Issues, ValidationIssue{
		Severity: severity, Field: field, Message: message,
	})
}
