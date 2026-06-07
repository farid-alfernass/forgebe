package adapters

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

type adapterSpec struct {
	Filename string
	Template string
}

var adapterTemplates = map[string]adapterSpec{
	"claude": {
		Filename: "CLAUDE.md",
		Template: `# Claude Code Context

Project: {{ .Project.Name }}
Primary language: {{ .Stack.PrimaryLanguage }}
Architecture: {{ .Stack.Architecture }}
Package manager: {{ .Stack.PackageManager }}
Testing strategy: {{ .Policy.Testing.Strategy }}
Dependency approval required: {{ .Policy.Dependencies.RequireApproval }}
Sensitive areas: {{ join .Areas.SensitiveAreas }}

Rules:
- Preserve existing structure unless approved.
- Prefer minimal, verifiable changes.
- Follow testing strategy above.
`,
	},
	"cursor": {
		Filename: ".cursorrules",
		Template: `Project={{ .Project.Name }}
Language={{ .Stack.PrimaryLanguage }}
Architecture={{ .Stack.Architecture }}
Testing={{ .Policy.Testing.Strategy }}
Rules=Preserve structure;Require approval for dependencies={{ .Policy.Dependencies.RequireApproval }}
`,
	},
	"copilot": {
		Filename: ".github/copilot-instructions.md",
		Template: `# Copilot Instructions

- Project: {{ .Project.Name }}
- Language: {{ .Stack.PrimaryLanguage }}
- Architecture: {{ .Stack.Architecture }}
- Testing: {{ .Policy.Testing.Strategy }}
- Sensitive Areas: {{ join .Areas.SensitiveAreas }}
`,
	},
	"hermes": {
		Filename: "AGENTS.md",
		Template: `# Hermes Agent Context

Project: {{ .Project.Name }}
Language: {{ .Stack.PrimaryLanguage }}
Architecture: {{ .Stack.Architecture }}
Delivery mode: {{ .Policy.Delivery.Mode }}
Priority: {{ .Policy.Delivery.Priority }}
Sensitive areas: {{ join .Areas.SensitiveAreas }}
`,
	},
	"generic": {
		Filename: "AI_CONTEXT.md",
		Template: `# AI Context

Project: {{ .Project.Name }}
Primary language: {{ .Stack.PrimaryLanguage }}
Framework: {{ .Stack.Framework }}
Architecture: {{ .Stack.Architecture }}
Testing strategy: {{ .Policy.Testing.Strategy }}
Dependency policy: approval required={{ .Policy.Dependencies.RequireApproval }}
Sensitive areas: {{ join .Areas.SensitiveAreas }}
`,
	},
}

// Render renders a full adapter file for a target tool.
func Render(tool string, p profile.ProjectProfile) (string, string, error) {
	spec, ok := adapterTemplates[strings.ToLower(tool)]
	if !ok {
		return "", "", fmt.Errorf("adapter: unsupported tool %q", tool)
	}
	out, err := execute(spec.Template, p)
	if err != nil {
		return "", "", err
	}
	return out, spec.Filename, nil
}

// RenderBrief renders a compact ready-to-paste brief.
func RenderBrief(tool string, p profile.ProjectProfile) (string, error) {
	_, _, err := Render(tool, p)
	if err != nil {
		return "", err
	}
	brief := fmt.Sprintf(
		"Project: %s\nLanguage: %s\nArchitecture: %s\nTesting: %s\nDependency approval required: %t\nSensitive areas: %s\n",
		p.Project.Name,
		p.Stack.PrimaryLanguage,
		p.Stack.Architecture,
		p.Policy.Testing.Strategy,
		p.Policy.Dependencies.RequireApproval,
		strings.Join(p.Areas.SensitiveAreas, ", "),
	)
	return brief, nil
}

func execute(tmplSrc string, p profile.ProjectProfile) (string, error) {
	funcMap := template.FuncMap{
		"join": func(items []string) string {
			if len(items) == 0 {
				return "none"
			}
			return strings.Join(items, ", ")
		},
	}
	tmpl, err := template.New("adapter").Funcs(funcMap).Parse(tmplSrc)
	if err != nil {
		return "", fmt.Errorf("adapter: parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("adapter: execute template: %w", err)
	}
	return buf.String(), nil
}
