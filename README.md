# ForgeBE — AI-Native Backend Blueprint

**Local-first, AI-agnostic backend engineering framework untuk menjaga kualitas dan konsistensi kode di era AI-assisted development.**

[![CI](https://github.com/farid-alfernass/forgebe/actions/workflows/ci.yml/badge.svg)](https://github.com/farid-alfernass/forgebe/actions/workflows/ci.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=farid-alfernass_forgebe&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=farid-alfernass_forgebe)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=farid-alfernass_forgebe&metric=coverage)](https://sonarcloud.io/summary/new_code?id=farid-alfernass_forgebe)

---

## Masalah

Developer backend di era AI menghadapi masalah nyata:

- **Context Drift**: AI mulai memberikan solusi yang melenceng dari standar project
- **Inkonsistensi**: Output AI berbeda-beda tiap sesi, tidak ada standar yang sama
- **Quality Gap**: Kode AI tidak selalu mengikuti best practice, security, dan arsitektur yang sudah ditetapkan
- **Lock-in**: Tools AI spesifik (Cursor, Claude Code, Copilot) — migrasi jadi sulit

**ForgeBE menyelesaikan ini dengan menjadi "The Developer's Brain" — menyimpan konteks project, standar kualitas, dan panduan AI di satu tempat, local-first, dan bisa dipakai oleh tools AI manapun.**

## Konsep Inti

### 🔷 AI-Native Backend Blueprint

ForgeBE bukan sekadar CLI tool. Ini adalah **framework metodologi** yang terdiri dari:

```
┌───────────────────────────────────────────────────┐
│                  THE BLUEPRINT                     │
├───────────────────────────────────────────────────┤
│ 1. Context Layer                                   │
│    → AGENTS.md, CLAUDE.md, .cursorrules            │
│    → Single Source of Truth untuk AI               │
├───────────────────────────────────────────────────┤
│ 2. Architectural Framework                         │
│    → Hexagonal / Clean Architecture                │
│    → Pemisahan domain, application, infrastructure │
├───────────────────────────────────────────────────┤
│ 3. TDD-AI Loop                                     │
│    → Definition → Test-First → Implementation      │
│    → Refactor dengan AI-driven optimization        │
├───────────────────────────────────────────────────┤
│ 4. Quality Gates                                   │
│    → Pre-commit health check                       │
│    → Policy verification                           │
│    → Coverage minimum threshold                    │
│    → SonarCloud Quality Gate                       │
└───────────────────────────────────────────────────┘
```

### 🔷 Local-First & Git-Safe

Semua state ForgeBE disimpan di `~/.forgebe/`, **bukan di repository project**. Ini artinya:

- ✅ Tidak ada polusi di repo project
- ✅ Bisa dipakai di project baru atau project yang sudah berjalan
- ✅ Aman untuk CI/CD — tidak ada perubahan file yang tidak diinginkan
- ✅ Deterministic project ID berdasarkan path repository

### 🔷 AI-Agnostic

ForgeBE tidak terikat ke satu AI tool. Satu profile project bisa digunakan oleh:

| AI Tool | Status |
|---------|--------|
| Claude (Claude Code, Claude.ai) | ✅ Supported |
| Cursor IDE | ✅ Supported |
| GitHub Copilot | ✅ Supported |
| Hermes Agent | ✅ Supported |
| Generic / Prompt-based | ✅ Supported |

Cukup jalankan `forgebe export adapter <tool>` dan dapatkan konfigurasi yang sesuai.

## Quick Start

### Instalasi

```bash
# macOS (Homebrew)
brew tap farid-alfernass/tap
brew install forgebe

# Linux (Binary)
curl -fsSL https://github.com/farid-alfernass/forgebe/releases/latest/download/forgebe_Linux_x86_64.tar.gz | tar xz
sudo mv forgebe /usr/local/bin/

# From source
git clone https://github.com/farid-alfernass/forgebe.git
cd forgebe
make build
sudo make install
```

### Setup Project Baru

```bash
# 1. Inisialisasi project dengan guided onboarding
cd my-project
forgebe init

# 2. Lihat profile project
forgebe profile show

# 3. Export adapter untuk AI tool yang kamu pakai
forgebe export adapter cursor  # atau claude, copilot, hermes

# 4. Sync konteks AI ke root project
forgebe sync
```

### Setup Project Lama (Sudah Berjalan)

```bash
# 1. Scan project yang sudah ada, ForgeBE akan mendeteksi stack dan arsitektur
cd existing-project
forgebe adopt

# 2. Verifikasi profile sudah sesuai dengan engineering policy
forgebe verify

# 3. Siap digunakan
forgebe brief claude
```

## Workflow Harian dengan AI

### TDD-AI Loop

```
┌───────────────────────────────────────────────────┐
│                   TDD-AI LOOP                      │
│                                                     │
│  1. DEFINITION                                      │
│     → Tulis spec/test case dulu                     │
│     → minta AI: "Buat interface + domain model"     │
│                                                     │
│  2. TEST-FIRST (RED)                                │
│     → Minta AI buat unit test dari spec             │
│     → Test harus FAIL dulu (Red phase)              │
│                                                     │
│  3. IMPLEMENTATION (GREEN)                          │
│     → Minta AI buat logic sampe test PASS           │
│                                                     │
│  4. REFACTOR                                        │
│     → Minta AI optimasi kode                        │
│     → Test masih harus PASS setelah refactor        │
└───────────────────────────────────────────────────┘
```

### Menggunakan Prompt Generator

ForgeBE punya prompt generator yang menghasilkan prompt konteks-spesifik untuk berbagai mode pekerjaan:

```bash
# Implementasi fitur baru
forgebe prompt implement "Add user authentication endpoint"

# Code review
forgebe prompt review "Check payment handler for security issues"

# Debugging
forgebe prompt debug "Fix panic in request middleware"

# Planning
forgebe prompt plan "Migrate database from MySQL to PostgreSQL"
```

### Pre-Commit Check

Sebelum commit, jalankan health gate untuk memastikan kualitas:

```bash
forgebe check
```

Ini akan mengecek:
- ✅ Profile freshness (apakah ada perubahan?)
- ✅ Repo exists (path masih valid?)
- ✅ Verify status (policy compliance)
- ✅ Sync status (AI context up-to-date?)

Jika ada yang FAIL, exit code non-zero — CI/CD otomatis menolak.

## Perintah CLI Lengkap

| Command | Deskripsi |
|---------|-----------|
| `forgebe init` | Inisialisasi project baru dengan guided onboarding |
| `forgebe adopt` | Adopsi project yang sudah ada |
| `forgebe scan` | Scan codebase dan deteksi stack/arsitektur |
| `forgebe scaffold` | Generate boilerplate dari template |
| `forgebe verify` | Verifikasi profile terhadap engineering policies |
| `forgebe check` | Health gate untuk pre-commit / CI |
| `forgebe prompt` | Generate context-aware prompt untuk AI tasks |
| `forgebe doctor` | Diagnosa kesehatan instalasi dan profile |
| `forgebe profile show` | Lihat detail project profile |
| `forgebe brief` | Generate compact AI prompt untuk tool tertentu |
| `forgebe export` | Export project context (summary, adapter, bundle) |
| `forgebe import` | Import project bundle |
| `forgebe sync` | Sync AI context files ke project root |
| `forgebe watch` | Auto-sync saat ada perubahan file |
| `forgebe version` | Tampilkan versi |

Untuk detail lengkap setiap command, lihat [CLI Reference](docs/cli-reference.md).

## Integrasi CI/CD

ForgeBE sudah siap diintegrasikan ke pipeline CI manapun. Contoh GitHub Actions:

```yaml
# .github/workflows/ci.yml
jobs:
  quality-gate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test -coverprofile=coverage.out ./...
      - run: forgebe check    # Health gate — FAIL kalau ada masalah
      - run: forgebe verify   # Policy verification
```

ForgeBE juga terintegrasi dengan **SonarCloud** untuk quality gate otomatis:
- Coverage minimum: **80%**
- Duplicated lines: **<3%**
- Security hotspots: **100% reviewed**

## Arsitektur

```
┌────────────────────────────────────────────────────────────┐
│                       CLI Layer                             │
│  init  adopt  scan  scaffold  verify  check  prompt         │
│  doctor  profile  brief  export  import  sync  watch        │
└──────────────────────────┬─────────────────────────────────┘
                           │
┌──────────────────────────▼─────────────────────────────────┐
│                      Core Domain                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐  │
│  │Discovery │  │ Profile  │  │ Adapters │  │  Contract  │  │
│  │(scanner) │  │  (model) │  │(templates)│  │(extractor) │  │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐  │
│  │Scaffold  │  │ Verify   │  │  Check   │  │  Prompt    │  │
│  │(template)│  │(policies)│  │(health)  │  │(generator) │  │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘  │
└──────────────────────────┬─────────────────────────────────┘
                           │
┌──────────────────────────▼─────────────────────────────────┐
│                     Storage Layer                            │
│  ~/.forgebe/                                                │
│    ├── projects/<id>/    # Per-project state                │
│    │     ├── project-profile.yaml                           │
│    │     └── engineering-contract.md                        │
│    ├── exports/           # Generated adapters & bundles    │
│    └── cache/             # Temporary cache                 │
└─────────────────────────────────────────────────────────────┘
```

Untuk detail arsitektur lengkap, lihat [Architecture](docs/architecture.md).

## Project Status

**v1.2.0** — Production-ready. Lihat [CHANGELOG.md](CHANGELOG.md) untuk riwayat rilis.

### Quality Metrics (Terbaru)

| Metric | Value | Target |
|--------|-------|--------|
| Test Coverage | 85.2% | ≥80% ✅ |
| Duplicated Lines | 1.2% | <3% ✅ |
| Sonar Quality Gate | OK | ✅ |

## Dokumentasi

| Dokumen | Deskripsi |
|---------|-----------|
| [Blueprint Guide](docs/blueprint-guide.md) | Deep dive metodologi AI-Native Backend Engineering |
| [Real-World Examples](docs/real-world-examples.md) | Contoh penggunaan di project nyata: Go, Node.js, legacy, multi-team, CI/CD |
| [Architecture](docs/architecture.md) | System design, data flow, extension points |
| [CLI Reference](docs/cli-reference.md) | Complete command reference with examples |
| [Security](docs/security.md) | Threat model, security boundaries |
| [Contributing](CONTRIBUTING.md) | Panduan kontribusi, testing, PR, release process |

## License

MIT © 2026 Farid Al Fernass
