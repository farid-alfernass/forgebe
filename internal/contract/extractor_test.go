package contract

import (
	"strings"
	"testing"
)

func TestExtractFromMarkdown_FindsPoliciesAndSensitiveAreas(t *testing.T) {
	input := `# Project Rules

## Testing Policy
- Strategy: tdd
- Required: true

## Dependency Rules
- Require approval for new dependencies

## Sensitive Areas
- Authentication / Authorization
- Database migrations
`

	result := ExtractFromMarkdown("README.md", input)

	if result.TestingStrategy != "tdd" {
		t.Fatalf("expected testing strategy tdd, got %q", result.TestingStrategy)
	}
	if !result.DependencyApprovalRequired {
		t.Fatal("expected dependency approval to be required")
	}
	if len(result.SensitiveAreas) != 2 {
		t.Fatalf("expected 2 sensitive areas, got %+v", result.SensitiveAreas)
	}
}

func TestExtractFromMarkdown_DetectsAIInstructionFiles(t *testing.T) {
	input := `# CLAUDE.md
Do not add dependencies without explicit approval.
Always write tests before implementation.
Preserve existing structure.
`

	result := ExtractFromMarkdown("CLAUDE.md", input)

	if !result.IsAIInstructionFile {
		t.Fatal("expected file to be recognized as AI instruction file")
	}
	if !result.DependencyApprovalRequired {
		t.Fatal("expected dependency approval to be required")
	}
	if !result.PreserveStructure {
		t.Fatal("expected preserve structure to be true")
	}
	if result.TestingStrategy == "" {
		t.Fatal("expected testing strategy to be inferred")
	}
}

func TestExtractFromMarkdown_IgnoresEmptyContent(t *testing.T) {
	result := ExtractFromMarkdown("README.md", "   \n\t ")
	if strings.TrimSpace(result.SourceFile) == "" {
		t.Fatal("expected source file to be populated")
	}
	if result.TestingStrategy != "" || len(result.SensitiveAreas) != 0 {
		t.Fatalf("expected empty extraction result, got %+v", result)
	}
}
