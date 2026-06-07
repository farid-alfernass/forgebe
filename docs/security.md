# ForgeBE Security Model

## Overview

ForgeBE is designed with security as a core principle. This document describes the security boundaries, threat model, and mitigations implemented in ForgeBE.

## Security Principles

1. **Local-first by default** — No data leaves the user's machine unless explicitly exported.
2. **Git-safe** — No writes to tracked repository files without explicit user intent.
3. **Validate all inputs** — User input, file paths, and archive contents are validated.
4. **Fail secure** — Errors in security-critical operations result in safe failures, not silent degradation.

---

## Threat Model

### In Scope

| Threat | Mitigation |
|--------|------------|
| **Path Traversal (ZipSlip)** | Archive imports validate all extracted paths against a resolved, cleaned target directory. |
| **Symlink Escape** | Symlink components in extraction paths are resolved and checked against the target root. |
| **Data Corruption** | Atomic writes (temp file + rename) ensure no partial state on crash. |
| **Secrets Exposure** | ForgeBE does not read or store secrets from environment or credential files. |
| **Injection via AI Files** | Contract extractor treats markdown content as data, not instructions. |

### Out of Scope

| Threat | Reason |
|--------|--------|
| **Compromised local machine** | ForgeBE assumes the local environment is trusted. |
| **Network attacks** | ForgeBE does not open network ports or make outbound requests by default. |
| **Malicious Go dependencies** | Standard Go module security applies; use `go mod verify`. |

---

## Security Boundaries

### 1. Filesystem Operations

**Profile Storage**
- All state is stored under `~/.forgebe/`.
- Project directories are keyed by deterministic IDs, not user-controlled names.
- File writes use `AtomicWrite` (write to temp, rename) to prevent corruption.

**Archive Extraction**
```go
// ZipSlip protection
absPath := filepath.Join(targetDir, entry.Name)
if !strings.HasPrefix(filepath.Clean(absPath), targetDir) {
    return error // reject traversal
}

// Symlink protection
realPath, _ := filepath.EvalSymlinks(filepath.Dir(absPath))
if !strings.HasPrefix(realPath, targetDir) {
    return error // reject symlink escape
}
```

### 2. Input Validation

**CLI Arguments**
- Paths are resolved to absolute form with `filepath.Abs()`.
- Project IDs are validated as filesystem-safe strings.

**Archive Contents**
- Entry paths are checked for `..` components and null bytes.
- Extraction is aborted on any path violation.

**Discovery Scanner**
- Skips `testdata/`, `fixtures/`, `examples/` to prevent fixture-based detection errors.
- Reads file structure only, not arbitrary file contents.

### 3. Secrets Handling

**ForgeBE does NOT:**
- Read `.env` files, credential files, or secret stores.
- Store API keys, tokens, or passwords in profiles.
- Transmit any project context over the network.

**Users should:**
- Keep secrets in secure credential stores.
- Never commit `.forgebe/exports/` or bundles to public repositories if they contain sensitive context.

---

## Secure Usage Guidelines

### For Users

1. **Review exports before sharing**
   - `forgebe export summary` and `forgebe export adapter` produce files intended for AI context.
   - Ensure they don't contain sensitive business logic or internal details you don't want shared.

2. **Treat bundles as sensitive**
   - A `.forgebe.zip` bundle contains the full project profile and context.
   - Share only with trusted collaborators.

3. **Don't commit ForgeBE state to public repos**
   - `~/.forgebe/` is local state.
   - If you must version it, use a private repository or encrypted storage.

### For Contributors

1. **No hardcoded secrets**
   - Never commit API keys, tokens, or passwords.
   - Use environment variables for test fixtures if needed.

2. **Validate all file paths**
   - When writing files, ensure the destination is within an expected root.
   - Use `filepath.Clean()` and prefix checks.

3. **Test security regressions**
   - Add tests for path traversal, symlink escape, and malformed inputs.
   - Run `go test -count=1 ./...` before merging.

---

## Security Audit Checklist

Before each major release, verify:

- [ ] No hardcoded secrets in codebase.
- [ ] `go vet ./...` passes without warnings.
- [ ] Archive import tests cover ZipSlip and symlink cases.
- [ ] All file writes use `AtomicWrite` or equivalent safe patterns.
- [ ] No shell execution with user-controlled input.
- [ ] Discovery scanner skips test fixtures correctly.
- [ ] CLI error messages don't leak sensitive paths or data.

---

## Reporting Security Issues

If you discover a security vulnerability in ForgeBE, please report it responsibly:

1. **Do not** open a public issue.
2. Email the maintainer directly or use GitHub Security Advisories.
3. Include:
   - A description of the vulnerability.
   - Steps to reproduce.
   - Affected versions.
   - Suggested fix (if any).

Security reports will be acknowledged within 48 hours and addressed with priority.

---

## Security History

| Version | Date | Issue | Status |
|---------|------|-------|--------|
| v0.5.0 | 2026-06-07 | Added symlink traversal hardening | Fixed |
| v0.5.0 | 2026-06-07 | Added explicit close error checking in bundle import | Fixed |

---

## References

- [OWASP Path Traversal](https://owasp.org/www-community/attacks/Path_Traversal)
- [CWE-22: Improper Limitation of a Pathname](https://cwe.mitre.org/data/definitions/22.html)
- [Go Security Policy](https://go.dev/security)
