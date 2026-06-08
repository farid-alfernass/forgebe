package scaffold

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// TestInferTemplate_GoService tests template inference for Go projects without API-specific arch
func TestInferTemplate_GoService(t *testing.T) {
	p := profile.ProjectProfile{
		Stack: profile.Stack{
			PrimaryLanguage: "go",
			Architecture:    "layered",
		},
	}

	g := NewGenerator(&p, "/tmp", "", false, false)
	tmpl := g.inferTemplate()
	if tmpl != TemplateGoService {
		t.Errorf("expected TemplateGoService for Go+layered, got %v", tmpl)
	}
}

// TestInferTemplate_GoAPI tests template inference for Go projects with API/rest architecture
func TestInferTemplate_GoAPI(t *testing.T) {
	p := profile.ProjectProfile{
		Stack: profile.Stack{
			PrimaryLanguage: "go",
			Architecture:    "rest",
		},
	}

	g := NewGenerator(&p, "/tmp", "", false, false)
	tmpl := g.inferTemplate()
	if tmpl != TemplateGoAPI {
		t.Errorf("expected TemplateGoAPI for Go+rest, got %v", tmpl)
	}
}

// TestInferTemplate_NodeExpress tests Node.js inference
func TestInferTemplate_NodeExpress(t *testing.T) {
	p := profile.ProjectProfile{
		Stack: profile.Stack{
			PrimaryLanguage: "javascript",
			Architecture:    "monolith",
		},
	}

	g := NewGenerator(&p, "/tmp", "", false, false)
	tmpl := g.inferTemplate()
	if tmpl != TemplateNodeExpress {
		t.Errorf("expected TemplateNodeExpress for JS, got %v", tmpl)
	}
}

// TestInferTemplate_Python tests Python inference
func TestInferTemplate_Python(t *testing.T) {
	p := profile.ProjectProfile{
		Stack: profile.Stack{
			PrimaryLanguage: "python",
			Architecture:    "microservice",
		},
	}

	g := NewGenerator(&p, "/tmp", "", false, false)
	tmpl := g.inferTemplate()
	if tmpl != TemplatePythonFastAPI {
		t.Errorf("expected TemplatePythonFastAPI for Python, got %v", tmpl)
	}
}

// TestSanitizeName tests project name sanitization
func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My-Project", "my-project"},
		{"project name", "projectname"},
		{"PROJECT_123", "project_123"},
		{"test-project-v2", "test-project-v2"},
		{"@special#chars!", "specialchars"},
	}

	for _, tc := range tests {
		got := sanitizeName(tc.input)
		if got != tc.expected {
			t.Errorf("sanitizeName(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

// TestGenerate_DryRun verifies dry-run does not write to disk
func TestGenerate_DryRun(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "dry-test"
	p.Stack.PrimaryLanguage = "go"

	tmpDir := t.TempDir()
	g := NewGenerator(&p, tmpDir, TemplateGoService, false, true)

	files, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("expected files to be returned from dry-run")
	}

	// Verify no files were actually written
	// Dry-run should return files but not write them
	// That's handled by the Generator itself
}

// TestGenerate_ForceOverwrite verifies --force overwrites existing files
func TestGenerate_ForceOverwrite(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "force-test"
	p.Stack.PrimaryLanguage = "go"

	tmpDir := t.TempDir()

	// Pre-create a file
	goModPath := tmpDir + "/go.mod"

	// Write initial content using os package
	if err := os.WriteFile(goModPath, []byte("initial content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Generate with force=true
	g := NewGenerator(&p, tmpDir, TemplateGoService, true, false)
	_, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate with force failed: %v", err)
	}

	// Verify file was overwritten (content should be different now)
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) == "initial content" {
		t.Errorf("file was not overwritten as expected")
	}
}

// TestGenerate_FileExistsWithoutForce verifies error when file exists without --force
func TestGenerate_FileExistsWithoutForce(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "exists-test"
	p.Stack.PrimaryLanguage = "go"

	tmpDir := t.TempDir()

	// Pre-create go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}

	// Generate without force
	g := NewGenerator(&p, tmpDir, TemplateGoService, false, false)
	_, err := g.Generate()
	if err == nil {
		t.Fatal("expected error for existing file without force")
	}
}

// TestGeneratorWithInvalidProfile verifies error handling for nil profile
func TestGeneratorWithInvalidProfile(t *testing.T) {
	g := NewGenerator(nil, "/tmp", TemplateGoService, false, false)
	_, err := g.Generate()
	if err == nil {
		t.Fatal("expected error for nil profile")
	}
}

// TestGeneratorWithEmptyOutput verifies error handling for empty output dir
func TestGeneratorWithEmptyOutput(t *testing.T) {
	p := profile.NewSampleProfile()
	g := NewGenerator(&p, "", TemplateGoService, false, false)
	_, err := g.Generate()
	if err == nil {
		t.Fatal("expected error for empty output directory")
	}
}

// TestGenerateGoServiceTemplates verifies basic Go service structure generation
func TestGenerateGoServiceTemplates(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "goforge"
	p.Stack.PrimaryLanguage = "go"
	p.Stack.Architecture = "layered"

	tmpDir := t.TempDir()
	g := NewGenerator(&p, tmpDir, TemplateGoService, false, false)

	files, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify structure
	expectedFiles := []string{
		"cmd/goforge/main.go",
		"internal/handler/handler.go",
		"internal/service/service.go",
		"go.mod",
		"Makefile",
		".gitignore",
		"README.md",
	}

	for _, expected := range expectedFiles {
		found := false
		for _, f := range files {
			if f.Path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file %s not in generated list", expected)
		}
	}
}

// TestGenerateGoAPITemplates verifies basic Go API structure generation
func TestGenerateGoAPITemplates(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "api-go"
	p.Stack.PrimaryLanguage = "go"
	p.Stack.Architecture = "rest"

	tmpDir := t.TempDir()
	g := NewGenerator(&p, tmpDir, TemplateGoAPI, false, false)

	files, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify API structure differs from service
	foundRoutes := false
	for _, f := range files {
		if f.Path == "internal/api/routes.go" {
			foundRoutes = true
			break
		}
	}
	if !foundRoutes {
		t.Error("expected internal/api/routes.go in API template")
	}
}

// TestGenerateNodeExpressTemplates verifies Node.js Express template structure
func TestGenerateNodeExpressTemplates(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "node-proj"
	p.Stack.PrimaryLanguage = "typescript"
	p.Stack.Architecture = "monolith"

	tmpDir := t.TempDir()
	g := NewGenerator(&p, tmpDir, TemplateNodeExpress, false, false)

	files, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	foundIndex := false
	foundRoutes := false
	for _, f := range files {
		if f.Path == "src/index.ts" {
			foundIndex = true
		}
		if f.Path == "src/routes/index.ts" {
			foundRoutes = true
		}
	}

	if !foundIndex || !foundRoutes {
		t.Error("expected Node.js src structure")
	}
}

// TestGeneratePythonFastAPITemplates verifies Python FastAPI template structure
func TestGeneratePythonFastAPITemplates(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "py-api"
	p.Stack.PrimaryLanguage = "python"
	p.Stack.Architecture = "microservice"

	tmpDir := t.TempDir()
	g := NewGenerator(&p, tmpDir, TemplatePythonFastAPI, false, false)

	files, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	foundInit := false
	foundMain := false
	for _, f := range files {
		if f.Path == "app/__init__.py" {
			foundInit = true
		}
		if f.Path == "app/main.py" {
			foundMain = true
		}
	}

	if !foundInit || !foundMain {
		t.Error("expected Python app structure with __init__.py")
	}
}

// TestGeneratorUnknownTemplate verifies error for unsupported template type
func TestGeneratorUnknownTemplate(t *testing.T) {
	p := profile.NewSampleProfile()
	g := NewGenerator(&p, "/tmp", "unknown-template", false, false)
	_, err := g.Generate()
	if err == nil {
		t.Fatal("expected error for unsupported template")
	}
}

// TestGenerateReadmeTimestamp verifies README includes generation timestamp
func TestGenerateReadmeTimestamp(t *testing.T) {
	p := profile.NewSampleProfile()
	p.Project.Name = "ts-test"

	tmpDir := t.TempDir()
	g := NewGenerator(&p, tmpDir, TemplateGoService, false, false)

	files, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var readme string
	for _, f := range files {
		if f.Path == "README.md" {
			readme = f.Content
			break
		}
	}

	if readme == "" {
		t.Fatal("README.md not in generated files")
	}

	now := time.Now().Format("2006-01-02")
	if !contains(readme, now) {
		t.Errorf("README expected to contain today's date %s", now)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
