package contract

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

const defaultTemplate = `# Engineering Contract

## Project Identity
- Name: {{ .Project.Name }}
- Repository: {{ .Metadata.RepoPath }}
- Primary language: {{ .Stack.PrimaryLanguage }}
- Architecture: {{ .Stack.Architecture }}
- Package manager: {{ .Stack.PackageManager }}

## Delivery Rules
- Testing strategy: {{ .Policy.Testing.Strategy }}
- Dependency policy: AllowAdd={{ .Policy.Dependencies.AllowAddition }}, ReqApproval={{ .Policy.Dependencies.RequireApproval }}
- Change policy: PreserveStruct={{ .Policy.Changes.PreserveStructure }}, AllowRefactor={{ .Policy.Changes.AllowRefactor }}
- Quality priority: {{ .Policy.Delivery.Priority }}
- Delivery mode: {{ .Policy.Delivery.Mode }}
- AI model strategy: {{ .AI.ModelStrategy }}

## Sensitive Areas
{{- if .Areas.SensitiveAreas }}
{{- range .Areas.SensitiveAreas }}
- {{ . }}
{{- end }}
{{- else }}
- No sensitive areas declared yet.
{{- end }}

## Guardrails
- Preserve existing repository structure unless explicitly approved.
- Do not add new dependencies without following the dependency policy.
- Prefer minimal, verifiable changes.
- New business logic should be covered by tests according to policy.
- Treat auth, billing, migrations, infra, and public APIs conservatively.
`

// Render returns the engineering contract markdown for a profile.
func Render(p profile.ProjectProfile) (string, error) {
	tmpl, err := template.New("contract").Parse(defaultTemplate)
	if err != nil {
		return "", fmt.Errorf("contract: parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("contract: render template: %w", err)
	}
	return strings.TrimSpace(buf.String()) + "\n", nil
}
