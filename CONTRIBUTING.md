# Contributing to ForgeBE

Terima kasih sudah tertarik berkontribusi ke ForgeBE — framework AI-native backend engineering yang local-first dan AI-agnostic.

---

## Prinsip Utama

Semua kontribusi harus menjaga prinsip ini:

1. **Local-first** — State disimpan di `~/.forgebe/`, bukan remote service.
2. **Git-safe** — Tidak menulis file ke repo target tanpa izin eksplisit user.
3. **AI-agnostic** — Tidak bergantung pada satu vendor AI tertentu.
4. **Project-aware** — Menghormati arsitektur dan konvensi project yang sudah ada.
5. **Secure by default** — Tidak ada hardcoded secrets, tidak ada shell execution dengan input untrusted.

---

## Cara Berkontribusi

### Melaporkan Bug

1. Cek apakah bug sudah dilaporkan di [Issues](https://github.com/farid-alfernass/forgebe/issues).
2. Jika belum, buat issue baru dengan template:
   - **Describe the bug**: Apa yang terjadi?
   - **To reproduce**: Langkah-langkah untuk mereproduksi
   - **Expected behavior**: Apa yang seharusnya terjadi?
   - **Environment**: OS, Go version, ForgeBE version
   - **Logs/Output**: Error message atau output yang relevan

### Mengusulkan Fitur

1. Buat issue dengan label `enhancement`.
2. Jelaskan:
   - **Problem**: Masalah apa yang ingin diselesaikan?
   - **Proposed solution**: Solusi yang diusulkan
   - **Alternatives**: Alternatif yang sudah dipertimbangkan
   - **Impact**: Siapa yang terpengaruh?

### Mengirim Pull Request

1. Fork repository
2. Buat branch dari `main`:
   ```bash
   git checkout -b feat/my-feature main
   ```
3. Lakukan perubahan (ikuti panduan di bawah)
4. Pastikan quality gate lokal PASS
5. Push dan buat PR

---

## Development Setup

### Prerequisites

- Go 1.22+
- Git
- Make (opsional, untuk shortcut commands)

### Clone & Build

```bash
git clone https://github.com/farid-alfernass/forgebe.git
cd forgebe

go mod download
go test ./...
go build -o bin/forgebe ./cmd/forgebe
```

### Verify Instalasi

```bash
./bin/forgebe version
./bin/forgebe doctor
```

---

## Struktur Project

```
forgebe/
├── cmd/forgebe/          # CLI entrypoint
├── internal/
│   ├── cli/             # Cobra command definitions
│   ├── discovery/       # Codebase scanning & detection
│   ├── profile/         # Project profile model & store
│   ├── contract/        # Engineering contract renderer
│   ├── adapters/        # AI tool adapters (claude, cursor, etc.)
│   ├── export/          # Export utilities (summary, bundle)
│   ├── onboarding/      # Guided init flow
│   ├── scaffold/        # Project scaffolding / template generator
│   ├── adopt/           # Adopt existing project
│   ├── verify/          # Policy verification engine
│   ├── check/           # Health check / quality gate
│   ├── prompt/          # AI prompt generator
│   └── storage/         # Local storage primitives
├── docs/                # Dokumentasi
├── testdata/            # Test fixtures
└── .github/workflows/   # CI/CD
```

### Konvensi Penting

- **Domain logic** di `internal/<package>/` — setiap package punya tanggung jawab tunggal.
- **CLI wiring** di `internal/cli/` — Cobra command hanya memanggil domain logic.
- **Test** di file `*_test.go` sejajar dengan source.
- **Fixtures** di `testdata/` — excluded dari scanning.

---

## Quality Gate (Wajib Sebelum PR)

Jalankan semua ini sebelum push:

```bash
# Format
gofmt -w ./cmd ./internal

# Lint
go vet ./...

# Test (harus semua PASS)
go test -count=1 ./...

# Build (harus sukses)
go build -o bin/forgebe ./cmd/forgebe

# Coverage check (target: >=80%)
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep '^total'
```

### Untuk Perubahan CLI

Tambahan E2E smoke test:

```bash
./bin/forgebe scan .
./bin/forgebe doctor
./bin/forgebe profile show
./bin/forgebe brief generic
./bin/forgebe check
./bin/forgebe verify
./bin/forgebe scaffold --list-templates
```

---

## Panduan Menulis Kode

### Menambahkan Command CLI Baru

1. Buat file `internal/cli/newcmd.go`:
   ```go
   func newNewcmdCmd() *cobra.Command {
       cmd := &cobra.Command{
           Use:   "newcmd [args]",
           Short: "Deskripsi singkat",
           RunE: func(cmd *cobra.Command, args []string) error {
               // Logic di sini
               return nil
           },
       }
       cmd.Flags().Bool("json", false, "Output JSON")
       return cmd
   }
   ```

2. Register di `root.go`:
   ```go
   cmd.AddCommand(newNewcmdCmd())
   ```

3. Buat test di `internal/cli/newcmd_test.go`
4. Update `docs/cli-reference.md`

### Menambahkan Adapter AI Baru

1. Edit `internal/adapters/adapter.go`:
   ```go
   adapterRegistry["newtool"] = adapterSpec{
       Filename: "NEWTOOL_CONTEXT.md",
       Template: `Project: {{.Project.Name}}
   Language: {{.Stack.PrimaryLanguage}}
   Architecture: {{.Stack.Architecture}}`,
   }
   ```

2. Test: `./bin/forgebe export adapter newtool <project-id>`
3. Tambahkan ke daftar supported tools di `README.md`

### Menambahkan Policy Verification Baru

1. Buat validator di `internal/verify/`:
   ```go
   type MyPolicyValidator struct{}

   func (v *MyPolicyValidator) Name() string { return "my_policy" }
   func (v *MyPolicyValidator) Validate(p *profile.ProjectProfile) Result {
       if condition {
           return Result{Status: StatusPass}
       }
       return Result{Status: StatusFail, Message: "alasan gagal"}
   }
   ```

2. Register di `verify.go`
3. Buat test di `verify_test.go`

### Menambahkan Template Scaffold Baru

1. Tambah constant di `internal/scaffold/scaffold.go`:
   ```go
   const TemplateNewStack TemplateType = "new-stack"
   ```

2. Implementasi generator:
   ```go
   func (g *Generator) generateNewStack() []File {
       projectName := sanitizeName(g.profile.Project.Name)
       return []File{
           {Path: "src", IsDir: true},
           {Path: "src/main.ext", Content: g.newStackMain()},
           {Path: pathGitignore, Content: newStackGitignore},
           {Path: pathReadme, Content: g.readme(projectName)},
       }
   }
   ```

3. Tambahkan ke switch di `Generate()` dan ke `Templates()`
4. Buat test

---

## Testing Guidelines

### Unit Test

- Setiap file `xxx.go` harus punya `xxx_test.go`.
- Gunakan `t.TempDir()` untuk isolasi filesystem.
- Gunakan `t.Setenv()` untuk environment variable testing.
- Prefer table-driven tests:
  ```go
  tests := []struct {
      name     string
      input    string
      expected string
  }{
      {"basic", "input", "expected"},
      {"edge case", "", "default"},
  }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          // test logic
      })
  }
  ```

### Integration Test (CLI)

- Gunakan helper `newRootCmdWithProject(t, args)` untuk test CLI commands.
- Test minimal: help output, flag parsing, basic execution.
- Gunakan `--dry-run` flag untuk test tanpa side effects.

### Coverage Target

- Total coverage: **≥80%**
- Package baru: **≥85%**
- Setiap PR harus menjaga coverage tidak turun

---

## Commit Style

Gunakan [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new scaffold template for Rust
fix: prevent panic when profile is nil
docs: update CLI reference with verify examples
test: add integration tests for adopt command
refactor: extract constants from scaffold generator
ci: pin GitHub Actions to full SHA
```

### Scope (opsional)

```
feat(scaffold): add java-spring template
fix(verify): handle missing testing policy gracefully
test(cli): strengthen Phase 10 coverage
```

---

## Pull Request Checklist

Sebelum submit PR, pastikan:

- [ ] Code formatted dengan `gofmt`
- [ ] `go vet ./...` pass
- [ ] `go test -count=1 ./...` pass (semua test hijau)
- [ ] Coverage ≥80% (tidak turun dari sebelumnya)
- [ ] CLI perubahan sudah di-smoke test
- [ ] Dokumentasi diupdate jika ada perubahan behavior
- [ ] Security-sensitive changes punya explicit test
- [ ] Commit messages mengikuti conventional commits
- [ ] PR description menjelaskan: apa yang berubah, kenapa, dan bagaimana di-test

---

## Release Process

### Versioning

ForgeBE menggunakan [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes di CLI interface atau profile format
- **MINOR**: Fitur baru yang backward-compatible
- **PATCH**: Bug fix dan improvement kecil

### Steps

1. Update `CHANGELOG.md`
2. Jalankan full quality gate
3. Tag dan push:
   ```bash
   git tag v1.3.0
   git push origin main --tags
   ```
4. GitHub Actions akan otomatis build release binary untuk semua platform

---

## Bantuan

- **Issues**: https://github.com/farid-alfernass/forgebe/issues
- **Discussions**: https://github.com/farid-alfernass/forgebe/discussions

Terima kasih sudah membantu membuat ForgeBE lebih baik! 🙏
