# ForgeBE — Awareness Review (`forgebe review`) — Design Spec

> **Date:** 2026-06-24
> **Status:** Approved design, pending implementation plan
> **Scope:** Phase 1 of the "Aware-by-default, easy-to-adopt" roadmap (Direction A + B)

## Context

ForgeBE today is strong on the **input** side of its mission ("AI mengikuti kemauan
developer"): it captures project context + an engineering contract in `~/.forgebe/`
and exports it to any AI tool (Claude, Cursor, Copilot, Hermes, Generic).

It is thin on the **output/feedback** side ("developer tetap aware dengan hasil AI").
The existing `verify` package validates the *profile itself*, and
`verify.VerifyAgainstRepo` is explicitly *"not yet implemented"*. There is no command
that inspects **actual code changes** and tells the developer what changed and what
deviates from their own policy.

`forgebe review` (Awareness Review) fills exactly this gap. It is the differentiating
"aha" feature for community adoption: demo-able, screenshot-able, and uniquely tied to
ForgeBE's policy model.

### Non-negotiable principle: AI-agnostic

The whole point of ForgeBE is to **not** lock to a single AI vendor. Because this
feature reads a **git diff** — the result of changes, not the agent that made them — it
is inherently tool-agnostic. It works identically whether the diff came from Claude
Code, Cursor, Copilot, Hermes, or hand-typing.

**Hard rule for implementation:** no code path in the `review` or `git` packages may
import or reference any AI vendor SDK, assume any vendor-specific session format, or
branch on "which AI tool". Violating this defeats the product's reason to exist.

## Goal

A new CLI command `forgebe review` that:

1. Computes the git diff for a configurable range (working tree / staged / since a ref).
2. Loads the project's policy from its ForgeBE profile.
3. Runs a set of **diff ↔ policy** rules and produces a prioritized **awareness report**
   ("here is what changed, and here is what you need to be consciously aware of").
4. Is informational by default (exit 0) and can act as a gate under `--strict` (exit 1
   on any FAIL) so it can later power a pre-commit hook and GitHub Action (Phase 2).

## Posture: "aware by default, gate on demand"

- **Default** (`forgebe review`): informational. Always exit 0. The purpose is to *make
  the developer aware*, not to judge. This honors the "tetap aware" philosophy.
- **`--strict`**: turns FAIL findings into a non-zero exit, for CI/hook gating (Phase 2).

## Architecture

Follows the existing package layout and the `verify` / `check` patterns.

```
forgebe review            internal/cli/review.go      (cobra cmd, mirrors verify.go)
        |
        v
internal/review/          NEW package — the "awareness brain"
   reviewer.go              Reviewer{profile, diff}; Run() -> []Finding
   rules.go                 individual diff<->policy rules
   report.go                Text() + JSON(), mirrors verify/check formatting
        |  uses
        v
internal/git/             NEW package (currently empty dir)
   diff.go                  ChangedFiles(rangeSpec) via `git diff --name-status --numstat`
   provider.go              DiffProvider interface (testability seam)
        |  loads
        v
internal/profile (Policy)  EXISTING — source of the rules
```

### Key types (shape, not final signatures)

```go
// internal/git
type FileChange struct {
    Path    string // new path
    OldPath string // for renames
    Status  string // A, M, D, R
    Added   int    // lines added (numstat)
    Deleted int    // lines deleted (numstat)
}

type DiffProvider interface {
    ChangedFiles(spec RangeSpec) ([]FileChange, error)
}

type RangeSpec struct {
    Staged bool   // git diff --cached
    Since  string // git diff <ref>...; empty => working tree vs HEAD
}

// internal/review
type Finding struct {
    Rule     string // forbidden_path, dependency_added, missing_test, ...
    Severity string // FAIL, WARN, INFO
    Path     string // file the finding relates to (optional)
    Message  string
}

type Reviewer struct {
    profile *profile.ProjectProfile
    changes []git.FileChange
}

func (r *Reviewer) Run() []Finding
```

`Finding` deliberately mirrors `verify.Result` / `check.Result` so reporting and JSON
output stay consistent across the tool. JSON output and exit codes are built in from day
one so Phase 2 (hook + GitHub Action) needs no rework.

## Rule set (MVP)

High-value, low-cost rules that map directly onto fields already present in
`profile.Policy` and `profile.Areas`:

| Rule key            | Policy source                                                   | Default severity | Notes |
|---------------------|-----------------------------------------------------------------|------------------|-------|
| `forbidden_path`    | `Policy.Changes.ForbiddenPaths`                                 | FAIL             | Any changed file matching a forbidden glob/pattern. |
| `sensitive_area`    | `Areas.SensitiveAreas`, `Areas.PublicAPIs`, `Changes.RequireApproval` | WARN       | Touched code that needs the developer's conscious sign-off. |
| `dependency_added`  | `Policy.Dependencies.AllowAddition` / `RequireApproval`         | WARN or FAIL     | Detected via changes to manifest files (go.mod, package.json, requirements.txt, etc.). FAIL if `AllowAddition=false`. |
| `dependency_forbidden` | `Policy.Dependencies.Forbidden`                              | FAIL             | A forbidden dependency name appears in a manifest diff. |
| `missing_test`      | `Policy.Testing.Required` / `UnitRequired` + `Areas.TestRoots`  | WARN             | A new/modified source file with no corresponding test file changed in the same diff. |
| `summary`           | —                                                               | INFO             | Counts: N files changed, +X/-Y lines, grouped by area (source root). |

### Matching rules

- Path patterns reuse the same matching approach already used for forbidden/require
  patterns in the profile (glob-style). Keep matching logic in one helper to avoid
  divergence.
- Manifest detection is a fixed map of filename → ecosystem (go.mod, package.json,
  requirements.txt, pyproject.toml, Cargo.toml, pom.xml, build.gradle). Parsing depth is
  minimal for MVP: detect *that* a dependency line was added, and extract the dependency
  token for `dependency_forbidden`. Deep semantic parsing is out of scope.
- `missing_test`: a changed file under a `SourceRoots` entry with a code extension, where
  no changed file in a `TestRoots` entry (or matching `*_test.*` convention) references
  it. MVP uses a heuristic (same base name / same package dir), not call-graph analysis.

## Output / UX

Two formats, mirroring `verify` / `check`:

**Text (default, human-facing):**

```
ForgeBE Awareness Review
========================
Range: working tree (vs HEAD)   •   12 files changed  (+340 / -58)

❌ FAIL  forbidden_path        internal/auth/token.go touched (Areas.Sensitive)
⚠️  WARN  dependency_added      go.mod: + github.com/foo/bar (needs approval)
⚠️  WARN  missing_test          internal/payment/charge.go is new, no test
ℹ️  INFO  summary               domain:3  app:5  infra:4

---
3 need your awareness  •  Result: REVIEW NEEDED (exit 0)
```

**JSON (`--json`):** array of `Finding` plus a summary object — for Phase 2 (hook / CI /
PR comment).

### Flags

| Flag             | Meaning                                                          |
|------------------|-----------------------------------------------------------------|
| `[project-id]`   | Positional, resolved like `verify` (`resolveProjectIDArg`).      |
| `--staged`       | Review staged changes (`git diff --cached`) — for pre-commit.    |
| `--since <ref>`  | Review changes since a ref (`git diff <ref>...HEAD`) — for CI/PR.|
| `--json`, `-j`   | JSON output.                                                    |
| `--strict`       | Exit 1 if any FAIL finding (gate mode).                          |

Default range (no flag) = working tree vs HEAD.

### Exit codes

- Default: always `0` (informational).
- `--strict`: `1` if any `FAIL` finding, else `0`.
- Operational errors (not a git repo, profile not found): non-zero with a clear message.

## Data flow

```
git diff (range) ──▶ internal/git.ChangedFiles() ──┐
                                                    ├─▶ review.Reviewer.Run() ─▶ []Finding ─▶ report.Text()/JSON()
profile.Policy (load by project-id) ───────────────┘                                           │
                                                                                         exit code (posture)
```

## CLI integration

- New file `internal/cli/review.go` with `newReviewCmd() *cobra.Command`, mirroring
  `verify.go` (uses `resolveProjectIDArg`, `loadProfileByID`, `cmd.OutOrStdout()`).
- Register in `internal/cli/root.go`: `cmd.AddCommand(newReviewCmd())`.

## Testing strategy (TDD — per CLAUDE.md)

- **Test-first, table-driven** (repo convention). Each rule gets a table of PASS/WARN/FAIL
  cases driven by a synthetic `[]git.FileChange` + a crafted `Policy`.
- **`internal/git`** tested two ways: (1) against a real temp repo (`t.TempDir()` +
  `git init` + commits) for the real provider; (2) reviewer logic tested against a fake
  `DiffProvider` so rule tests need no real git.
- Reuse `testdata/repos/go-layered` as a realistic fixture.
- **Coverage ≥ 80%** to satisfy the existing SonarCloud quality gate.
- **AI-agnostic guard test:** a test (or CI grep) asserting no AI-vendor identifier
  appears in `internal/review` / `internal/git` source.

## Non-goals (Phase 1)

- ❌ Cross-layer architecture import analysis (per-language parsing) → Phase 3.
- ❌ Human-vs-AI attribution / provenance / chain-of-truth → deferred (set aside).
- ❌ Pre-commit hook & GitHub Action → Phase 2 (but Phase 1 is built CI-ready so Phase 2
  is wiring only).
- ❌ Auto-fix. Review only *raises awareness*; it never edits code.
- ❌ Deep semantic dependency parsing / version constraint analysis.

## Acceptance criteria

1. `forgebe review` on a dirty working tree prints an awareness report and exits 0.
2. `forgebe review --strict` exits 1 when a forbidden path or forbidden dependency is in
   the diff, else 0.
3. `forgebe review --staged` and `forgebe review --since <ref>` scope the diff correctly.
4. `forgebe review --json` emits valid JSON (findings + summary).
5. All MVP rules covered by table-driven tests; package coverage ≥ 80%.
6. No AI-vendor reference anywhere in `internal/review` or `internal/git`.
7. `go vet ./...` clean; `gofmt` clean; binary builds.

## Roadmap context (for reference, not this spec)

| Phase | Name | This spec? |
|-------|------|-----------|
| **1** | **Awareness Review** (`forgebe review`) | **Yes** |
| 2 | Distribute the wedge: pre-commit hook + GitHub Action posting the report on PRs | No |
| 3 | Deeper awareness: over-time drift + provenance/journal (absorbs chain_of_truth) | No |
| 4 | Frictionless onboarding: instant-run/zero-config, plugin packaging, docs/website | No |
