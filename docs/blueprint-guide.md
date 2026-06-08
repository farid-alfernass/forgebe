# Blueprint Guide — AI-Native Backend Engineering

## Pendahuluan

Ini adalah panduan mendalam tentang **AI-Native Backend Blueprint**, sebuah metodologi dan framework yang dirancang untuk menjaga kualitas, konsistensi, dan keamanan kode di era AI-assisted development.

### Masalah yang Diselesaikan

| Masalah | Dampak | Solusi ForgeBE |
|---------|--------|----------------|
| **Context Drift** | AI mulai memberikan solusi yang melenceng dari standar project | Single Source of Truth via project profile + AI context files |
| **Inkonsistensi Output** | Tiap sesi AI hasilnya berbeda-beda | TDD-AI Loop yang terstruktur |
| **Quality Gap** | Kode AI tidak selalu mengikuti best practice | Quality Gates + Policy Verification |
| **Tool Lock-in** | Migrasi antar AI tool susah | AI-agnostic adapter system |
| **No Traceability** | Tidak tahu standar apa yang dipakai | `verify`, `check`, `prompt` — semuanya teraudit |

---

## 1. The Context Layer

### Filosofi

Seorang developer senior tidak menulis ulang semua aturan setiap kali mulai coding. Begitu juga AI. Kita perlu **Single Source of Truth** yang dibaca AI sebelum mulai bekerja.

### Implementasi

ForgeBE mengelola 3 file konteks utama:

#### `AGENTS.md`

Definisi peran AI dalam project. Apa yang boleh dan tidak boleh dilakukan.

```markdown
# Hermes Agent Context

Project: my-service
Language: go
Architecture: hexagonal
Delivery mode: hybrid
Priority: safety
Sensitive areas: Authentication, Database migrations, Public API contracts
```

#### `CLAUDE.md`

Aturan teknis yang harus diikuti AI.

```markdown
# Backend Engineering Guideline

## Architecture
- Gunakan Hexagonal Architecture
- Domain layer TIDAK boleh import infrastructure layer

## Code Style
- Error handling: selalu return custom error objects
- Setiap public function harus punya GoDoc

## Testing
- Coverage minimal 80%
- Setiap use case harus punya unit test
```

#### `.cursorrules`

Instruksi spesifik untuk Cursor IDE (auto-complete, chat behavior).

### Cara Kerja

```bash
# Sync semua AI context files ke root project
forgebe sync

# Atau auto-sync saat coding
forgebe watch
```

---

## 2. TDD-AI Loop

Ini adalah workflow inti dari blueprint ini. Tujuannya: **jangan biarkan AI menulis kode besar sekaligus tanpa pengawasan.**

### Flow Diagram

```
┌──────────────────────────────────────────────────────────┐
│                     TDD-AI LOOP                           │
├──────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐                                         │
│  │  DEFINITION  │  → Tulis spec / interface               │
│  │  (Stage 1)   │  → Minta AI buat domain model           │
│  └──────┬───────┘                                         │
│         ▼                                                 │
│  ┌──────────────┐                                         │
│  │  TEST-FIRST  │  → Minta AI buat unit test dari spec    │
│  │  (Stage 2)   │  → **Test harus FAIL** (RED phase)      │
│  └──────┬───────┘                                         │
│         ▼                                                 │
│  ┌──────────────┐                                         │
│  │ IMPLEMENT.   │  → Minta AI buat logic                  │
│  │  (Stage 3)   │  → **Test harus PASS** (GREEN phase)    │
│  └──────┬───────┘                                         │
│         ▼                                                 │
│  ┌──────────────┐                                         │
│  │  REFACTOR    │  → Minta AI optimasi (Big O, memory)    │
│  │  (Stage 4)   │  → **Test masih PASS**                   │
│  └──────────────┘                                         │
│                                                           │
└──────────────────────────────────────────────────────────┘
```

### Contoh Penggunaan dengan AI Prompt

**Stage 1 (Definition):**
```
forgebe prompt implement "Add user subscription feature"
```

AI akan menerima:
- Struktur project saat ini
- Standar arsitektur
- Aturan penulisan kode
- Sensitive areas

**Stage 2 (Test-First):**
```
Buat unit test untuk Subscribe method dengan aturan:
- Gunakan testify
- Coverage untuk success case, expired, payment failed
- Jangan gunakan mock untuk domain logic
```

**Stage 3 (Implementation):**
```
Implementasi Subscribe method sampai semua test PASS.
Gunakan repository pattern untuk database access.
```

**Stage 4 (Refactor):**
```
Review kode yang sudah ditulis:
- Apakah ada concurrent access issues?
- Apakah error handling sudah sesuai standar?
- Optimasi database query (N+1 problem?)
```

---

## 3. Quality Gates

### Pre-Commit Gate

Sebelum commit, jalankan:

```bash
forgebe check
```

Output:
```
ForgeBE Check
=============

✅ profile_freshness: PASS
✅ repo_exists: PASS
✅ sync_status: PASS
✅ verify: PASS

Result: PASS (4 passed, 0 warnings, 0 failures)
```

Jika ada masalah, exit code non-zero → commit dibatalkan otomatis oleh git hooks atau CI.

### Policy Verification

```bash
forgebe verify
```

Memvalidasi project profile terhadap **8 engineering policies**:

| Policy | Deskripsi |
|--------|-----------|
| testing_policy | Apakah testing strategy sudah didefinisikan? |
| dependency_policy | Apakah dependency management terkontrol? |
| change_policy | Apakah ada aturan perubahan kode? |
| sensitive_areas | Apakah area sensitif sudah ditandai? |
| source_roots | Apakah direktori source sudah terdefinisi? |
| architecture | Apakah pola arsitektur sudah dipilih? |
| forbidden_paths | Apakah ada path yang dilarang? |

### CI Quality Gates (SonarCloud)

ForgeBE terintegrasi dengan SonarCloud untuk quality gate otomatis:

```yaml
# Dari .github/workflows/ci.yml
- name: SonarCloud Scan
  uses: SonarSource/sonarcloud-github-action@v3
  env:
    SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Standar minimum yang diterapkan:

| Metric | Threshold | Ketika Gagal |
|--------|-----------|--------------|
| Coverage | ≥80% | PR tidak bisa di-merge |
| Duplicated Lines | <3% | Wajib refactor |
| Security Hotspots | 100% reviewed | Harus review manual |
| Reliability Rating | A (1) | Bug wajib diperbaiki |
| Security Rating | A (1) | Vulnerability wajib diperbaiki |

### Setup Branch Protection

Agar quality gates menjadi blocker untuk merge PR:

1. GitHub → Settings → Branches → Add rule
2. Set branch pattern: `main`
3. Centang: **Require status checks to pass before merging**
4. Pilih status checks:
   - `test` (dari CI workflow)
   - `sonar` (dari SonarCloud)
5. Centang: **Require branches to be up to date**

---

## 4. Architectural Framework

### Clean Architecture / Hexagonal

ForgeBE mendorong penggunaan pola yang sangat terstruktur karena AI bekerja paling baik dengan **pemisahan tanggung jawab yang jelas**.

```
src/
├── domain/           # Business logic murni
│   └── entity.go     # → AI: Tidak import database
│
├── application/      # Use Cases / Service layer
│   └── service.go    # → AI: Orchestrasi domain + infrastructure
│
├── infrastructure/   # Database, External API, Redis
│   └── repository.go # → AI: Hanya implementasi interface
│
└── interface/        # Controller, Message Broker, CLI
    └── handler.go    # → AI: Validasi input, format output
```

### Kenapa Ini Penting untuk AI?

Ketika kamu bilang ke AI:
> "Buat use case baru di folder **application**, panggil interface dari **domain**, jangan sentuh **infrastructure**."

AI dengan arsitektur yang jelas bisa mengerjakan tugas dengan presisi tinggi karena lingkup perubahannya terbatas dan terdefinisi.

---

## 5. Project Lifecycle dengan ForgeBE

### A. Project Baru

```bash
forgebe init                     # Setup profile + konfigurasi AI
forgebe scaffold --template go-api --output ./new-service  # Generate boilerplate
forgebe sync                     # Sync AI context
```

### B. Adopsi Project Lama

```bash
forgebe adopt /path/to/existing    # Scan + buat profile
forgebe verify                      # Validasi policy compliance
forgebe export adapter cursor       # Siapkan konfigurasi untuk Cursor
```

### C. Development Harian

```bash
forgebe prompt implement "Add endpoint"  # Generate AI prompt
# ... coding dengan AI ...
go test ./...                             # Pastikan test PASS
forgebe check                             # Health gate
git commit -m "feat: add endpoint"        # Commit (dengan pre-commit hook)
```

### D. Code Review

```bash
forgebe prompt review "Check auth handler for security"
# Dapatkan prompt review yang sudah terisi konteks
```

### E. CI/CD Pipeline

```yaml
name: CI
on: [push, pull_request]
jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test -coverprofile=coverage.out ./...
      - run: forgebe check
      - run: forgebe verify
      - uses: SonarSource/sonarcloud-github-action@v3
```

---

## 6. Security & Compliance

### Sensitive Areas

ForgeBE memungkinkan kamu menandai area kode yang sensitif:

```yaml
# project-profile.yaml
sensitive_areas:
  - Authentication / Authorization
  - Database migrations
  - Public API contracts
  - Payment processing
```

Ketika AI menerima konteks ini (via `brief` atau `prompt`), AI akan ekstra hati-hati di area tersebut.

### Engineering Contract

ForgeBE menghasilkan **engineering contract** — dokumen yang berisi semua aturan dan standar project. Ini bisa dipakai sebagai referensi untuk:

- Onboarding developer baru
- Audit kualitas
- Review AI output

---

## 7. Extension Points

### Menambahkan AI Tool Adapter Baru

Setiap adapter adalah template Go `text/template`:

```go
adapterRegistry["newtool"] = adapterSpec{
    Filename: "NEWTOOL_CONTEXT.md",
    Template: `Project: {{.Project.Name}}
Language: {{.Stack.PrimaryLanguage}}
Sensitive: {{range .SensitiveAreas}}{{.}} {{end}}`,
}
```

### Menambahkan Policy Verifikasi Baru

Buat validator baru di `internal/verify/`:

```go
type MyPolicyValidator struct{}

func (v *MyPolicyValidator) Name() string { return "my_policy" }
func (v *MyPolicyValidator) Validate(proj *profile.ProjectProfile) verify.Result {
    // Logika validasi
}
```

### Menambahkan Template Scaffold Baru

Tambah case baru di `internal/scaffold/scaffold.go`:

```go
case TemplateJavaSpring:
    files = g.generateJavaSpring()
```

---

## 8. Best Practices

### Untuk Developer

1. **Always verify**: Jangan pernah percaya output AI mentah-mentah. Selalu jalankan `forgebe verify` dan `forgebe check`.
2. **TDD-AI Loop wajib**: Tidak boleh lompat langsung ke implementasi tanpa test case dulu.
3. **Sync sebelum coding**: Jalankan `forgebe sync` di awal sesi coding untuk memastikan AI punya konteks terbaru.

### Untuk Tim

1. **Standardisasi**: Semua project dalam tim harus pakai ForgeBE dengan konfigurasi yang konsisten.
2. **Branch Protection**: Quality gates sebagai syarat merge PR.
3. **Audit trail**: Semua perubahan AI-driven harus ter-trace via ForgeBE profile.

### Untuk Platform Engineer

1. **Internal template**: Buat template scaffold khusus yang sesuai dengan standar perusahaan.
2. **Custom policies**: Tambahkan policy verifikasi yang spesifik untuk domain bisnis.
3. **Self-hosted SonarCloud**: Jika perlu data governance yang lebih ketat.

---

## Referensi

- [Architecture](architecture.md) — Detail teknis arsitektur ForgeBE
- [CLI Reference](cli-reference.md) — Semua perintah lengkap dengan contoh
- [Security](security.md) — Threat model dan security guidelines
- [Contributing](../CONTRIBUTING.md) — Panduan kontribusi
