# Phase 10 — Feature Expansion Guide

**Release:** v1.2.0  
**Date:** June 8, 2026  
**Status:** ✅ Complete (All 5 commands implemented & tested)

ForgeBE Phase 10 adds 5 new operational commands that make ForgeBE a complete backend engineering framework — not just a profile/context manager.

---

## Commands Overview

| Command | Purpose | Use Case |
|---------|---------|----------|
| `scaffold` | Generate boilerplate | Start new projects with consistent structure |
| `adopt` | Migrate existing projects | Onboard legacy/external projects to ForgeBE |
| `verify` | Validate profile policies | Check alignment with engineering standards |
| `check` | Health gate for CI/pre-commit | Fast validation for automation workflows |
| `prompt` | Context-aware AI prompts | Task-specific prompts for Claude/Cursor/Copilot |

---

## Phase 10.1: `forgebe scaffold`

Generate project boilerplate from profile or template override.

### Supported Templates
- `go-service` — Go service with layered architecture
- `go-api` — Go API with handlers, models, repository
- `node-express` — Node.js + TypeScript + Express
- `node-nestjs` — Node.js + TypeScript + NestJS
- `python-fastapi` — Python + FastAPI

### Examples

#### Generate using project profile (inferred language)
```bash
forgebe scaffold
# Output: Scaffold complete: 10 file(s) generated in /Users/raihan/forgebe
```

#### Generate specific template to output directory
```bash
forgebe scaffold --template node-express --output ./new-backend
# Output: Scaffold complete: 10 file(s) generated in ./new-backend
```

#### Preview files (dry-run) before generation
```bash
forgebe scaffold --template python-fastapi --dry-run
# Output:
# Dry-run: 12 file(s) to generate
#   app/
#   app/__init__.py
#   app/main.py
#   app/api/
#   ...
```

#### Force overwrite existing files
```bash
forgebe scaffold --force
# Overwrites existing files without prompting
```

### What Gets Generated

**Go Service:**
- `cmd/{project}/main.go` — Entry point
- `internal/handler/handler.go` — HTTP handlers
- `internal/service/service.go` — Business logic
- `go.mod` — Module definition
- `Makefile` — Build/test/run commands
- `.gitignore` — Standard Go ignores
- `README.md` — Project documentation

**Node Express:**
- `src/index.ts` — Application setup
- `src/routes/index.ts` — Route definitions
- `src/middleware/errorHandler.ts` — Error handling
- `package.json` — Dependencies and scripts
- `tsconfig.json` — TypeScript configuration
- `.gitignore` — Node ignores
- `README.md` — Getting started guide

**Python FastAPI:**
- `app/main.py` — Application entry
- `app/api/routes.py` — API endpoints
- `app/models/schemas.py` — Data schemas
- `requirements.txt` — Python dependencies
- `.gitignore` — Python ignores
- `README.md` — Setup instructions

---

## Phase 10.2: `forgebe adopt`

Adopt existing projects into ForgeBE without full onboarding.

### How It Works
1. Scans project directory
2. Auto-detects language, framework, CI, architecture
3. Creates ForgeBE profile automatically
4. Sets sensible defaults (80% coverage, TDD, safety-first)

### Examples

#### Adopt current directory
```bash
cd /path/to/existing-project
forgebe adopt
# Output:
# Project adopted successfully
#   Project ID: my-project_a1b2c3d4
#   Profile:    ~/.forgebe/projects/my-project_a1b2c3d4/project-profile.yaml
#   Language:   go
```

#### Adopt specific directory
```bash
forgebe adopt /path/to/project
```

#### Preview what would be adopted (dry-run)
```bash
forgebe adopt --dry-run
# Output:
# Dry-run: adopt /path/to/project
#   Project ID: project_xyz789
#   Language:   python
#
# Actions that would be performed:
#   - Create profile: project_xyz789
#   - Language detected: python
#   - Architecture: unknown
#   - Test framework: pytest
```

#### Force overwrite existing profile
```bash
forgebe adopt --force
# Re-scans and updates the profile
```

### What Gets Detected
- **Languages:** Go, Python, JavaScript/TypeScript, Java, Rust
- **Frameworks:** FastAPI, Express, Spring, Gin, etc.
- **Testing:** go test, pytest, Jest, JUnit
- **CI/CD:** GitHub Actions, GitLab CI, Jenkins
- **Package Managers:** Go modules, npm, pip, Maven, Cargo

---

## Phase 10.3: `forgebe verify`

Validate project profile against engineering policies.

### What Gets Checked
- Testing policy (required, coverage target, strategy)
- Dependency management (approval required, forbidden packages)
- Change policies (structure preservation, cross-file changes)
- Sensitive areas definition
- Source roots configuration
- Architecture pattern defined
- Forbidden paths

### Examples

#### Verify project profile
```bash
forgebe verify
# Output:
# ForgeBE Verify Report
# =====================
#
# ✅ [PASS] testing_policy
#    Testing is yes, strategy: tdd, coverage target: 80%
#
# ✅ [PASS] dependency_policy
#    Allow addition: no, require approval: yes
#
# ⚠️  [WARN] change_policy
#    Structure preservation with restricted cross-file changes may limit refactoring
#
# ✅ [PASS] sensitive_areas
#    Sensitive areas: Authentication / Authorization, Database migrations
#
# ⚠️  [WARN] source_roots
#    No source roots defined; AI context may miss key directories
#
# ---
# Summary: 3 PASS / 2 WARN / 0 FAIL
# Result: WARN
```

#### Verify specific project by ID
```bash
forgebe verify my-backend_abc123
```

#### Get JSON output for CI integration
```bash
forgebe verify --json
# Output: [{"check":"testing_policy","status":"PASS",...}, ...]
```

### Use in CI Pipeline
```bash
#!/bin/bash
forgebe verify --json | jq '.[] | select(.status=="FAIL")'
if [ $? -eq 0 ]; then
  echo "Policy violations found!"
  exit 1
fi
```

---

## Phase 10.4: `forgebe check`

Health gate for pre-commit hooks and CI pipelines.

### What Gets Checked
- Profile freshness (>7 days old = WARN)
- Repository exists and is accessible
- Basic verification passes
- Sync status up-to-date

### Examples

#### Run all health checks
```bash
forgebe check
# Output:
# ForgeBE Check
# =============
#
# ✅ profile_freshness: PASS
# ✅ repo_exists: PASS
# ✅ sync_status: PASS
# ✅ verify: PASS
#
# ---
# Result: PASS (4 passed, 0 warnings, 0 failures)
```

#### Get JSON for scripting
```bash
forgebe check --json
# Output:
# {
#   "status": "PASS",
#   "checks": {
#     "profile_freshness": "PASS",
#     "repo_exists": "PASS",
#     "verify": "PASS",
#     "sync_status": "PASS"
#   },
#   "total_checks": 4,
#   "passed": 4,
#   "warned": 0,
#   "failed": 0,
#   "message": "4 passed, 0 warnings, 0 failures"
# }
```

#### Use in pre-commit hook
```bash
#!/bin/bash
# .git/hooks/pre-commit

forgebe check
if [ $? -ne 0 ]; then
  echo "ForgeBE health check failed. Commit blocked."
  exit 1
fi
```

#### Use in GitHub Actions
```yaml
name: Health Gate
on: [pull_request, push]
jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: forgebe check
        env:
          FORGEBE_HOME: ~/.forgebe
```

---

## Phase 10.5: `forgebe prompt`

Generate context-aware prompts for AI tools (Claude, Cursor, Copilot, Hermes).

### Supported Modes
- `implement` — Code implementation prompt
- `review` — Code review checklist prompt
- `debug` — Debugging/troubleshooting prompt
- `plan` — Implementation planning prompt

### Examples

#### Generate implementation prompt
```bash
forgebe prompt implement "Add user authentication endpoint"
# Output:
# # Implement Task
#
# ## Project Context
# - Project: forgebe
# - Language: go
# - Framework: go test
# - Architecture: clean/hexagonal
# - Testing: tdd
# - Sensitive Areas: Authentication / Authorization, Database migrations
# - Priority: safety
#
# ## Task
# Add user authentication endpoint
#
# ## Instructions
# - Follow the project's architecture pattern.
# - Write tests first (TDD) if testing strategy is tdd.
# - Keep changes within source roots.
# - Do NOT modify sensitive areas without explicit approval.
# - Do NOT add new dependencies without approval.
# - Ensure code compiles and tests pass before submitting.
```

Copy to Claude/Cursor/Hermes prompt box.

#### Generate review prompt
```bash
forgebe prompt review "Check the payment handler"
# Output:
# # Code Review
#
# ## Project Context
# - Project: forgebe
# - Language: go
# - Architecture: clean/hexagonal
# - Testing: tdd
# - Sensitive Areas: Authentication / Authorization, Database migrations
#
# ## Review Focus
# Check the payment handler
#
# ## Checklist
# - Does the code follow the project architecture?
# - Are tests present and meaningful?
# - Any security concerns in sensitive areas?
# - Does it follow the naming and style conventions?
# - Does it pass linting?
# - Is it properly formatted?
# - Are error paths handled?
```

#### Generate plan prompt
```bash
forgebe prompt plan "Migrate database from MySQL to PostgreSQL"
```

#### Generate debug prompt
```bash
forgebe prompt debug "Fix panic in request middleware"
```

### Why Context-Aware?

Traditional prompts lose context. `forgebe prompt` includes:
- **Project specifics:** Language, framework, architecture
- **Constraints:** Testing strategy, dependency approval, sensitive areas
- **Standards:** Linting requirements, code review rules, delivery priority
- **Warnings:** Do not modify sensitive areas, approval requirements

This ensures AI-generated code aligns with your project's standards automatically.

---

## Integration Patterns

### Local Development Workflow
```bash
# 1. Generate new feature scaffold
forgebe scaffold --template go-api --output ./new-feature

# 2. Generate implementation prompt
forgebe prompt implement "Add user endpoints to new-feature"

# 3. Paste prompt into Claude/Cursor

# 4. Verify generated code against policies
forgebe verify

# 5. Check health before commit
forgebe check
```

### CI/CD Pipeline
```yaml
# .github/workflows/quality.yml
name: Quality Gate
on: [pull_request]

jobs:
  forgebe:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: forgebe verify --json | jq 'if any(.status=="FAIL") then error("Policy violations") else . end'
      - run: forgebe check
```

### Pre-commit Hook
```bash
#!/bin/bash
# .git/hooks/pre-commit

# Verify no policy violations
forgebe verify --json | jq -e 'all(.status != "FAIL")' || {
  echo "ForgeBE verification failed"
  exit 1
}

# Health check must pass
forgebe check || exit 1
```

### Adoption Workflow (New Team Member)
```bash
# 1. Clone existing project
git clone https://github.com/company/backend.git
cd backend

# 2. Adopt into ForgeBE
forgebe adopt

# 3. Generate onboarding prompt
forgebe prompt plan "Get familiar with project structure and architecture"

# 4. Paste into Claude for guidance
```

---

## Best Practices

### Scaffold
- Use templates matching your stack exactly
- Review generated code before committing
- Customize boilerplate for your policies
- Commit scaffold as baseline, then build on top

### Adopt
- Run adopt in project root (where `.git/` exists)
- Review generated profile — adjust sensitive areas
- Run `forgebe verify` to check policy alignment
- Commit profile to git for team consistency

### Verify
- Run before every commit
- Integrate into CI pipeline
- Use JSON output for automation
- Address WARN items proactively

### Check
- Use in pre-commit hooks
- Use in CI quality gates
- Fail CI if check fails
- Monitor freshness — update stale profiles

### Prompt
- Use for all coding tasks (implement, review, debug, plan)
- Copy entire output into AI tool
- Include generated prompt in PR description for transparency
- Review AI output against context before merging

---

## Troubleshooting

### `forgebe scaffold` fails with "file already exists"
Use `--force` to overwrite or clean target directory first.

### `forgebe adopt` shows "project already has a profile"
Use `--force` to re-scan and update, or remove the existing profile from `~/.forgebe/projects/`.

### `forgebe verify` shows multiple WARN
Review the specific warnings — they indicate missing context or non-standard setup. Update profile accordingly.

### `forgebe check` fails in CI
Ensure `FORGEBE_HOME` is set correctly and profile is committed to repo.

### `forgebe prompt` output looks generic
Verify profile has sensible values — especially language, architecture, sensitive areas. Run `forgebe profile show` to inspect.

---

## Summary

Phase 10 transforms ForgeBE from a **profile/context tool** into a **complete backend engineering framework**:

- **Scaffold** — Standardize new projects instantly
- **Adopt** — Onboard existing projects with one command
- **Verify** — Enforce engineering standards automatically
- **Check** — Quality gate for automation workflows
- **Prompt** — AI-generated code that respects your standards

**Coverage:** 78.4% across all packages  
**Tests:** 60+ new tests (RED-GREEN-REFACTOR TDD)  
**Commands:** 13 total (7 existing + 5 new + 1 root)  
**Documentation:** Full with 20+ usage examples  

Ready for production. 🚀
