package prompt

import (
	"fmt"
	"strings"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// Generator produces context-aware prompts for AI tools.
type Generator struct {
	profile *profile.ProjectProfile
}

// NewGenerator creates a new prompt generator.
func NewGenerator(p *profile.ProjectProfile) (*Generator, error) {
	if p == nil {
		return nil, fmt.Errorf("prompt: profile is nil")
	}
	return &Generator{profile: p}, nil
}

// Generate produces a prompt for the given mode and task.
func (g *Generator) Generate(mode, task string) (string, error) {
	switch mode {
	case "implement":
		return g.implementPrompt(task), nil
	case "review":
		return g.reviewPrompt(task), nil
	case "debug":
		return g.debugPrompt(task), nil
	case "plan":
		return g.planPrompt(task), nil
	default:
		return "", fmt.Errorf("prompt: unsupported mode %q (supported: implement, review, debug, plan)", mode)
	}
}

func (g *Generator) implementPrompt(task string) string {
	var buf strings.Builder
	buf.WriteString("# Implement Task\n\n")
	buf.WriteString(g.contextBlock())
	buf.WriteString("\n## Task\n")
	if task != "" {
		buf.WriteString(task)
		buf.WriteString("\n")
	}
	buf.WriteString("\n## Instructions\n")
	buf.WriteString("- Follow the project's architecture pattern.\n")
	buf.WriteString("- Write tests first (TDD) if testing strategy is tdd.\n")
	buf.WriteString("- Keep changes within source roots.\n")
	buf.WriteString("- Do NOT modify sensitive areas without explicit approval.\n")
	if g.profile.Policy.Dependencies.RequireApproval {
		buf.WriteString("- Do NOT add new dependencies without approval.\n")
	}
	buf.WriteString("- Ensure code compiles and tests pass before submitting.\n")
	return buf.String()
}

func (g *Generator) reviewPrompt(task string) string {
	var buf strings.Builder
	buf.WriteString("# Code Review\n\n")
	buf.WriteString(g.contextBlock())
	buf.WriteString("\n## Review Focus\n")
	if task != "" {
		buf.WriteString(task)
		buf.WriteString("\n\n")
	}
	buf.WriteString("## Checklist\n")
	buf.WriteString("- Does the code follow the project architecture?\n")
	buf.WriteString("- Are tests present and meaningful?\n")
	buf.WriteString("- Any security concerns in sensitive areas?\n")
	buf.WriteString("- Does it follow the naming and style conventions?\n")
	if g.profile.Policy.Quality.LintRequired {
		buf.WriteString("- Does it pass linting?\n")
	}
	if g.profile.Policy.Quality.FormatRequired {
		buf.WriteString("- Is it properly formatted?\n")
	}
	buf.WriteString("- Are error paths handled?\n")
	return buf.String()
}

func (g *Generator) debugPrompt(task string) string {
	var buf strings.Builder
	buf.WriteString("# Debug Task\n\n")
	buf.WriteString(g.contextBlock())
	buf.WriteString("\n## Problem\n")
	if task != "" {
		buf.WriteString(task)
		buf.WriteString("\n")
	}
	buf.WriteString("\n## Approach\n")
	buf.WriteString("1. Reproduce the issue (write a failing test if possible).\n")
	buf.WriteString("2. Identify the root cause — don't guess, verify.\n")
	buf.WriteString("3. Fix minimally — avoid unrelated changes.\n")
	buf.WriteString("4. Verify the fix passes the test.\n")
	buf.WriteString("5. Check for regressions in related code.\n")
	return buf.String()
}

func (g *Generator) planPrompt(task string) string {
	var buf strings.Builder
	buf.WriteString("# Implementation Plan\n\n")
	buf.WriteString(g.contextBlock())
	buf.WriteString("\n## Goal\n")
	if task != "" {
		buf.WriteString(task)
		buf.WriteString("\n")
	}
	buf.WriteString("\n## Planning Guidelines\n")
	buf.WriteString("- Break into bite-sized tasks (2-5 min each).\n")
	buf.WriteString("- Include exact file paths.\n")
	buf.WriteString("- Include TDD cycle for each task.\n")
	buf.WriteString("- Consider the architecture pattern.\n")
	buf.WriteString("- Identify sensitive areas that need approval.\n")
	buf.WriteString("- Estimate impact on existing tests.\n")
	return buf.String()
}

func (g *Generator) contextBlock() string {
	var buf strings.Builder
	buf.WriteString("## Project Context\n")
	buf.WriteString(fmt.Sprintf("- Project: %s\n", g.profile.Project.Name))
	buf.WriteString(fmt.Sprintf("- Language: %s\n", g.profile.Stack.PrimaryLanguage))

	if g.profile.Stack.Framework != "" {
		buf.WriteString(fmt.Sprintf("- Framework: %s\n", g.profile.Stack.Framework))
	}
	if g.profile.Stack.Architecture != "" {
		buf.WriteString(fmt.Sprintf("- Architecture: %s\n", g.profile.Stack.Architecture))
	}
	if g.profile.Policy.Testing.Strategy != "" {
		buf.WriteString(fmt.Sprintf("- Testing: %s\n", g.profile.Policy.Testing.Strategy))
	}
	if g.profile.Stack.PackageManager != "" {
		buf.WriteString(fmt.Sprintf("- Package Manager: %s\n", g.profile.Stack.PackageManager))
	}
	if len(g.profile.Areas.SensitiveAreas) > 0 {
		buf.WriteString(fmt.Sprintf("- Sensitive Areas: %s\n", strings.Join(g.profile.Areas.SensitiveAreas, ", ")))
	}
	if len(g.profile.Areas.SourceRoots) > 0 {
		buf.WriteString(fmt.Sprintf("- Source Roots: %s\n", strings.Join(g.profile.Areas.SourceRoots, ", ")))
	}
	if g.profile.Policy.Delivery.Priority != "" {
		buf.WriteString(fmt.Sprintf("- Priority: %s\n", g.profile.Policy.Delivery.Priority))
	}

	return buf.String()
}
