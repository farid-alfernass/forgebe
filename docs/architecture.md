# ForgeBE Architecture

## Overview

ForgeBE is a Go CLI tool that helps backend engineers maintain consistent AI-assisted delivery. It stores project context locally and exports tool-specific adapters for various AI assistants.

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLI Layer                               │
│  init │ scan │ doctor │ profile │ brief │ export │ import      │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                      Core Domain                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │  Discovery  │  │   Profile   │  │       Adapters          │  │
│  │  (scanner)  │  │   (model)   │  │  (claude/cursor/...)    │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │  Contract   │  │   Export    │  │      Onboarding         │  │
│  │ (extractor) │  │(summary/zip)│  │   (guided questions)    │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                     Storage Layer                               │
│  ~/.forgebe/                                                    │
│    ├── projects/<id>/       # Per-project state                 │
│    │     ├── project-profile.yaml                               │
│    │     └── engineering-contract.md                            │
│    ├── exports/             # Generated adapters & bundles      │
│    └── cache/               # Temporary scan cache              │
└─────────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
forgebe/
├── cmd/
│   └── forgebe/
│       └── main.go           # CLI entrypoint
├── internal/
│   ├── cli/                  # Cobra command definitions
│   │   ├── root.go
│   │   ├── init.go
│   │   ├── scan.go
│   │   ├── doctor.go
│   │   ├── profile.go
│   │   ├── brief.go
│   │   ├── export.go
│   │   ├── import.go
│   │   └── output.go         # JSON/file output helpers
│   ├── discovery/            # Codebase scanning
│   │   ├── scanner.go
│   │   ├── detector_language.go
│   │   ├── detector_architecture.go
│   │   ├── detector_testing.go
│   │   ├── detector_dependencies.go
│   │   └── detector_ci.go
│   ├── profile/              # Project profile model & store
│   │   ├── model.go
│   │   ├── store.go
│   │   ├── validator.go
│   │   └── project_id.go
│   ├── contract/             # Engineering contract
│   │   ├── renderer.go
│   │   └── extractor.go
│   ├── adapters/             # AI tool adapters
│   │   └── adapter.go
│   ├── export/               # Export utilities
│   │   ├── summary.go
│   │   └── bundle.go
│   ├── onboarding/           # Guided init flow
│   │   ├── questions.go
│   │   ├── answers.go
│   │   ├── guided.go
│   │   └── defaults.go
│   └── storage/              # Local storage primitives
│       ├── paths.go
│       └── atomic.go
├── testdata/                 # Test fixtures (excluded from scan)
├── docs/                     # Documentation
├── .github/workflows/        # CI/CD
├── Makefile
├── README.md
├── CONTRIBUTING.md
├── CHANGELOG.md
├── LICENSE
└── SECURITY.md
```

## Data Flow

### 1. Initialization (`forgebe init`)

```
User runs `forgebe init`
        │
        ▼
┌───────────────┐     ┌─────────────────┐     ┌────────────────┐
│   Discovery   │────▶│   Onboarding    │────▶│  Profile Store │
│   (scanner)   │     │  (questions)    │     │   (YAML save)  │
└───────────────┘     └─────────────────┘     └────────────────┘
        │                     │                       │
        ▼                     ▼                       ▼
  Detects lang,         Fills gaps,            Writes to
  arch, testing,        user confirms          ~/.forgebe/projects/<id>/
  CI, deps
```

### 2. Export Adapter (`forgebe export adapter claude`)

```
User runs `forgebe export adapter claude <id>`
        │
        ▼
┌───────────────┐     ┌─────────────────┐     ┌────────────────┐
│ Profile Store │────▶│    Adapters     │────▶│  File Writer   │
│  (YAML load)  │     │  (templates)    │     │  (atomic)      │
└───────────────┘     └─────────────────┘     └────────────────┘
        │                     │                       │
        ▼                     ▼                       ▼
  Loads profile,        Renders tool-           Writes to
  validates             specific markdown       ~/.forgebe/exports/
```

### 3. Bundle Export/Import

```
Export:
  Project Dir ──▶ ZIP Writer ──▶ .forgebe.zip

Import:
  .forgebe.zip ──▶ ZIP Reader ──▶ Target Dir
                       │
                       ▼
              Path validation
              (ZipSlip + Symlink)
```

## Key Design Decisions

### Local-First Storage

All ForgeBE state lives under `~/.forgebe/`:
- **Isolation**: Each project has its own directory keyed by a deterministic ID.
- **Portability**: Bundles can be exported and shared.
- **Git-safe**: No writes to the target repository by default.

### Deterministic Project ID

```go
ProjectID(repoPath, repoRemote) = hash(absPath + remote)[:8]
```

This ensures the same project always maps to the same ID across sessions.

### Discovery vs Onboarding

- **Discovery** (automatic): Reads file structure, manifests, and CI configs.
- **Onboarding** (interactive): Asks clarifying questions for ambiguous fields.
- **Defaults**: Discovery results pre-fill onboarding answers.

### Adapter Templates

Each adapter (Claude, Cursor, Copilot, Hermes, Generic) has a Go text/template that renders the profile into tool-specific instructions:

```go
adapters.Render("claude", profile) → (string, filename, error)
```

### Security Boundaries

- **Bundle import**: Path validation, symlink resolution, no overwrites without explicit targeting.
- **Atomic writes**: Temp file + rename to prevent partial writes.
- **No shell execution**: All operations are pure Go.

## Extension Points

### Adding a New Adapter

1. Add template to `internal/adapters/adapter.go`:
   ```go
   adapterRegistry["newtool"] = adapterSpec{
       Filename: "NEWTOOL_CONTEXT.md",
       Template: `...`,
   }
   ```

2. Test: `forgebe export adapter newtool <id>`

### Adding a New Detector

1. Create `internal/discovery/detector_xxx.go`.
2. Implement detection function: `DetectXxx(root string, files []string) XxxDetection`.
3. Wire into `scanner.go` `Scan()` method.
4. Add tests in `detector_xxx_test.go`.

### Adding a New CLI Command

1. Create `internal/cli/newcmd.go`.
2. Implement `newNewcmdCmd() *cobra.Command`.
3. Register in `root.go`: `cmd.AddCommand(newNewcmdCmd())`.
4. Add `--json` support via `OutputJSON(cmd)` if applicable.
