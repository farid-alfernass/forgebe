# Phase 9 Quality, Distribution, Sonar, and Coverage Plan

> **For Hermes:** Execute in small verified increments. After each implementation chunk: run focused tests, then repo-wide regression, then update this plan if scope changes.

**Goal:** Bring ForgeBE from production-ready to continuously quality-governed: distribution-complete, Sonar-enabled, recursively improved, and with automated test coverage raised to **at least 80%**.

**Architecture:** Work in concentric layers. First stabilize release/distribution metadata. Then wire static analysis (SonarQube Cloud). Then attack low-coverage/high-risk packages with tests and small refactors. After each wave, run full regression + smoke tests. Do not chase 80% by adding brittle tests around implementation details; prefer stable behavior-driven tests around CLI commands, profile/store behavior, export/import edge cases, and discovery logic.

**Tech Stack:** Go 1.22, Cobra CLI, GitHub Actions, GoReleaser, Homebrew tap, SonarQube Cloud.

---

## Phase A — Audit and Planning

### Task A1: Record baseline state
**Objective:** Capture repo/tag/coverage baseline before modifications.

**Files:**
- Update: `docs/plans/2026-06-07-phase9-quality-distribution-sonar-plan.md`

**Steps:**
1. Record latest commit, latest tags, current coverage total, and package coverage.
2. Record distribution status: release assets, Homebrew formula status, release workflow status.
3. Record Sonar status: configured or not configured, token available, org/project keys still needed or not.

**Verification:**
- `git status --short` clean or understood
- `go test -coverprofile=/tmp/forgebe.out ./... && go tool cover -func=/tmp/forgebe.out | tail -1`

**Baseline captured on 2026-06-08:**
- Latest commit before further Phase 9 work: `2d48058` (`feat: enrich context status summary and JSON output`)
- Latest tag: `v1.1.0`
- Coverage total: `54.3%`
- Lowest-coverage/high-value areas observed:
  - `internal/cli/context_cmd.go` (`newStatusCmd`, `newWatchCmd`)
  - `internal/cli/doctor.go`
  - `internal/cli/profile.go`
  - `internal/profile/store.go` metadata/discovery persistence paths
- Distribution status:
  - Release workflow exists: `.github/workflows/release.yml`
  - GoReleaser config exists: `.goreleaser.yaml`
  - Homebrew tap publishing configured via `HOMEBREW_TAP_GITHUB_TOKEN`
  - Current branch has understood but uncommitted non-Phase-1.1 changes (`.github/workflows/ci.yml`, `internal/profile/validator_test.go`, untracked planning/sonar files)
- Sonar status:
  - CI workflow includes Sonar scan + quality gate steps
  - `sonar-project.properties` is present in working tree as untracked work and still needs review/finalization
  - Expected secret: `SONAR_TOKEN`
  - SonarCloud host configured: `https://sonarcloud.io`

---

## Phase B — Distribution Completion

### Task B1: Verify version/build metadata end-to-end
**Objective:** Ensure binary version data is injected correctly for local build and release build.

**Files:**
- Review/modify: `internal/cli/version.go`
- Review/modify: `Makefile`
- Review/modify: `cmd/forgebe/main.go`

**Steps:**
1. Build with `make build`.
2. Verify `forgebe version` text output includes version, commit, build date.
3. Verify `forgebe version --json` emits structured output.
4. Add or adjust tests for version formatting if practical.

**Verification:**
- `./bin/forgebe version`
- `./bin/forgebe version --json | jq .`

### Task B2: Verify GoReleaser output naming and checksums
**Objective:** Ensure `.goreleaser.yaml` matches actual release assets and Homebrew expectations.

**Files:**
- Review/modify: `.goreleaser.yaml`
- Review/modify: `.github/workflows/release.yml`

**Steps:**
1. Confirm archive names for darwin/linux amd64/arm64.
2. Confirm checksum file name and brew stanza compatibility.
3. If needed, run local snapshot build.

**Verification:**
- `goreleaser release --snapshot --clean` or equivalent if available
- compare output names to formula URLs

### Task B3: Finalize Homebrew tap flow
**Objective:** Make brew installation reproducible and documented.

**Files:**
- Review/modify: `README.md`
- Review/modify: `docs/brew-tap.md`
- Review/modify: `.goreleaser.yaml`
- External repo context: `homebrew-tap/Formula/forgebe.rb`

**Steps:**
1. Verify formula URLs and sha256 values against released assets.
2. Verify local install from tap works.
3. Document exact one-time setup for `homebrew-tap` repo and token.
4. If full auto-publish still depends on external secret/repo setup, explicitly document remaining manual step.

**Verification:**
- `brew tap farid-alfernass/tap <path-or-remote>`
- `brew install farid-alfernass/tap/forgebe`
- `forgebe version`

---

## Phase C — SonarQube Cloud Enablement

### Task C1: Add scanner configuration
**Objective:** Make Sonar analysis reproducible in CI and locally.

**Files:**
- Create/modify: `sonar-project.properties`
- Review/modify: `.github/workflows/ci.yml`
- Review/modify: `.github/workflows/release.yml` (only if needed)

**Steps:**
1. Configure sources, tests, exclusions, coverage report path.
2. Add CI job or step to run Sonar scanner.
3. Use secret token in CI, never commit token.

**Verification:**
- config file valid
- scanner command composes successfully

### Task C2: Run baseline Sonar analysis
**Objective:** Produce first issue list for recursive remediation.

**Files:**
- No required file changes initially.

**Steps:**
1. Run scanner locally if possible using env var for token.
2. Capture issues: bugs, code smells, vulnerabilities, hotspots, duplication.
3. Prioritize real quality risks over cosmetic warnings.

**Verification:**
- successful scanner upload or clear blocker message
- issue summary recorded in plan or follow-up doc

---

## Phase D — Coverage Ramp to 80%+

### Task D1: Raise `internal/cli` coverage above 70%
**Objective:** Cover command construction and behavior helpers.

**Files:**
- Modify: `internal/cli/cli_test.go`
- Add tests for: `version.go`, `scan.go`, `doctor.go`, `profile.go`, `export.go`, `import.go`, `brief.go`

**Steps:**
1. Prefer command execution tests with temp ForgeBE roots.
2. Capture stdout/stderr and assert stable output fragments / JSON fields.
3. Cover success and failure branches.
4. Refactor minimal seams only if needed for testability.

**Verification:**
- `go test -cover ./internal/cli`

### Task D2: Raise `internal/export` coverage above 80%
**Objective:** Cover more bundle import/export edge cases and summary rendering branches.

**Files:**
- Modify: `internal/export/bundle_test.go`
- Modify/add: `internal/export/summary_test.go`

**Steps:**
1. Add malformed zip path tests.
2. Add symlink/traversal regression tests if possible.
3. Add missing/empty bundle content scenarios.
4. Add summary rendering expectations for sparse profiles.

**Verification:**
- `go test -cover ./internal/export`

### Task D3: Raise `internal/profile` coverage above 85%
**Objective:** Cover validators, stores, sample profile paths, and edge behavior.

**Files:**
- Modify/add: `internal/profile/*_test.go`

**Steps:**
1. Test invalid profile cases and expected issues.
2. Test store save/load round-trip more deeply.
3. Test project id determinism edge cases.
4. Test sample profile assumptions.

**Verification:**
- `go test -cover ./internal/profile`

### Task D4: Raise `internal/discovery` coverage above 80%
**Objective:** Cover skipped dirs, fallback detection, edge manifests, and empty repos.

**Files:**
- Modify/add: `internal/discovery/*_test.go`
- Add fixtures under `testdata/` as needed

**Steps:**
1. Add empty repo scan case.
2. Add skip-dir regression for `testdata`, `fixtures`, `examples`.
3. Add root-manifest prioritization case.
4. Add unknown language fallback cases.

**Verification:**
- `go test -cover ./internal/discovery`

### Task D5: Optionally add `cmd/forgebe` smoke coverage
**Objective:** Lift overall total modestly by testing entrypoint behavior if cheap.

**Files:**
- Add: `cmd/forgebe/main_test.go` only if practical and non-brittle.

**Verification:**
- `go test -cover ./cmd/forgebe`

### Task D6: Recompute total coverage and close gaps recursively
**Objective:** Iterate until total >= 80%.

**Steps:**
1. Run full repo coverage.
2. Identify next-lowest meaningful package.
3. Add tests or tiny refactors.
4. Repeat until target reached or blocker documented.

**Verification:**
- `go test -coverprofile=/tmp/forgebe.out ./...`
- `go tool cover -func=/tmp/forgebe.out | tail -1`

---

## Phase E — Recursive Quality Remediation

### Task E1: Fix code review findings
**Objective:** Run another independent review after coverage improvements and fix real maintainability/security issues.

**Files:**
- As indicated by findings.

**Verification:**
- Review summary with no open high-severity issues.

### Task E2: Fix Sonar issues recursively
**Objective:** Work down Sonar issues in priority order.

**Process:**
1. Bugs and vulnerabilities first.
2. Reliability/maintainability next.
3. Minor smells last if they improve readability with low risk.
4. Re-run tests after every wave.

**Verification:**
- issue count drops materially
- no new regressions

---

## Phase F — Final Verification and Routine

### Task F1: Full regression and smoke suite
**Objective:** Verify all commands still work after refactors.

**Commands:**
- `gofmt -w ./cmd ./internal`
- `go vet ./...`
- `go test -count=1 ./...`
- `go build -o bin/forgebe ./cmd/forgebe`
- command smoke suite: `version`, `scan`, `doctor`, `profile show`, `brief`, `export summary`, `export adapter`, `export bundle`, `import`

### Task F2: Add routine quality process docs
**Objective:** Document ongoing development cadence.

**Files:**
- Update: `CONTRIBUTING.md`
- Update: `README.md` or create `docs/quality-routine.md`

**Content:**
- Weekly regression cadence
- Pre-release checklist
- Sonar review cadence
- Coverage floor policy
- Bug triage rhythm

### Task F3: Commit and push with clear milestones
**Objective:** Keep changes reviewable.

**Commit strategy:**
- docs/plan
- sonar/config
- coverage wave 1
- coverage wave 2
- quality fixes
- final verification/docs

---

## Success Criteria

- Homebrew path documented and reproducible.
- SonarQube Cloud wired into repo and baseline analysis completed.
- Review findings recursively addressed.
- Total test coverage **>= 80%**.
- Full regression and smoke test pass.
- Ongoing quality routine documented.

---

## Notes

- Do **not** commit secrets. Sonar token must be used through environment variables / CI secrets only.
- Favor stable black-box tests over white-box brittle assertions.
- If 80% total coverage becomes blocked by low-value or brittle areas, first refactor for testability rather than writing meaningless tests.
- Keep Homebrew tap content in the tap repo, not the main ForgeBE repo, except docs/config references.
