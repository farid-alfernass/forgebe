package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faridtriwicaksono/forgebe/internal/profile"
)

// TemplateType represents a project template type.
type TemplateType string

const (
	TemplateGoService     TemplateType = "go-service"
	TemplateGoAPI         TemplateType = "go-api"
	TemplateNodeExpress   TemplateType = "node-express"
	TemplateNodeNestJS    TemplateType = "node-nestjs"
	TemplatePythonFastAPI TemplateType = "python-fastapi"
	TemplateJavaSpring    TemplateType = "java-spring"
)

// File represents a file to be generated.
type File struct {
	Path    string
	Content string
	IsDir   bool
}

// Generator generates project scaffolding based on profile.
type Generator struct {
	profile      *profile.ProjectProfile
	outputDir    string
	templateType TemplateType
	force        bool
	dryRun       bool
}

// NewGenerator creates a new scaffold generator.
func NewGenerator(p *profile.ProjectProfile, outputDir string, templateType TemplateType, force, dryRun bool) *Generator {
	return &Generator{
		profile:      p,
		outputDir:    outputDir,
		templateType: templateType,
		force:        force,
		dryRun:       dryRun,
	}
}

// Generate creates project scaffolding.
func (g *Generator) Generate() ([]File, error) {
	if g.profile == nil {
		return nil, fmt.Errorf("scaffold: profile is nil")
	}

	if g.outputDir == "" {
		return nil, fmt.Errorf("scaffold: output directory is empty")
	}

	// Resolve template type if not specified
	tmpl := g.templateType
	if tmpl == "" {
		tmpl = g.inferTemplate()
	}

	// Generate files based on template
	var files []File
	switch tmpl {
	case TemplateGoService:
		files = g.generateGoService()
	case TemplateGoAPI:
		files = g.generateGoAPI()
	case TemplateNodeExpress:
		files = g.generateNodeExpress()
	case TemplatePythonFastAPI:
		files = g.generatePythonFastAPI()
	default:
		return nil, fmt.Errorf("scaffold: unsupported template type: %v", tmpl)
	}

	// Write files if not dry-run
	if !g.dryRun {
		if err := g.writeFiles(files); err != nil {
			return nil, err
		}
	}

	return files, nil
}

// inferTemplate infers the template from the profile.
func (g *Generator) inferTemplate() TemplateType {
	lang := g.profile.Stack.PrimaryLanguage
	arch := g.profile.Stack.Architecture

	switch lang {
	case "go":
		if strings.Contains(arch, "api") || strings.Contains(arch, "rest") {
			return TemplateGoAPI
		}
		return TemplateGoService
	case "typescript", "javascript":
		if strings.Contains(arch, "nest") {
			return TemplateNodeNestJS
		}
		return TemplateNodeExpress
	case "python":
		return TemplatePythonFastAPI
	case "java":
		return TemplateJavaSpring
	default:
		return TemplateGoService
	}
}

// generateGoService generates a Go service template.
func (g *Generator) generateGoService() []File {
	projectName := sanitizeName(g.profile.Project.Name)
	return []File{
		{Path: "cmd/" + projectName, IsDir: true},
		{Path: "cmd/" + projectName + "/main.go", Content: g.goServiceMain(projectName)},
		{Path: "internal/handler", IsDir: true},
		{Path: "internal/handler/handler.go", Content: g.goServiceHandler()},
		{Path: "internal/service", IsDir: true},
		{Path: "internal/service/service.go", Content: g.goServiceService()},
		{Path: "go.mod", Content: g.goMod(projectName)},
		{Path: "Makefile", Content: g.goMakefile(projectName)},
		{Path: ".gitignore", Content: goGitignore},
		{Path: "README.md", Content: g.readme(projectName)},
	}
}

// generateGoAPI generates a Go API template.
func (g *Generator) generateGoAPI() []File {
	projectName := sanitizeName(g.profile.Project.Name)
	return []File{
		{Path: "cmd/api", IsDir: true},
		{Path: "cmd/api/main.go", Content: g.goAPIMain(projectName)},
		{Path: "internal/api", IsDir: true},
		{Path: "internal/api/routes.go", Content: g.goAPIRoutes()},
		{Path: "internal/models", IsDir: true},
		{Path: "internal/models/models.go", Content: goAPIModels},
		{Path: "internal/repository", IsDir: true},
		{Path: "internal/repository/repository.go", Content: goAPIRepository},
		{Path: "go.mod", Content: g.goMod(projectName)},
		{Path: "Makefile", Content: g.goMakefile(projectName)},
		{Path: ".gitignore", Content: goGitignore},
		{Path: "README.md", Content: g.readme(projectName)},
	}
}

// generateNodeExpress generates a Node.js Express template.
func (g *Generator) generateNodeExpress() []File {
	projectName := sanitizeName(g.profile.Project.Name)
	return []File{
		{Path: "src", IsDir: true},
		{Path: "src/index.ts", Content: g.nodeExpressIndex()},
		{Path: "src/routes", IsDir: true},
		{Path: "src/routes/index.ts", Content: nodeExpressRoutes},
		{Path: "src/middleware", IsDir: true},
		{Path: "src/middleware/errorHandler.ts", Content: nodeExpressErrorHandler},
		{Path: "package.json", Content: g.nodePackageJson(projectName)},
		{Path: "tsconfig.json", Content: nodeTsconfig},
		{Path: ".gitignore", Content: nodeGitignore},
		{Path: "README.md", Content: g.readme(projectName)},
	}
}

// generatePythonFastAPI generates a Python FastAPI template.
func (g *Generator) generatePythonFastAPI() []File {
	projectName := sanitizeName(g.profile.Project.Name)
	return []File{
		{Path: "app", IsDir: true},
		{Path: "app/__init__.py", Content: ""},
		{Path: "app/main.py", Content: g.pythonFastAPIMain()},
		{Path: "app/api", IsDir: true},
		{Path: "app/api/__init__.py", Content: ""},
		{Path: "app/api/routes.py", Content: pythonFastAPIRoutes},
		{Path: "app/models", IsDir: true},
		{Path: "app/models/__init__.py", Content: ""},
		{Path: "app/models/schemas.py", Content: pythonFastAPISchemas},
		{Path: "requirements.txt", Content: pythonRequirements},
		{Path: ".gitignore", Content: pythonGitignore},
		{Path: "README.md", Content: g.readme(projectName)},
	}
}

// writeFiles writes generated files to disk.
func (g *Generator) writeFiles(files []File) error {
	for _, f := range files {
		path := filepath.Join(g.outputDir, f.Path)

		if f.IsDir {
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("scaffold: mkdir %s: %w", path, err)
			}
			continue
		}

		// Check if file exists and force is not set
		if _, err := os.Stat(path); err == nil && !g.force {
			return fmt.Errorf("scaffold: file already exists: %s (use --force to overwrite)", path)
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("scaffold: mkdir parent: %w", err)
		}

		// Write file
		if err := os.WriteFile(path, []byte(f.Content), 0644); err != nil {
			return fmt.Errorf("scaffold: write %s: %w", path, err)
		}
	}
	return nil
}

// Helper functions for template generation

func (g *Generator) goServiceMain(projectName string) string {
	return fmt.Sprintf(`package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"status\":\"ok\"}"))
	})

	log.Println("Starting %s on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
`, projectName)
}

func (g *Generator) goAPIMain(projectName string) string {
	return `package main

import (
	"log"
	"net/http"

	"internal/api"
)

func main() {
	router := http.NewServeMux()
	api.RegisterRoutes(router)

	log.Println("Starting API server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
`
}

func (g *Generator) goServiceHandler() string {
	return `package handler

import (
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\":\"ok\"}"))
}
`
}

func (g *Generator) goServiceService() string {
	return `package service

// Service provides business logic.
type Service struct {
	// Add dependencies here
}

// NewService creates a new service.
func NewService() *Service {
	return &Service{}
}
`
}

func (g *Generator) goAPIRoutes() string {
	return `package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"status\":\"ok\"}"))
	})
}
`
}

func (g *Generator) goMod(projectName string) string {
	return fmt.Sprintf(`module github.com/your-org/%s

go 1.22
`, projectName)
}

func (g *Generator) goMakefile(projectName string) string {
	return fmt.Sprintf(`.PHONY: build test run fmt vet

build:
	go build -o bin/%s ./cmd/%s

test:
	go test -v -cover ./...

run: build
	./bin/%s

fmt:
	go fmt ./...

vet:
	go vet ./...
`, projectName, projectName, projectName)
}

func (g *Generator) nodeExpressIndex() string {
	return `import express, { Express } from 'express';
import routes from './routes';

const app: Express = express();
const PORT = process.env.PORT || 8080;

app.use(express.json());
app.use('/api', routes);

app.listen(PORT, () => {
  console.log('Server running on port ' + PORT);
});
`
}

func (g *Generator) nodePackageJson(projectName string) string {
	return fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "main": "dist/index.js",
  "scripts": {
    "dev": "ts-node src/index.ts",
    "build": "tsc",
    "start": "node dist/index.js",
    "test": "jest"
  },
  "dependencies": {
    "express": "^4.18.0"
  },
  "devDependencies": {
    "@types/express": "^4.17.0",
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0",
    "ts-node": "^10.0.0"
  }
}
`, projectName)
}

func (g *Generator) pythonFastAPIMain() string {
	return `from fastapi import FastAPI
from app.api import routes

app = FastAPI(title="API", version="1.0.0")

@app.get("/health")
def health():
    return {"status": "ok"}

app.include_router(routes.router)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
`
}

func (g *Generator) readme(projectName string) string {
	return fmt.Sprintf(`# %s

Generated with ForgeBE on %s

## Getting Started

### Prerequisites
- Go 1.22+ (for Go projects)
- Node.js 18+ (for Node.js projects)
- Python 3.9+ (for Python projects)

### Build

make build

### Run

make run

### Test

make test

## Architecture

This project follows clean architecture principles with clear separation of concerns.

## License

MIT
`, projectName, time.Now().Format("2006-01-02"))
}

func sanitizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return -1
	}, name)
	return name
}

// Template constants

const goGitignore = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.so.*
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# Build output
bin/
dist/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`

const goAPIModels = `package models

type User struct {
	ID    string
	Name  string
	Email string
}
`

const goAPIRepository = `package repository

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}
`

const nodeExpressRoutes = `import { Router } from 'express';

const router = Router();

router.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

export default router;
`

const nodeExpressErrorHandler = `import { Request, Response, NextFunction } from 'express';

export const errorHandler = (
  err: Error,
  req: Request,
  res: Response,
  next: NextFunction
) => {
  console.error(err.stack);
  res.status(500).json({ error: 'Internal Server Error' });
};
`

const nodeTsconfig = `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules"]
}
`

const nodeGitignore = `node_modules/
dist/
*.log
.env
.env.local
.DS_Store
`

const pythonFastAPIRoutes = `from fastapi import APIRouter

router = APIRouter(prefix="/api", tags=["items"])

@router.get("/items")
def list_items():
    return {"items": []}

@router.get("/items/{item_id}")
def get_item(item_id: str):
    return {"item_id": item_id}
`

const pythonFastAPISchemas = `from pydantic import BaseModel

class Item(BaseModel):
    id: str
    name: str
    description: str = None

class User(BaseModel):
    id: str
    name: str
    email: str
`

const pythonRequirements = `fastapi==0.104.0
uvicorn==0.24.0
pydantic==2.5.0
python-dotenv==1.0.0
`

const pythonGitignore = `__pycache__/
*.py[cod]
*$py.class
*.so
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
.venv
env/
venv/
.env
.DS_Store
`
