package contract

import (
	"path/filepath"
	"regexp"
	"strings"
)

// ExtractionResult captures policies inferred from markdown/instruction files.
type ExtractionResult struct {
	SourceFile                 string   `yaml:"source_file" json:"source_file"`
	IsAIInstructionFile        bool     `yaml:"is_ai_instruction_file" json:"is_ai_instruction_file"`
	TestingStrategy            string   `yaml:"testing_strategy,omitempty" json:"testing_strategy,omitempty"`
	DependencyApprovalRequired bool     `yaml:"dependency_approval_required" json:"dependency_approval_required"`
	PreserveStructure          bool     `yaml:"preserve_structure" json:"preserve_structure"`
	SensitiveAreas             []string `yaml:"sensitive_areas,omitempty" json:"sensitive_areas,omitempty"`
	RawSignals                 []string `yaml:"raw_signals,omitempty" json:"raw_signals,omitempty"`
}

var (
	reTDD             = regexp.MustCompile(`(?i)\b(tdd|tests? before implementation|write tests? before)\b`)
	reDepApproval     = regexp.MustCompile(`(?i)(do not add dependencies without explicit approval|require approval for new dependencies|dependencies? .*approval)`)
	rePreserveStruct  = regexp.MustCompile(`(?i)(preserve existing structure|minimal changes|do not restructure)`)
	reSensitiveBullet = regexp.MustCompile(`(?m)^\s*[-*]\s+(.+)$`)
)

// ExtractFromMarkdown extracts engineering contract hints from markdown or AI instruction files.
func ExtractFromMarkdown(filename, content string) ExtractionResult {
	result := ExtractionResult{SourceFile: filename}
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return result
	}

	base := strings.ToLower(filepath.Base(filename))
	switch base {
	case "claude.md", "agents.md", ".cursorrules", "ai_context.md", "copilot-instructions.md":
		result.IsAIInstructionFile = true
	}

	if reTDD.MatchString(content) {
		result.TestingStrategy = "tdd"
		result.RawSignals = append(result.RawSignals, "tdd")
	}
	if reDepApproval.MatchString(content) {
		result.DependencyApprovalRequired = true
		result.RawSignals = append(result.RawSignals, "dependency-approval")
	}
	if rePreserveStruct.MatchString(content) {
		result.PreserveStructure = true
		result.RawSignals = append(result.RawSignals, "preserve-structure")
	}

	lines := strings.Split(content, "\n")
	inSensitiveSection := false
	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))
		if strings.Contains(lower, "sensitive areas") || strings.Contains(lower, "critical areas") {
			inSensitiveSection = true
			continue
		}
		if inSensitiveSection {
			if strings.HasPrefix(lower, "##") || strings.HasPrefix(lower, "# ") {
				inSensitiveSection = false
				continue
			}
			if matches := reSensitiveBullet.FindStringSubmatch(line); len(matches) == 2 {
				result.SensitiveAreas = append(result.SensitiveAreas, strings.TrimSpace(matches[1]))
			}
		}
	}
	return result
}
