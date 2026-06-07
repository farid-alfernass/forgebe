package adapters

import (
	"strings"
	"testing"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

func TestRenderAdapter_Claude(t *testing.T) {
	p := profile.NewSampleProfile()
	output, filename, err := Render("claude", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filename != "CLAUDE.md" {
		t.Fatalf("expected CLAUDE.md, got %s", filename)
	}
	if !strings.Contains(output, "Primary language") {
		t.Fatalf("expected claude output to contain profile summary, got: %s", output)
	}
}

func TestRenderAdapter_Cursor(t *testing.T) {
	p := profile.NewSampleProfile()
	_, filename, err := Render("cursor", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filename != ".cursorrules" {
		t.Fatalf("expected .cursorrules, got %s", filename)
	}
}

func TestRenderAdapter_UnknownTool(t *testing.T) {
	p := profile.NewSampleProfile()
	_, _, err := Render("unknown-tool", p)
	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestRenderBrief_Generic(t *testing.T) {
	p := profile.NewSampleProfile()
	brief, err := RenderBrief("generic", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(brief, "Project:") {
		t.Fatalf("expected compact brief output, got %s", brief)
	}
}
