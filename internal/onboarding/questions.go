package onboarding

import (
	"github.com/AlecAivazis/survey/v2"
)

// Questions defines all the questions for the guided onboarding.
var Questions = []*survey.Question{
	{
		Name: "projectName",
		Prompt: &survey.Input{
			Message: "What is the name of this project?",
			Help:    "Used for display purposes only.",
		},
		Validate: survey.Required,
	},
	{
		Name: "projectMaturity",
		Prompt: &survey.Select{
			Message: "Is this a new or existing project?",
			Options: []string{"new project", "existing codebase", "legacy codebase", "refactor in progress"},
			Default: "existing codebase",
			Help:    "Helps ForgeBE set appropriate defaults.",
		},
	},
	{
		Name: "projectType",
		Prompt: &survey.Select{
			Message: "What type of project is this?",
			Options: []string{"monolith", "service", "library", "other"},
			Default: "service",
			Help:    "Service = HTTP API, Library = reusable package, Monolith = all-in-one.",
		},
	},
	{
		Name: "primaryLanguage",
		Prompt: &survey.Input{
			Message: "What is the primary programming language?",
			Help:    "e.g., go, javascript, typescript, python, java, rust",
		},
		Validate: survey.Required,
	},
	{
		Name: "framework",
		Prompt: &survey.Input{
			Message: "What main framework or runtime do you use? (optional)",
			Help:    "e.g., gin, express, react, spring, django, etc.",
		},
	},
	{
		Name: "packageManager",
		Prompt: &survey.Input{
			Message: "What package manager or build tool do you use? (optional)",
			Help:    "e.g., go modules, npm, yarn, pip, maven, gradle, cargo",
		},
	},
	{
		Name: "architecture",
		Prompt: &survey.Select{
			Message: "What architecture style does this project follow?",
			Options: []string{"clean/hexagonal", "layered (service-repository)", "mvc", "event-driven", "microservices", "other / not sure"},
			Default: "layered (service-repository)",
			Help:    "Helps ForgeBE enforce boundaries.",
		},
	},
	{
		Name: "testingPolicy",
		Prompt: &survey.Select{
			Message: "What is the testing policy for new logic?",
			Options: []string{"TDD (tests required before implementation)", "Test after implementation", "No strict policy yet", "Tests only for critical paths"},
			Default: "TDD (tests required before implementation)",
			Help:    "Determines whether ForgeBE will suggest writing tests first.",
		},
	},
	{
		Name: "dependencyPolicy",
		Prompt: &survey.Select{
			Message: "Can ForgeBE suggest adding new dependencies?",
			Options: []string{"Only with explicit approval", "Yes, for small utility packages", "Yes, freely", "No, never add dependencies"},
			Default: "Only with explicit approval",
			Help:    "Controls how ForgeBE interacts with dependency files.",
		},
	},
	{
		Name: "changePolicy",
		Prompt: &survey.Select{
			Message: "How aggressive can ForgeBE be when modifying code?",
			Options: []string{"Preserve structure, minimal changes", "Allow refactoring within same layer", "Allow cross-file refactoring", "Allow broad refactor with approval"},
			Default: "Preserve structure, minimal changes",
			Help:    "Limits the scope of suggested changes.",
		},
	},
	{
		Name: "qualityPriority",
		Prompt: &survey.Select{
			Message: "What is your top priority for quality?",
			Options: []string{"Safety / minimal regression", "Code quality / maintainability", "Delivery speed", "Performance"},
			Default: "Safety / minimal regression",
			Help:    "Influences ForgeBE's strictness.",
		},
	},
	{
		Name: "deliveryMode",
		Prompt: &survey.Select{
			Message: "How should ForgeBE work with you?",
			Options: []string{"Plan-first (propose plan, then implement)", "Direct execution (implement immediately if confident)", "Hybrid (plan for complex, direct for simple)"},
			Default: "Plan-first (propose plan, then implement)",
			Help:    "Determines ForgeBE's interaction style.",
		},
	},
	{
		Name: "aiModelStrategy",
		Prompt: &survey.Select{
			Message: "How do you plan to use AI models?",
			Options: []string{"Single model (e.g., only Claude)", "Multi-model (different models for different tasks)", "Round-robin (rotate models)", "Best available model"},
			Default: "Multi-model (different models for different tasks)",
			Help:    "ForgeBE can adapt its suggestions based on model strengths.",
		},
	},
	{
		Name: "sensitiveAreas",
		Prompt: &survey.MultiSelect{
			Message: "Which areas are considered sensitive? (choose all that apply)",
			Options: []string{"Authentication / Authorization", "Payments / Billing", "Database migrations", "Public API contracts", "Infrastructure / DevOps", "Background jobs / Workers", "Other (specify)"},
			Default: []string{"Authentication / Authorization", "Payments / Billing"},
		},
	},
	{
		Name: "otherSensitive",
		Prompt: &survey.Input{
			Message: "If you selected 'Other (specify)' above, please describe:",
			Help:    "Optional.",
		},
	},
}
