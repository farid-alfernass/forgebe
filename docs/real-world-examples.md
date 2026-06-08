# Real-World Examples

Panduan ini menunjukkan cara menggunakan ForgeBE di berbagai skenario project nyata — dari microservice Go, Node.js API, sampai adopsi project legacy.

---

## Contoh 1: Go Microservice (Payment Gateway)

### Skenario

Tim backend membangun **Payment Gateway** microservice menggunakan Go. Tim ini memakai Claude Code dan Cursor IDE sebagai AI assistant.

### Setup Awal

```bash
cd ~/projects/payment-gateway
forgebe init

# Jawab pertanyaan guided onboarding:
# Language: go
# Architecture: hexagonal
# Testing: tdd
# Sensitive areas: Payment processing, PCI compliance, Database migrations
```

### Generate AI Context

```bash
# Untuk developer yang pakai Cursor
forgebe export adapter cursor
# Output: ~/.forgebe/exports/payment-gateway/.cursorrules

# Untuk developer yang pakai Claude Code
forgebe export adapter claude
# Output: ~/.forgebe/exports/payment-gateway/CLAUDE.md

# Sync ke root project (opsional, jika tim mau commit AI context)
forgebe sync --dry-run
forgebe sync
```

### Workflow TDD-AI Loop

```bash
# 1. Generate prompt untuk fitur baru
forgebe prompt implement "Add refund endpoint with idempotency key"

# Output prompt yang dihasilkan (copy ke AI):
# ---
# Project: payment-gateway
# Language: go
# Architecture: hexagonal
# Sensitive areas: Payment processing, PCI compliance
#
# Task: Implement "Add refund endpoint with idempotency key"
#
# Guidelines:
# - Follow hexagonal architecture (domain → application → infrastructure)
# - Write unit tests FIRST (TDD red phase)
# - Use repository pattern for database access
# - Handle errors with custom error types
# - Consider idempotency for payment operations
# ---

# 2. Setelah AI selesai menulis kode, verifikasi
go test ./...
forgebe check
forgebe verify
```

### Pre-Commit Hook

```bash
# .git/hooks/pre-commit
#!/bin/bash
set -e
gofmt -w ./cmd ./internal
go vet ./...
go test -count=1 ./...
forgebe check --json | jq -e '.result == "PASS"'
```

### CI Pipeline

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go test -coverprofile=coverage.out ./...
      - run: forgebe check
      - run: forgebe verify
```

---

## Contoh 2: Node.js API (E-Commerce Backend)

### Skenario

Tim full-stack membangun **E-Commerce Backend** dengan Node.js + TypeScript + NestJS. Developer menggunakan Cursor IDE.

### Setup Awal

```bash
cd ~/projects/ecommerce-api
forgebe init

# Language: typescript
# Architecture: nestjs (modular)
# Testing: jest
# Sensitive areas: User authentication, Payment integration, PII handling
```

### Scaffold Module Baru

```bash
# Generate boilerplate untuk service baru
forgebe scaffold --template node-nestjs --output ./services/inventory-service

# Hasil:
# services/inventory-service/
# ├── src/
# │   ├── index.ts
# │   ├── routes/
# │   └── middleware/
# ├── package.json
# ├── tsconfig.json
# └── .gitignore
```

### AI Prompt untuk Code Review

```bash
forgebe prompt review "Check order processing module for race conditions"

# Output:
# ---
# Project: ecommerce-api
# Language: typescript
# Architecture: nestjs (modular)
# Sensitive areas: Payment integration
#
# Task: Review "Check order processing module for race conditions"
#
# Review focus:
# - Concurrent access patterns (database locks, optimistic locking)
# - Transaction boundaries
# - Error rollback behavior
# - Test coverage for edge cases
# ---
```

### Monitoring dengan Watch Mode

```bash
# Auto-sync AI context setiap ada perubahan file
forgebe watch --dry-run  # Preview dulu
forgebe watch            # Jalankan daemon
```

---

## Contoh 3: Adopsi Project Legacy (Monolith Java → Microservice)

### Skenario

Tim sedang memecah **monolith Java** menjadi microservice. Project sudah berjalan 3 tahun tanpa standar AI. Mereka ingin mulai menggunakan AI dengan aman.

### Langkah Adopsi

```bash
# 1. Scan project yang sudah ada
cd ~/projects/legacy-monolith
forgebe adopt

# ForgeBE akan:
# - Scan struktur folder, detect Java + Spring Boot
# - Detect testing framework (JUnit 5)
# - Detect CI (Jenkins)
# - Buat profile berdasarkan hasil scan
# - Tanyakan konfirmasi untuk area sensitif

# 2. Review hasil adopsi
forgebe profile show

# Output:
# Project:    legacy-monolith
# Profile ID: legacy-monolith_a3f2b1c4
# Language:   java
# Framework:  spring-boot
# Arch:       layered (controller-service-repository)
# Testing:    junit5
# CI:         jenkins
# Sensitive:  [Database migrations, Authentication, External API contracts]

# 3. Verifikasi apakah project sudah comply dengan policy
forgebe verify

# Output:
# ForgeBE Verify Report
# =====================
# ✅ testing_policy: PASS (junit5 detected)
# ✅ architecture: PASS (layered pattern detected)
# ⚠️  dependency_policy: WARNING (no lock file found)
# ✅ sensitive_areas: PASS (3 areas defined)
# ✅ source_roots: PASS (src/main/java detected)
#
# Result: PASS (4 passed, 1 warning, 0 failures)

# 4. Export brief untuk AI
forgebe brief claude
```

### Workflow: AI-Assisted Refactoring Monolith

```bash
# Generate prompt untuk decompose service
forgebe prompt plan "Extract user-service from monolith: UserController, UserService, UserRepository"

# Generate prompt untuk implementasi
forgebe prompt implement "Create standalone user-service microservice with Spring Boot 3"

# Setelah AI generate kode, jalankan verification
forgebe check
```

---

## Contoh 4: Multi-Team dengan Standar Terpusat

### Skenario

Perusahaan punya 5 tim backend. Setiap tim punya project masing-masing tapi harus mengikuti standar yang sama.

### Setup Template Terpusat

```bash
# Platform engineer membuat template standar
forgebe scaffold --list-templates

# Output:
# Available templates:
#   go-service       - Go microservice (hexagonal)
#   go-api           - Go REST API (clean arch)
#   node-express     - Node.js Express (TypeScript)
#   node-nestjs      - Node.js NestJS (modular)
#   python-fastapi   - Python FastAPI
```

### Setiap Tim Membuat Project Baru

```bash
# Tim A: User Service
forgebe scaffold --template go-api --output ./user-service
cd user-service && forgebe init --non-interactive

# Tim B: Notification Service
forgebe scaffold --template go-service --output ./notification-service
cd notification-service && forgebe init --non-interactive

# Tim C: Analytics API
forgebe scaffold --template python-fastapi --output ./analytics-api
cd analytics-api && forgebe init --non-interactive
```

### Verifikasi Compliance Terpusat

```bash
# Script CI yang dijalankan untuk SEMUA project
for project in user-service notification-service analytics-api; do
  echo "=== Checking $project ==="
  cd ~/projects/$project
  forgebe check --json
  forgebe verify --json
  cd -
done
```

---

## Contoh 5: CI/CD Pipeline Lengkap dengan SonarCloud

### GitHub Actions Workflow Production-Ready

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Go mod tidy check
        run: |
          go mod tidy
          git diff --exit-code

      - name: Format check
        run: |
          gofmt -w ./cmd ./internal
          git diff --exit-code

      - name: Vet
        run: go vet ./...

      - name: Test with coverage
        run: go test -coverprofile=coverage.out ./...

      - name: ForgeBE Health Check
        run: forgebe check

      - name: ForgeBE Policy Verification
        run: forgebe verify

      - name: Build
        run: go build -o bin/forgebe ./cmd/forgebe

      - name: Upload coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage-out
          path: coverage.out

  sonar:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/download-artifact@v4
        with:
          name: coverage-out

      - name: SonarCloud Scan
        uses: SonarSource/sonarcloud-github-action@v3
        env:
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Branch Protection Rules (Wajib)

Setelah CI berjalan, set di GitHub → Settings → Branches → main:

- ✅ Require status checks: `test`, `sonar`
- ✅ Require branches to be up to date
- ✅ Require pull request reviews (minimal 1)
- ✅ Dismiss stale reviews on new push

---

## Tips & Tricks

### 1. Combine dengan Git Hooks

```bash
# Install pre-commit hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
set -e
echo "Running ForgeBE quality gate..."
forgebe check || { echo "❌ ForgeBE check failed"; exit 1; }
echo "✅ Quality gate passed"
EOF
chmod +x .git/hooks/pre-commit
```

### 2. Quick Brief untuk Chat AI

```bash
# Copy brief langsung ke clipboard (macOS)
forgebe brief claude | pbcopy

# Paste ke Claude.ai, ChatGPT, atau AI tool lainnya
```

### 3. Dry-Run untuk Semua Command

```bash
# Semua command yang menulis file mendukung --dry-run
forgebe scaffold --template go-api --output ./test --dry-run
forgebe sync --dry-run
forgebe watch --dry-run
```

### 4. JSON Output untuk Scripting

```bash
# Semua command mendukung --json untuk integrasi scripting
forgebe check --json | jq '.result'
forgebe verify --json | jq '.policies[] | select(.status == "FAIL")'
forgebe scan --json | jq '.detections.languages'
```
