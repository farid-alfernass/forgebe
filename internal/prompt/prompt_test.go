package prompt

import (
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestNewGenerator_NilProfile(t *testing.T) {
	_, err := NewGenerator(nil)
	if err == nil {
		t.Fatal("expected error for nil profile")
	}
}

func TestNewGenerator_ValidProfile(t *testing.T) {
	p := profile.NewSampleProfile()
	g, err := NewGenerator(&p)
	if err != nil {
		t.Fatalf("NewGenerator failed: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil generator")
	}
}

func TestGenerate_ImplementMode(t *testing.T) {
	p := profile.NewSampleProfile()
	g, _ := NewGenerator(&p)

	output, err := g.Generate("implement", "Add a new REST endpoint for user CRUD")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if output == "" {
		t.Fatal("expected non-empty output")
	}
	if !contains(output, "implement") && !contains(output, "Implement") {
		t.Error("expected output to reference implementation mode")
	}
	if !contains(output, "go") {
		t.Error("expected output to contain language from profile")
	}
}

func TestGenerate_ReviewMode(t *testing.T) {
	p := profile.NewSampleProfile()
	g, _ := NewGenerator(&p)

	output, err := g.Generate("review", "Review the authentication middleware")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if !contains(output, "review") && !contains(output, "Review") {
		t.Error("expected output to reference review mode")
	}
}

func TestGenerate_DebugMode(t *testing.T) {
	p := profile.NewSampleProfile()
	g, _ := NewGenerator(&p)

	output, err := g.Generate("debug", "Fix the panic in request handler")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if !contains(output, "debug") && !contains(output, "Debug") {
		t.Error("expected output to reference debug mode")
	}
}

func TestGenerate_PlanMode(t *testing.T) {
	p := profile.NewSampleProfile()
	g, _ := NewGenerator(&p)

	output, err := g.Generate("plan", "Migrate from REST to gRPC")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if !contains(output, "plan") && !contains(output, "Plan") {
		t.Error("expected output to reference plan mode")
	}
}

func TestGenerate_UnknownMode(t *testing.T) {
	p := profile.NewSampleProfile()
	g, _ := NewGenerator(&p)

	_, err := g.Generate("unknown-mode", "some task")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestGenerate_IncludesProjectContext(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "my-backend"
	p.Stack.PrimaryLanguage = "go"
	p.Stack.Architecture = "clean/hexagonal"
	p.Policy.Testing.Strategy = "tdd"
	g, _ := NewGenerator(&p)

	output, _ := g.Generate("implement", "Add health endpoint")
	if !contains(output, "my-backend") {
		t.Error("expected output to contain project name")
	}
	if !contains(output, "clean/hexagonal") {
		t.Error("expected output to contain architecture")
	}
	if !contains(output, "tdd") {
		t.Error("expected output to contain testing strategy")
	}
}

func TestGenerate_IncludesSensitiveAreas(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Areas.SensitiveAreas = []string{"Authentication", "Payment"}
	g, _ := NewGenerator(&p)

	output, _ := g.Generate("implement", "Modify auth flow")
	if !contains(output, "Authentication") {
		t.Error("expected output to include sensitive areas")
	}
}

func TestGenerate_EmptyTask(t *testing.T) {
	p := profile.NewSampleProfile()
	g, _ := NewGenerator(&p)

	output, err := g.Generate("implement", "")
	if err != nil {
		t.Fatalf("Generate should not error on empty task: %v", err)
	}
	if output == "" {
		t.Fatal("expected non-empty output even without task")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
