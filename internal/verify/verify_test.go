package verify

import (
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestNewValidator_InvalidProfile(t *testing.T) {
	_, err := NewValidator(nil)
	if err == nil {
		t.Fatal("expected error for nil profile")
	}
}

func TestValidate_DefaultProfile(t *testing.T) {
	p := profile.NewSampleProfile()
	v, err := NewValidator(&p)
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	results := v.Validate()
	if len(results) == 0 {
		t.Fatal("expected at least one check result")
	}
}

func TestValidate_HasTestingCheck(t *testing.T) {
	p := profile.NewSampleProfile()
	v, _ := NewValidator(&p)

	results := v.Validate()
	found := false
	for _, r := range results {
		if r.Check == "testing_policy" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected testing_policy check in results")
	}
}

func TestValidate_HasDependencyCheck(t *testing.T) {
	p := profile.NewSampleProfile()
	v, _ := NewValidator(&p)

	results := v.Validate()
	found := false
	for _, r := range results {
		if r.Check == "dependency_policy" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected dependency_policy check in results")
	}
}

func TestValidate_WithDefaultPasses(t *testing.T) {
	p := profile.NewSampleProfile()
	v, _ := NewValidator(&p)

	results := v.Validate()
	for _, r := range results {
		if r.Status == "FAIL" {
			t.Errorf("expected no FAIL results with default profile, got %s: %s", r.Check, r.Message)
		}
	}
}

func TestValidate_SensitiveAreasPresent(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Areas.SensitiveAreas = []string{"Authentication", "Authorization", "Database"}
	v, _ := NewValidator(&p)

	results := v.Validate()
	foundSensitiveCheck := false
	for _, r := range results {
		if r.Check == "sensitive_areas" {
			foundSensitiveCheck = true
			if r.Status != "PASS" {
				t.Errorf("expected PASS for sensitive areas check, got %s: %s", r.Status, r.Message)
			}
			break
		}
	}

	if !foundSensitiveCheck {
		t.Error("expected sensitive_areas check when areas are defined")
	}
}

func TestValidate_ChangePolicyWithRestrictedPaths(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Policy.Changes.ForbiddenPaths = []string{"vendor/", "node_modules/"}
	v, _ := NewValidator(&p)

	results := v.Validate()
	foundForbiddenCheck := false
	for _, r := range results {
		if r.Check == "forbidden_paths" {
			foundForbiddenCheck = true
			break
		}
	}

	if !foundForbiddenCheck {
		t.Error("expected forbidden_paths check when forbidden paths defined")
	}
}

func TestValidator_JSONOutput(t *testing.T) {
	p := profile.NewSampleProfile()
	v, _ := NewValidator(&p)

	results := v.Validate()
	jsonStr := ResultsToJSON(results)

	if len(jsonStr) == 0 {
		t.Fatal("expected non-empty JSON output")
	}

	if jsonStr[:1] != "[" {
		t.Error("expected JSON array output")
	}
}

func TestResultsToText(t *testing.T) {
	p := profile.NewSampleProfile()
	v, _ := NewValidator(&p)

	results := v.Validate()
	text := ResultsToText(results)

	if len(text) == 0 {
		t.Fatal("expected non-empty text output")
	}
	if !contains(text, "ForgeBE Verify Report") {
		t.Error("expected header in text output")
	}
	if !contains(text, "Summary:") {
		t.Error("expected summary in text output")
	}
}

func TestResultsToText_WithFailure(t *testing.T) {
	results := []Result{
		{Check: "test_check", Status: "FAIL", Message: "Something failed"},
	}
	text := ResultsToText(results)
	if !contains(text, "FAIL") {
		t.Error("expected FAIL in text output")
	}
	if !contains(text, "Result: FAIL") {
		t.Error("expected Result: FAIL at end")
	}
}

func TestResultsToText_WithWarn(t *testing.T) {
	results := []Result{
		{Check: "test_check", Status: "WARN", Message: "Warning"},
	}
	text := ResultsToText(results)
	if !contains(text, "WARN") {
		t.Error("expected WARN in text output")
	}
	if !contains(text, "Result: WARN") {
		t.Error("expected Result: WARN at end")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestValidate_EmptySensitiveAreas(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Areas.SensitiveAreas = nil
	v, _ := NewValidator(&p)

	results := v.Validate()
	for _, r := range results {
		if r.Check == "sensitive_areas" {
			if r.Status == "PASS" {
				t.Error("expected WARN for empty sensitive areas")
			}
			return
		}
	}
}
