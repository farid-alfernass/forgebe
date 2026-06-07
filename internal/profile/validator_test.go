package profile

import "testing"

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
}
