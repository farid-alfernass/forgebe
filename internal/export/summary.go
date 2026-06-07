package export

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

const summaryTemplate = `# ForgeBE Project Summary

## Overview
- **Project**: {{ .Project.Name }}
- **Profile ID**: {{ .Metadata.ProfileID }}
- **Repository**: {{ .Metadata.RepoPath }}
- **Created**: {{ .Metadata.CreatedAt.Format "2006-01-02 15:04" }}
- **Updated**: {{ .Metadata.UpdatedAt.Format "2006-01-02 15:04" }}

## Stack
- **Primary Language**: {{ .Stack.PrimaryLanguage }}
- **Framework**: {{ .Stack.Framework }}
- **Architecture**: {{ .Stack.Architecture }}
- **Package Manager**: {{ .Stack.PackageManager }}
- **Test Framework**: {{ .Stack.TestFramework }}
- **CI**: {{ formatSlice .Stack.CI }}

## Policies

### Testing
- **Strategy**: {{ .Policy.Testing.Strategy }}
- **Required**: {{ .Policy.Testing.Required }}
- **Coverage Target**: {{ .Policy.Testing.Coverage }}%

### Dependencies
- **Allow Addition**: {{ .Policy.Dependencies.AllowAddition }}
- **Require Approval**: {{ .Policy.Dependencies.RequireApproval }}
- **Forbidden**: {{ formatSlice .Policy.Dependencies.Forbidden }}

### Changes
- **Preserve Structure**: {{ .Policy.Changes.PreserveStructure }}
- **Allow Refactor**: {{ .Policy.Changes.AllowRefactor }}
- **Allow Cross-File**: {{ .Policy.Changes.AllowCrossFile }}

### Delivery
- **Mode**: {{ .Policy.Delivery.Mode }}
- **Priority**: {{ .Policy.Delivery.Priority }}

## Sensitive Areas
{{- if .Areas.SensitiveAreas }}
{{ range .Areas.SensitiveAreas }}- {{ . }}
{{ end -}}
{{- else }}
- None declared
{{- end }}

## AI Configuration
- **Interaction Mode**: {{ .AI.InteractionMode }}
- **Model Strategy**: {{ .AI.ModelStrategy }}
`

func RenderSummary(p profile.ProjectProfile) (string, error) {
	funcMap := template.FuncMap{
		"formatSlice": func(items []string) string {
			if len(items) == 0 {
				return "none"
			}
			return strings.Join(items, ", ")
		},
	}
	tmpl, err := template.New("summary").Funcs(funcMap).Parse(summaryTemplate)
	if err != nil {
		return "", fmt.Errorf("export: parse summary template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("export: execute summary template: %w", err)
	}
	return strings.TrimSpace(buf.String()) + "\n", nil
}
