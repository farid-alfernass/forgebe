package profile

import (
	"os"
	"testing"
)

func TestValidateProjectProfile_Valid(t *testing.T) {
	p := NewSampleProfile()
	result := ValidateProjectProfile(&p)
	if !result.Valid {
		t.Fatalf("expected valid profile, got issues: %+v", result.Issues)
	}
}

func TestValidateProjectProfile_Nil(t *testing.T) {
	result := ValidateProjectProfile(nil)
	if result.Valid {
		t.Fatal("expected invalid for nil profile")
	}
}

func TestValidateProjectProfile_MissingFields(t *testing.T) {
	p := ProjectProfile{}
	result := ValidateProjectProfile(&p)
	if result.Valid {
		t.Fatal("expected invalid for empty profile")
	}
	// Check for version error
	found := false
	for _, issue := range result.Issues {
		if issue.Field == "version" && issue.Severity == "error" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected version error issue")
	}
}

func TestValidateProjectProfile_RepoPathChecks(t *testing.T) {
	p := NewSampleProfile()

	// Non-existent path
	p.Metadata.RepoPath = "/non/existent/path/for/forgebe/test"
	result := ValidateProjectProfile(&p)
	found := false
	for _, issue := range result.Issues {
		if issue.Field == "metadata.repo_path" && issue.Severity == "warning" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for non-existent repo path")
	}

	// Not a directory
	tmpFile, err := os.CreateTemp("", "forgebe-not-a-dir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	p.Metadata.RepoPath = tmpFile.Name()
	result = ValidateProjectProfile(&p)
	found = false
	for _, issue := range result.Issues {
		if issue.Field == "metadata.repo_path" && issue.Message == "Repo path is not a directory" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for repo path being a file")
	}
}

func TestValidationResult_SummaryFormatting(t *testing.T) {
	// Only warnings
	p := NewSampleProfile()
	p.Project.Name = "" // warning
	result := ValidateProjectProfile(&p)
	if !result.Valid {
		t.Error("profile should be valid with only warnings")
	}
	if result.Summary == "" || result.Summary == "All checks passed" {
		t.Errorf("unexpected summary: %s", result.Summary)
	}

	// Mixed errors and warnings
	p.Version = "" // error
	result = ValidateProjectProfile(&p)
	if result.Valid {
		t.Error("profile should be invalid with errors")
	}
	if !result.Valid || result.Summary[:7] == "Passed" {
		// should start with Failed
	}
}
