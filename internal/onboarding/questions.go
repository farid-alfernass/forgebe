package onboarding

import (
	"github.com/AlecAivazis/survey/v2"
)

func inputQuestion(name, message, help string, required bool) *survey.Question {
	q := &survey.Question{
		Name: name,
		Prompt: &survey.Input{
			Message: message,
			Help:    help,
		},
	}
	if required {
		q.Validate = survey.Required
	}
	return q
}

func selectQuestion(name, message string, options []string, defaultValue, help string) *survey.Question {
	return &survey.Question{
		Name: name,
		Prompt: &survey.Select{
			Message: message,
			Options: options,
			Default: defaultValue,
			Help:    help,
		},
	}
}

func multiSelectQuestion(name, message string, options, defaults []string) *survey.Question {
	return &survey.Question{
		Name: name,
		Prompt: &survey.MultiSelect{
			Message: message,
			Options: options,
			Default: defaults,
		},
	}
}

// Questions defines all the questions for the guided onboarding.
var Questions = []*survey.Question{
	inputQuestion("projectName", "What is the name of this project?", "Used for display purposes only.", true),
	selectQuestion(
		"projectMaturity",
		"Is this a new or existing project?",
		[]string{"new project", "existing codebase", "legacy codebase", "refactor in progress"},
		"existing codebase",
		"Helps ForgeBE set appropriate defaults.",
	),
	selectQuestion(
		"projectType",
		"What type of project is this?",
		[]string{"monolith", "service", "library", "other"},
		"service",
		"Service = HTTP API, Library = reusable package, Monolith = all-in-one.",
	),
	inputQuestion("primaryLanguage", "What is the primary programming language?", "e.g., go, javascript, typescript, python, java, rust", true),
	inputQuestion("framework", "What main framework or runtime do you use? (optional)", "e.g., gin, express, react, spring, django, etc.", false),
	inputQuestion("packageManager", "What package manager or build tool do you use? (optional)", "e.g., go modules, npm, yarn, pip, maven, gradle, cargo", false),
	selectQuestion(
		"architecture",
		"What architecture style does this project follow?",
		[]string{"clean/hexagonal", "layered (service-repository)", "mvc", "event-driven", "microservices", "other / not sure"},
		"layered (service-repository)",
		"Helps ForgeBE enforce boundaries.",
	),
	selectQuestion(
		"testingPolicy",
		"What is the testing policy for new logic?",
		[]string{"TDD (tests required before implementation)", "Test after implementation", "No strict policy yet", "Tests only for critical paths"},
		"TDD (tests required before implementation)",
		"Determines whether ForgeBE will suggest writing tests first.",
	),
	selectQuestion(
		"dependencyPolicy",
		"Can ForgeBE suggest adding new dependencies?",
		[]string{"Only with explicit approval", "Yes, for small utility packages", "Yes, freely", "No, never add dependencies"},
		"Only with explicit approval",
		"Controls how ForgeBE interacts with dependency files.",
	),
	selectQuestion(
		"changePolicy",
		"How aggressive can ForgeBE be when modifying code?",
		[]string{"Preserve structure, minimal changes", "Allow refactoring within same layer", "Allow cross-file refactoring", "Allow broad refactor with approval"},
		"Preserve structure, minimal changes",
		"Limits the scope of suggested changes.",
	),
	selectQuestion(
		"qualityPriority",
		"What is your top priority for quality?",
		[]string{"Safety / minimal regression", "Code quality / maintainability", "Delivery speed", "Performance"},
		"Safety / minimal regression",
		"Influences ForgeBE's strictness.",
	),
	selectQuestion(
		"deliveryMode",
		"How should ForgeBE work with you?",
		[]string{"Plan-first (propose plan, then implement)", "Direct execution (implement immediately if confident)", "Hybrid (plan for complex, direct for simple)"},
		"Plan-first (propose plan, then implement)",
		"Determines ForgeBE's interaction style.",
	),
	selectQuestion(
		"aiModelStrategy",
		"How do you plan to use AI models?",
		[]string{"Single model (e.g., only Claude)", "Multi-model (different models for different tasks)", "Round-robin (rotate models)", "Best available model"},
		"Multi-model (different models for different tasks)",
		"ForgeBE can adapt its suggestions based on model strengths.",
	),
	multiSelectQuestion(
		"sensitiveAreas",
		"Which areas are considered sensitive? (choose all that apply)",
		[]string{"Authentication / Authorization", "Payments / Billing", "Database migrations", "Public API contracts", "Infrastructure / DevOps", "Background jobs / Workers", "Other (specify)"},
		[]string{"Authentication / Authorization", "Payments / Billing"},
	),
	inputQuestion("otherSensitive", "If you selected 'Other (specify)' above, please describe:", "Optional.", false),
}
