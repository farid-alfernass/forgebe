# ForgeBE CLI Reference

## Global Flags

All commands support:

```
--json, -j       Output in JSON format (where applicable)
--help, -h       Show help for command
```

---

## Commands

### `forgebe init`

Initialize a ForgeBE project with guided onboarding or auto-discovery.

**Usage:**
```bash
forgebe init [path] [flags]
```

**Arguments:**
- `path` — Repository path (default: current directory)

**Flags:**
- `--non-interactive` — Skip interactive questions, use auto-discovery defaults

**Example:**
```bash
# Interactive onboarding
forgebe init

# Auto-discovery defaults
forgebe init --non-interactive .

# Initialize a remote path
forgebe init /path/to/project
```

**Output:**
```text
ForgeBE initialized successfully
Project ID: forgebe_df5dcae6
Profile: /Users/raihan/.forgebe/projects/forgebe_df5dcae6/project-profile.yaml
Contract: /Users/raihan/.forgebe/projects/forgebe_df5dcae6/engineering-contract.md
```

---

### `forgebe scan`

Scan a repository and infer its stack, architecture, testing framework, and CI setup.

**Usage:**
```bash
forgebe scan [path] [--json]
```

**Arguments:**
- `path` — Repository path (default: current directory)

**Output (YAML):**
```yaml
timestamp: 2026-06-07T09:00:00Z
repo_path: /Users/raihan/forgebe
detections:
  languages:
    - language: go
      confidence: high
      files: 48
      manifest: go.mod
  architecture:
    style: go-service-layout
    confidence: medium
  testing:
    framework: go test
  dependencies:
    package_manager: go modules
  ci:
    providers:
      - github-actions
```

**Output (JSON):**
```bash
forgebe scan --json
```

---

### `forgebe doctor`

Validate ForgeBE installation and check project profile health.

**Usage:**
```bash
forgebe doctor [project-id] [--json]
```

**Arguments:**
- `project-id` — Project ID (default: auto-detect from active projects)

**Output (Text):**
```text
ForgeBE Doctor — Project Health Check
  ✅ ForgeBE root: /Users/raihan/.forgebe
  ✅ Projects store: /Users/raihan/.forgebe/projects

  Profile: forgebe_df5dcae6
  Project: forgebe (go)
  Status: All checks passed
```

**Output (JSON):**
```bash
forgebe doctor --json
```

---

### `forgebe profile show`

Display the active project profile in human-readable format.

**Usage:**
```bash
forgebe profile show [project-id] [--json]
```

**Arguments:**
- `project-id` — Project ID (default: auto-detect)

**Output (Text):**
```text
Project:    forgebe
Profile ID: forgebe_df5dcae6
Language:   go
Framework:  go test
Arch:       go-service-layout
Maturity:   existing codebase
Type:       service
Testing:    tdd
Delivery:   hybrid (safety priority)
Repo:       /Users/raihan/forgebe
Sensitive:  [Authentication / Authorization Database migrations]
```

**Output (JSON):**
```bash
forgebe profile show --json
```

---

### `forgebe brief`

Generate a compact AI prompt for any supported tool.

**Usage:**
```bash
forgebe brief <tool> [project-id]
```

**Arguments:**
- `tool` — Tool name: `claude`, `cursor`, `copilot`, `hermes`, `generic`
- `project-id` — Project ID (default: auto-detect)

**Output:**
```text
Project: forgebe
Language: go
Architecture: go-service-layout
Testing: tdd
Dependency approval required: true
Sensitive areas: Authentication / Authorization, Database migrations, Public API contracts
```

**Example:**
```bash
forgebe brief claude forgebe_df5dcae6
# Copy output into your Claude chat
```

---

### `forgebe export`

Export project context in various formats.

#### `forgebe export summary`

Export a full markdown project summary.

**Usage:**
```bash
forgebe export summary [project-id] [--json] [--output FILE]
```

**Flags:**
- `--output, -o FILE` — Write to file instead of stdout

**Output:**
```markdown
# ForgeBE Project Summary

## Overview
- **Project**: forgebe
- **Profile ID**: forgebe_df5dcae6
- **Repository**: /Users/raihan/forgebe

## Stack
- **Primary Language**: go
- **Framework**: go test
...
```

**Example:**
```bash
# Write to file
forgebe export summary -o /tmp/summary.md

# Output as JSON
forgebe export summary --json
```

#### `forgebe export adapter`

Export a tool-specific adapter file.

**Usage:**
```bash
forgebe export adapter <tool> [project-id] [--output FILE]
```

**Arguments:**
- `tool` — Tool name: `claude`, `cursor`, `copilot`, `hermes`, `generic`
- `project-id` — Project ID (default: auto-detect)

**Flags:**
- `--output, -o FILE` — Write to custom path (default: `~/.forgebe/exports/`)

**Example:**
```bash
# Export Claude adapter to default location
forgebe export adapter claude

# Export with custom output
forgebe export adapter cursor -o ./AI_CONTEXT.md

# Export for Copilot
forgebe export adapter copilot
```

#### `forgebe export bundle`

Export a portable `.forgebe.zip` bundle containing all project context.

**Usage:**
```bash
forgebe export bundle [project-id] [--output FILE]
```

**Flags:**
- `--output, -o FILE` — Write to custom path (default: `~/.forgebe/exports/<id>.forgebe.zip`)

**Example:**
```bash
# Export to default location
forgebe export bundle

# Export to custom path
forgebe export bundle -o ~/Desktop/forgebe-backup.zip

# Transfer to another machine
scp forgebe-backup.zip user@remote:~/
```

---

### `forgebe sync`

Generate and update AI context files directly in the project root from the active ForgeBE profile.

**Usage:**
```bash
forgebe sync [project-id] [--dry-run] [--force] [--json]
```

**Arguments:**
- `project-id` — Project ID (default: auto-detect from local ForgeBE projects)

**Flags:**
- `--dry-run` — Show files that would be written without modifying the repository
- `--force, -f` — Rewrite files even when content is already up-to-date

**Generated files:**
- `CLAUDE.md`
- `AGENTS.md`
- `.cursorrules`
- `.github/copilot-instructions.md`
- `AI_CONTEXT.md`

**Safety note:**
`forgebe sync` writes to the repository root. Use `--dry-run` before first use on an existing project.

**Example:**
```bash
# Preview changes
forgebe sync --dry-run

# Write/update context files
forgebe sync
```

---

### `forgebe status`

Show whether generated AI context files are present and up-to-date for a project.

**Usage:**
```bash
forgebe status [project-id] [--json]
```

**Example output:**
```text
ForgeBE Context Status for project: forgebe_df5dcae6
Repository: /Users/raihan/forgebe
Last sync: 2026-06-07T18:09:55+07:00

Synced files:
  - CLAUDE.md (claude)
  - AGENTS.md (hermes)

All context files are up-to-date.
```

---

### `forgebe watch`

Run a long-lived watcher that performs an initial sync and automatically re-syncs AI context files when relevant project files change.

**Usage:**
```bash
forgebe watch [project-id]
```

**Watched file categories:**
- Language/package manifests: `go.mod`, `package.json`, `Cargo.toml`, `pyproject.toml`, etc.
- Source roots: `cmd/`, `internal/`, `pkg/`, `src/`, `api/`, etc.
- Build/config files: `Dockerfile`, `Makefile`, `.env`, `tsconfig.json`, etc.

**Example:**
```bash
forgebe watch forgebe_df5dcae6
```

Press `Ctrl+C` to stop.

---

### `forgebe import`

Restore a `.forgebe.zip` bundle into local ForgeBE storage.

**Usage:**
```bash
forgebe import <bundle-path> [target-dir]
```

**Arguments:**
- `bundle-path` — Path to the `.forgebe.zip` file
- `target-dir` — Extraction root (default: `~/.forgebe/projects/`)

**Example:**
```bash
# Import to default location
forgebe import ~/forgebe-backup.zip

# Import to custom directory
forgebe import ~/forgebe-backup.zip /tmp/restore

# Verify import
forgebe doctor
```

---

### `forgebe version`

Display ForgeBE version.

**Usage:**
```bash
forgebe version
```

**Output:**
```text
dev
```

---

## Environment Variables

ForgeBE respects these environment variables:

| Variable | Purpose | Default |
|----------|---------|---------|
| `FORGEBE_HOME` | Override ForgeBE home directory | `~/.forgebe` |

**Example:**
```bash
export FORGEBE_HOME=/custom/path
forgebe init
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error (missing project, invalid profile, etc.) |
| `2` | Invalid arguments or flags |

---

## Examples

### Complete Workflow

```bash
# 1. Initialize a project
cd my-project
forgebe init

# 2. Verify setup
forgebe doctor

# 3. Show profile
forgebe profile show

# 4. Generate AI brief
forgebe brief claude | pbcopy  # macOS
# or
forgebe brief claude > brief.txt  # Save to file

# 5. Export adapter for Cursor
forgebe export adapter cursor

# 6. Create a portable bundle for sharing
forgebe export bundle -o ~/my-project.forgebe.zip
```

### Scripting with JSON

```bash
# Get profile as JSON
PROFILE=$(forgebe profile show --json)
PROJECT_ID=$(echo $PROFILE | jq -r '.metadata.profile_id')
LANGUAGE=$(echo $PROFILE | jq -r '.stack.primary_language')

echo "Project ID: $PROJECT_ID"
echo "Language: $LANGUAGE"

# Export and check bundle
forgebe export bundle --json
if [ $? -eq 0 ]; then
  echo "Bundle exported successfully"
fi
```

---

## Troubleshooting

### "Project not found"

ForgeBE couldn't find a ForgeBE-initialized project. Run:

```bash
forgebe init
```

### "Profile validation failed"

The project profile is corrupt or incomplete. Try:

```bash
forgebe doctor
```

to see specific issues.

### "Permission denied" on export

Ensure `~/.forgebe/` is writable:

```bash
mkdir -p ~/.forgebe/exports
chmod 700 ~/.forgebe
```

## `forgebe review`

Compare the current git diff against your project policy and print an
awareness report — what changed and what needs your conscious review.
Works with any AI tool because it inspects the diff, not the agent.

```bash
forgebe review                 # working tree vs HEAD (informational, exit 0)
forgebe review --staged        # staged changes (for pre-commit)
forgebe review --since main    # changes since a ref (for CI/PR)
forgebe review --json          # machine-readable output
forgebe review --strict        # exit 1 if any FAIL finding (gate mode)
```

Rules checked (mapped to your profile policy):

| Rule | Policy source | Severity |
|------|---------------|----------|
| `forbidden_path` | `policy.changes.forbidden_paths` | FAIL |
| `sensitive_area` | `areas.sensitive_areas` / `areas.public_apis` / `policy.changes.require_approval` | WARN |
| `dependency_added` | `policy.dependencies.allow_addition` / `require_approval` | FAIL/WARN |
| `dependency_forbidden` | `policy.dependencies.forbidden` | FAIL |
| `missing_test` | `policy.testing.required` + `areas.test_roots` | WARN |
| `summary` | — | INFO |
