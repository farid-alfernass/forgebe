# Contributing to ForgeBE

Thank you for contributing to ForgeBE — an AI-agnostic, local-first backend engineering framework.

ForgeBE exists to help backend engineers use AI tools without losing consistency, safety, or delivery discipline.

## Principles

All contributions should preserve these core principles:

1. **Local-first**
   - Runtime state belongs in `~/.forgebe/`.
   - ForgeBE must not send project context to remote services by default.

2. **Git-safe**
   - ForgeBE must not write tracked files into a target repository without explicit user intent.
   - Adapter exports go to `~/.forgebe/exports/` unless the user passes `--output`.

3. **AI-agnostic**
   - ForgeBE must not depend on a single AI vendor or model.
   - Tool integrations should be adapters, not hard dependencies.

4. **Project-aware**
   - Existing architecture and conventions must be preserved by default.
   - Refactors and dependency additions should be explicit and reviewable.

5. **Secure by default**
   - No hardcoded secrets.
   - No unsafe archive extraction.
   - No shell execution with untrusted input.
   - File writes must validate destination paths.

## Development Setup

```bash
git clone https://github.com/farid-alfernass/forgebe.git
cd forgebe

go mod download
go test ./...
go build -o bin/forgebe ./cmd/forgebe
```

## Quality Gate

Before opening a PR or pushing a change, run:

```bash
gofmt -w ./cmd ./internal
go vet ./...
go test -count=1 ./...
go build -o bin/forgebe ./cmd/forgebe
```

For CLI-affecting changes, also run an E2E smoke test:

```bash
./bin/forgebe scan .
./bin/forgebe scan --json .
./bin/forgebe doctor
./bin/forgebe profile show
./bin/forgebe brief generic
./bin/forgebe export summary
./bin/forgebe export adapter generic
./bin/forgebe export bundle
```

## Testing Guidance

- Add unit tests for new logic under `internal/...`.
- Prefer deterministic fixture-based tests.
- Keep `testdata/`, `fixtures/`, and `examples/` out of production scanning logic unless explicitly intended.
- For archive import/export, include security regression tests for traversal and malformed paths.

## CLI UX Guidelines

- Commands should have predictable text output for humans.
- Machine-readable commands should support `--json` where useful.
- File-generating commands should support `--output`.
- Errors should include command context, e.g. `export bundle: ...`.

## Commit Style

Use clear, scoped messages:

```text
feat: add adapter export output flag
fix: block symlink traversal during bundle import
docs: add CLI reference
```

## Pull Request Checklist

- [ ] Code is formatted with `gofmt`.
- [ ] `go vet ./...` passes.
- [ ] `go test -count=1 ./...` passes.
- [ ] CLI changes have E2E verification.
- [ ] Security-sensitive changes have explicit tests or review notes.
- [ ] Documentation is updated when behavior changes.

## Release Notes

When preparing a release:

1. Update `CHANGELOG.md`.
2. Run the full quality gate.
3. Tag with semantic versioning:

```bash
git tag vX.Y.Z
git push origin main --tags
```

GitHub Actions will build release binaries for supported platforms.
