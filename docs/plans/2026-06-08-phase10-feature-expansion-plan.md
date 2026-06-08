# Phase 10 Feature Expansion Implementation Plan

> **For Hermes:** Follow `test-driven-development` and execute Phase 10 sequentially with explicit quality gates after each sub-phase.

**Goal:** Menambahkan 5 command baru yang membuat ForgeBE bukan hanya profile/context manager, tapi framework backend engineering yang lebih operasional untuk berbagai project: `scaffold`, `adopt`, `verify`, `check`, dan `prompt`.

**Architecture:** Setiap capability baru dipisah per package domain (`internal/scaffold`, `internal/adopt`, `internal/verify`, `internal/check`, `internal/prompt`) dan diintegrasikan tipis lewat `internal/cli`. Seluruh command harus local-first, git-safe, support `--dry-run` bila relevan, serta mengikuti profile/project contract yang sudah ada.

**Tech Stack:** Go 1.22, Cobra CLI, existing ForgeBE profile/store/discovery/export packages.

---

## Phase 10.1 — `forgebe scaffold`

**Outcome:** Generate boilerplate/backend skeleton berdasarkan project profile atau template override.

### Scope
- Supported awal:
  - Go service
  - Go API
  - Node Express
  - Python FastAPI
- Support:
  - infer template dari profile
  - `--template`
  - `--output`
  - `--force`
  - `--dry-run`

### Files
- Create: `internal/scaffold/scaffold.go`
- Create: `internal/scaffold/scaffold_test.go`
- Create: `internal/cli/scaffold.go`
- Create: `internal/cli/scaffold_test.go`
- Modify: `internal/cli/root.go`
- Modify: `docs/cli-reference.md`
- Modify: `README.md`

### TDD Tasks
1. Write failing tests untuk template inference dari profile language/architecture.
2. Write failing tests untuk dry-run menghasilkan daftar file tanpa write disk.
3. Write failing tests untuk actual write, force overwrite, dan existing-file protection.
4. Implement minimal `internal/scaffold` generator hingga test pass.
5. Write failing CLI tests untuk `forgebe scaffold --help`, success path, invalid template, dan dry-run output.
6. Implement CLI command + register ke root.
7. Run: `go test ./...`, `go build -o /tmp/forgebe ./cmd/forgebe`, smoke test `forgebe scaffold --help`.

### Acceptance Criteria
- Bisa generate file tree baru ke output directory.
- Tidak overwrite file existing tanpa `--force`.
- Dry-run tidak membuat file.
- Template bisa diinfer dari profile.

---

## Phase 10.2 — `forgebe adopt`

**Outcome:** Membantu existing project yang belum fully ForgeBE-aware untuk diadopt ke dalam workflow ForgeBE.

### Scope
- Detect current repo path
- Scan existing project
- Create/update profile bila belum ada
- Optional sync adapter/context files ke repo target
- Support `--dry-run`

### Files
- Create: `internal/adopt/adopt.go`
- Create: `internal/adopt/adopt_test.go`
- Create: `internal/cli/adopt.go`
- Create: `internal/cli/adopt_test.go`
- Modify: `internal/cli/root.go`
- Modify: `docs/cli-reference.md`

### TDD Tasks
1. Write failing tests untuk detect existing project without profile.
2. Write failing tests untuk adopt dry-run report.
3. Write failing tests untuk adopt success (create profile + optional sync metadata).
4. Implement adopt service.
5. Add CLI tests + command registration.
6. Regression + smoke test.

### Acceptance Criteria
- Existing repo bisa di-register ke ForgeBE tanpa onboarding manual penuh.
- Output menjelaskan project ID/result clearly.
- Dry-run aman tanpa perubahan disk.

---

## Phase 10.3 — `forgebe verify`

**Outcome:** Validasi perubahan atau codebase terhadap profile/engineering contract.

### Scope
- Check policy alignment:
  - testing strategy
  - dependency policy
  - sensitive areas warning
  - architecture hints
- Support repo path input
- Support JSON output

### Files
- Create: `internal/verify/verify.go`
- Create: `internal/verify/verify_test.go`
- Create: `internal/cli/verify.go`
- Create: `internal/cli/verify_test.go`
- Modify: `internal/cli/root.go`
- Modify: `docs/cli-reference.md`

### TDD Tasks
1. Write failing tests for violation detection.
2. Write failing tests for clean verification report.
3. Implement minimal verifier.
4. Add CLI tests for human-readable and JSON output.
5. Full regression.

### Acceptance Criteria
- Bisa memberi PASS/WARN/FAIL report terhadap repo/profile alignment.
- Menunjukkan violations secara eksplisit dan actionable.

---

## Phase 10.4 — `forgebe check`

**Outcome:** Fast validation command untuk local workflow / pre-commit / CI.

### Scope
- Aggregate:
  - profile presence
  - contract verification
  - sync status
  - optional test/lint hints
- Exit code non-zero when failing
- CI-friendly output and JSON mode

### Files
- Create: `internal/check/check.go`
- Create: `internal/check/check_test.go`
- Create: `internal/cli/check.go`
- Create: `internal/cli/check_test.go`
- Modify: `internal/cli/root.go`
- Modify: `docs/cli-reference.md`

### TDD Tasks
1. Write failing tests for pass/warn/fail exit behavior.
2. Implement checker aggregator.
3. Add CLI tests for output + exit code semantics.
4. Smoke test with real binary.

### Acceptance Criteria
- Bisa dipakai di script/pre-commit/CI.
- Exit behavior konsisten dan terdokumentasi.

---

## Phase 10.5 — `forgebe prompt`

**Outcome:** Generate prompt yang lebih kaya dan context-aware daripada `brief`, untuk task-specific AI execution.

### Scope
- Modes awal:
  - `implement`
  - `review`
  - `debug`
  - `plan`
- Inputs:
  - tool
  - project ID
  - optional task description
- Pull context dari profile + contract + status summary

### Files
- Create: `internal/prompt/prompt.go`
- Create: `internal/prompt/prompt_test.go`
- Create: `internal/cli/prompt.go`
- Create: `internal/cli/prompt_test.go`
- Modify: `internal/cli/root.go`
- Modify: `docs/cli-reference.md`
- Modify: `README.md`

### TDD Tasks
1. Write failing tests for each prompt mode.
2. Implement renderer with mode-specific sections.
3. Add CLI tests for output correctness and errors.
4. Full regression + smoke tests.

### Acceptance Criteria
- Output lebih task-aware daripada `brief`.
- Bisa dipakai langsung ke Claude/Cursor/Hermes/Copilot.

---

## Quality Gate After Each Sub-Phase

Sebelum lanjut ke sub-phase berikutnya, WAJIB:
1. `go fmt ./...`
2. `go vet ./...`
3. `go test ./...`
4. `go build -o /tmp/forgebe ./cmd/forgebe`
5. Minimal 1 smoke test binary untuk command phase itu
6. Update docs (`README.md` / `docs/cli-reference.md`) bila user-facing behavior berubah
7. Commit terpisah per sub-phase

---

## Delivery Order
1. Phase 10.1 `scaffold`
2. Phase 10.2 `adopt`
3. Phase 10.3 `verify`
4. Phase 10.4 `check`
5. Phase 10.5 `prompt`

---

## Notes
- Go version tetap 1.22.
- Semua generated output harus local-first dan git-safe.
- `--dry-run` diprioritaskan untuk command yang membuat perubahan disk.
- Hindari membuat generated AI context files menjadi source of truth di repo.
- Jaga duplikasi tetap rendah agar Sonar tidak regress.
- Bila template/code generation >300 lines, gunakan chunked write discipline.
- Phase dianggap selesai hanya jika artifact benar-benar bisa dijalankan/diverifikasi, bukan sekadar stub.
