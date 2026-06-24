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
