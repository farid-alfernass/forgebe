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
