# CHANGELOG

All notable changes to ForgeBE are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned

- Task-specific briefs (`forgebe task-brief`)
- Multi-project workspace support
- Configuration profiles for different AI workflows
- Integration with GitHub/GitLab APIs for automatic discovery

---

## [1.3.0] — 2026-06-08

### Added

#### Documentation Overhaul
- **Blueprint Guide** (`docs/blueprint-guide.md`) — Deep dive metodologi AI-Native Backend Engineering: TDD-AI Loop, Context Layer, Quality Gates, Architectural Framework, Project Lifecycle, Security, Extension Points
- **Real-World Examples** (`docs/real-world-examples.md`) — 5 skenario praktis: Go Payment Gateway, Node.js E-Commerce API, Java monolith adoption, multi-team setup, complete CI/CD pipeline
- **Release Guide** (`docs/release-guide.md`) — GitHub release process, versioning, binary distribution, automation
- **Enhanced Contributing Guide** (`CONTRIBUTING.md`) — Detailed guidelines: principles, setup, quality gates, testing, coding conventions, PR checklist, release process

#### Phase 10 Feature Expansion (Complete)
- `forgebe scaffold` — Project scaffolding with templates (go-service, go-api, node-express, node-nestjs, python-fastapi)
- `forgebe adopt` — Adopt existing projects with automatic discovery
- `forgebe verify` — Policy verification engine (8 policies: testing, dependency, architecture, sensitive areas, source roots, forbidden paths)
- `forgebe check` — Health check / quality gate for pre-commit and CI
- `forgebe prompt` — AI prompt generator (implement, review, debug, plan modes)

### Fixed
- Code duplication in scaffold generator reduced from 4.6% to 1.2%
- Security: pinned GitHub Actions to full SHA (`ci.yml`, `release.yml`)

### Changed
- README restructured with problem-solution approach
- Documentation reorganized for better navigation
- Architecture docs updated with Phase 10 components

### Removed
- Outdated planning docs (`docs/plans/`)
- Replaced `docs/phase10-feature-expansion.md` with blueprint-guide
- Removed `docs/brew-tap.md` and `docs/contract-extractor-roadmap.md`

### Quality Metrics
- **Coverage**: 81.5% (target >=80%) ✅
- **Duplicated Lines**: 1.2% (threshold <3%) ✅
- **Sonar Quality Gate**: OK ✅
- **All Tests**: PASS ✅

---

## [1.0.0] — 2026-06-07

### Added

#### Core Features
- `forgebe init` — Interactive and auto-discovery project initialization
- `forgebe scan` — Automatic codebase discovery (language, architecture, testing, CI, dependencies)
- `forgebe doctor` — Project profile health checks and validation
- `forgebe profile show` — Display canonical project profile
- `forgebe brief <tool>` — Compact AI-ready project context
- `forgebe export summary` — Full markdown project summary
- `forgebe export adapter <tool>` — Tool-specific adapter files (Claude, Cursor, Copilot, Hermes, Generic)
- `forgebe export bundle` — Portable .forgebe.zip bundles for team handoff
- `forgebe import` — Restore bundles with security validation
- Contract Extractor foundation — Extract engineering rules from markdown files
- `--json` flag — Machine-readable output for all major commands
- `--output` flag — Custom file destinations for exports

#### Documentation
- Comprehensive README with quickstart and feature overview
- CONTRIBUTING.md with development guidelines and quality gates
- docs/architecture.md — System design, data flow, extension points
- docs/security.md — Threat model, security boundaries, usage guidelines
- docs/cli-reference.md — Complete command reference and examples
- docs/contract-extractor-roadmap.md — Future roadmap for semantic extraction

#### Infrastructure
- GitHub Actions CI/CD with build matrix (Linux amd64, macOS amd64/arm64)
- Release workflow with checksum generation
- Local-first state management under `~/.forgebe/`
- Atomic writes to prevent corruption
- ZipSlip and symlink traversal protection

#### Quality & Security
- 100% code formatting (gofmt)
- Zero `go vet` warnings
- 7 core packages with comprehensive unit tests
- Deterministic project ID generation
- Guided onboarding with auto-discovery defaults
- Root-manifest prioritization in language detection
- Test fixture exclusion from production scanning

### Fixed

- Scanner now skips `testdata`, `fixtures`, `examples` directories
- Language detector prioritizes root manifests over nested ones
- Bundle import validates all extraction paths (ZipSlip protection)
- Bundle import resolves and validates symlinks
- Bundle export ensures parent directories exist
- Bundle export/import explicitly checks file close errors
- Fixture Go files now have valid package declarations

### Security

- Added symlink traversal hardening in bundle import
- Implemented ZipSlip path validation in archive extraction
- Added explicit close error handling for file I/O
- Secured all file writes with atomic operations
- No hardcoded secrets or credentials

---

## [0.5.0] — 2026-06-07

### Added

#### Phase 5: AI Adapters & Contract Extractor
- AI context adapters for Claude, Cursor, Copilot, Hermes, Generic
- `forgebe brief` command for compact AI prompts
- `forgebe export adapter` for tool-specific markdown files
- Contract Extractor foundation (`internal/contract/extractor.go`)
- Semantic extraction of testing strategy, dependencies, sensitive areas
- Roadmap documentation for future multi-file extraction

#### Phase 4: Doctor, Profile Show, Summary Export
- `forgebe doctor` command with profile validation
- `forgebe profile show` for readable profile display
- `forgebe export summary` for markdown summaries
- Profile validator with comprehensive checks
- Sample profile generator for testing

#### Phase 3: Guided Onboarding
- `forgebe init` with interactive questions
- `survey/v2` library integration for guided flow
- Auto-discovery defaults from `forgebe scan`
- Engineering contract markdown generation
- Onboarding answer collection and profile building

#### Phase 2: Codebase Discovery
- `forgebe scan` command for automatic detection
- Language detector (Go, Node, Python, Java, Rust, etc.)
- Architecture detector (layered, clean, hexagonal, Go-layout, MVC, DDD)
- Testing framework detector (go test, jest, pytest, vitest, etc.)
- CI platform detector (GitHub Actions, GitLab CI, Jenkins)
- Dependency detector (go.mod, package.json, requirements.txt, etc.)
- Discovery report with confidence scores and assumptions

#### Phase 1: Core Infrastructure
- Local-first storage at `~/.forgebe/`
- Deterministic project ID generation
- Profile storage layer (load/save YAML)
- Atomic write operations
- CLI skeleton with Cobra
- Version command
- GitHub Actions CI workflow
- Project files (Makefile, README, LICENSE, SECURITY.md, .gitignore)

### Security

- ZipSlip protection in bundle import (partial)
- Atomic file writes to prevent corruption
- Test fixture isolation in scanner

### Known Limitations

- Symlink traversal protection not yet complete (added in 1.0.0)
- Close error handling for bundle import incomplete (fixed in 1.0.0)
- Contract Extractor is foundation-only; semantic extraction deferred to v1.1+

---

## [0.1.0] — 2026-06-06

### Added

#### Phase 1: Initial Setup
- Go module initialization
- Core domain models (ProjectProfile, Stack, Policy, etc.)
- Storage layer (paths, atomic writes)
- CLI skeleton (Cobra/Viper integration)
- GitHub Actions CI workflow
- Project documentation (README, LICENSE, SECURITY.md)

### Status

- Proof of Concept (PoC): Models and storage infrastructure in place
- Ready for Phase 2 discovery implementation

---

## Migration Guide

### From v0.5.0 to v1.0.0

No breaking changes. All v0.5.0 bundles and profiles remain compatible.

**New in v1.0.0:**
- Use `--json` flag for machine-readable output
- Use `--output` flag for custom export paths
- Symlink traversal protection is now enabled by default in bundle import

---

## Release Schedule

| Version | Target Date | Focus |
|---------|-------------|-------|
| v0.1.0 | ✅ 2026-06-06 | PoC — core models & storage |
| v0.5.0 | ✅ 2026-06-07 | MVP — scan, init, export, adapters |
| v1.0.0 | ✅ 2026-06-07 | Production — security hardening, docs, v1 quality |
| v1.1.0 | 2026-07-XX | Semantic contract extraction from markdown |
| v2.0.0 | 2026-Q3 | Multi-project workspaces, team collaboration |

---

## Acknowledgments

ForgeBE was built to solve a real problem: backend engineers using AI tools need consistent, project-aware context without losing control or safety.

Thank you to all contributors and early testers who shaped this framework.
