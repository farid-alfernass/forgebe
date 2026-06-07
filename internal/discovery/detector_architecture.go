package discovery

import (
	"strings"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func DetectArchitecture(root string, dirs []string) profile.ArchitectureDetection {
	set := map[string]bool{}
	for _, d := range dirs {
		set[strings.ToLower(d)] = true
	}

	layers := []string{}
	patterns := []string{}
	style := "unknown"
	confidence := "low"

	if hasAny(set, "domain", "application", "infrastructure", "interface") {
		style = "clean-hexagonal"
		confidence = "high"
		layers = append(layers, "domain", "application", "infrastructure", "interface")
		patterns = append(patterns, "clean-architecture", "hexagonal")
		return profile.ArchitectureDetection{Style: style, Confidence: confidence, Layers: layers, Patterns: patterns}
	}

	if hasAny(set, "controllers", "services", "repositories") || hasAny(set, "controller", "service", "repository") {
		style = "layered-service-repository"
		confidence = "high"
		layers = append(layers, "controller", "service", "repository")
		patterns = append(patterns, "layered")
		return profile.ArchitectureDetection{Style: style, Confidence: confidence, Layers: layers, Patterns: patterns}
	}

	if hasAny(set, "app", "models", "views") || hasAny(set, "src/controllers", "src/models") {
		style = "mvc"
		confidence = "medium"
		layers = append(layers, "model", "view", "controller")
		patterns = append(patterns, "mvc")
		return profile.ArchitectureDetection{Style: style, Confidence: confidence, Layers: layers, Patterns: patterns}
	}

	if hasAny(set, "cmd", "internal") {
		style = "go-service-layout"
		confidence = "medium"
		layers = append(layers, "cmd", "internal")
		patterns = append(patterns, "idiomatic-go")
	}

	return profile.ArchitectureDetection{Style: style, Confidence: confidence, Layers: layers, Patterns: patterns}
}

func hasAny(set map[string]bool, values ...string) bool {
	for _, v := range values {
		if set[strings.ToLower(v)] {
			return true
		}
	}
	return false
}
