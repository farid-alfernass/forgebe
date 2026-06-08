# GitHub Release Guide untuk ForgeBE

Panduan lengkap untuk membuat dan mengelola release di GitHub.

---

## Prerequisites

- GitHub CLI (`gh`) sudah terinstall dan authenticated
- Repository sudah punya CI/CD workflow (GitHub Actions)
- `CHANGELOG.md` sudah diupdate

---

## Langkah-Langkah Release

### 1. Update CHANGELOG.md

```bash
# Edit CHANGELOG.md dan tambahkan versi baru di bagian atas
cat > /tmp/changelog_entry.txt << 'EOF'
## [1.3.0] - 2026-06-08

### Added
- docs: restructure documentation for AI-Native Backend Blueprint
- docs: add real-world examples (5 scenarios)
- docs: update contributing guide with detailed guidelines
- CLI: enhanced scaffold, adopt, verify, check, prompt commands (Phase 10)

### Fixed
- fix: reduce code duplication in scaffold generator
- fix: pin GitHub Actions to full SHA for security

### Changed
- docs: reorganize docs structure for better navigation
- ci: improved quality gate checks

### Quality
- Coverage: 85.2% (target >=80%) ✅
- Duplicated lines: 1.2% (threshold <3%) ✅
- Sonar Quality Gate: OK ✅
EOF

# Update CHANGELOG.md dengan entry baru
```

### 2. Verify All Quality Gates

```bash
cd /Users/raihan/forgebe

# Format
gofmt -w ./cmd ./internal

# Lint
go vet ./...

# Test
go test -count=1 ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep '^total'

# Build
go build -o bin/forgebe ./cmd/forgebe

# Smoke test
./bin/forgebe version
./bin/forgebe doctor
./bin/forgebe check
```

### 3. Create Release Tag

```bash
# Set version
VERSION="1.3.0"

# Create annotated tag (bukan lightweight tag)
git tag -a "v${VERSION}" -m "Release v${VERSION}

This release includes:
- Complete documentation restructuring for AI-Native Backend Blueprint
- Real-world examples (5 scenarios: Go, Node.js, legacy, multi-team, CI/CD)
- Enhanced contributing guidelines with detailed instructions
- Phase 10 feature expansion (scaffold, adopt, verify, check, prompt)
- Code quality improvements (coverage 85.2%, duplication 1.2%)
- Security: pinned GitHub Actions to full SHA

See CHANGELOG.md for full details."

# Verify tag dibuat
git tag -l -n10 | head -20

# Push tag ke remote
git push origin "v${VERSION}"
```

### 4. Create GitHub Release (via gh CLI)

```bash
VERSION="1.3.0"

# Buat release dengan notes
gh release create "v${VERSION}" \
  --title "ForgeBE v${VERSION}" \
  --notes-file /tmp/release_notes.md \
  --draft=false

# Atau lebih interaktif:
gh release create "v${VERSION}" \
  --title "ForgeBE v${VERSION}" \
  --notes "Release v${VERSION} - Complete documentation restructuring and Phase 10 features. See CHANGELOG.md for details."
```

### 5. Verifikasi Release di GitHub

```bash
# View releases
gh release list --limit 5

# View specific release
gh release view "v${VERSION}"

# View release assets (binaries)
gh release view "v${VERSION}" --json assets
```

---

## GitHub Actions Workflow untuk Auto-Release Binaries

File: `.github/workflows/release.yml` (sudah ada di repo)

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Workflow ini otomatis:
- Build binary untuk semua platform (Linux, macOS, Windows)
- Create checksum & sign files
- Upload ke GitHub release page

---

## Checklist Sebelum Release

```bash
#!/bin/bash
set -e

echo "🔍 ForgeBE Release Checklist"
echo "============================"

# 1. Check semua test pass
echo "✓ Running tests..."
go test -count=1 ./...

# 2. Check coverage
echo "✓ Checking coverage..."
go test -coverprofile=coverage.out ./...
COVERAGE=$(go tool cover -func=coverage.out | grep '^total' | awk '{print $3}' | sed 's/%//')
echo "  Coverage: ${COVERAGE}%"
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "  ❌ Coverage below 80%"
    exit 1
fi

# 3. Check git status clean
echo "✓ Checking git status..."
if [ -n "$(git status --porcelain)" ]; then
    echo "  ❌ Working directory not clean"
    exit 1
fi

# 4. Check CHANGELOG updated
echo "✓ Checking CHANGELOG..."
if ! grep -q "^## \[" CHANGELOG.md; then
    echo "  ❌ CHANGELOG not updated"
    exit 1
fi

# 5. Check README updated
echo "✓ Checking README..."
if grep -q "vX.Y.Z" README.md; then
    echo "  ⚠️  README has placeholder version"
fi

echo ""
echo "✅ All checks passed! Ready to release."
echo ""
echo "Next steps:"
echo "  1. git tag -a v1.3.0 -m '...'"
echo "  2. git push origin v1.3.0"
echo "  3. gh release create v1.3.0 --notes-file /tmp/notes.md"
```

Jalankan:
```bash
chmod +x .github/scripts/pre-release-check.sh
.github/scripts/pre-release-check.sh
```

---

## Download Release Binary

Setelah release dibuat, binary bisa di-download dari beberapa cara:

### Cara 1: GitHub CLI

```bash
# List releases
gh release list

# Download binary untuk macOS
gh release download v1.3.0 -p '*Darwin*'

# Download binary untuk Linux
gh release download v1.3.0 -p '*Linux*'

# Download binary untuk Windows
gh release download v1.3.0 -p '*Windows*'
```

### Cara 2: Direct URL

```bash
# Format: https://github.com/owner/repo/releases/download/tag/asset-name

# macOS
curl -L https://github.com/farid-alfernass/forgebe/releases/download/v1.3.0/forgebe_Darwin_x86_64.tar.gz

# Linux
curl -L https://github.com/farid-alfernass/forgebe/releases/download/v1.3.0/forgebe_Linux_x86_64.tar.gz

# Windows
curl -L https://github.com/farid-alfernass/forgebe/releases/download/v1.3.0/forgebe_Windows_x86_64.zip
```

### Cara 3: Homebrew (setelah setup tap)

```bash
brew tap farid-alfernass/tap
brew install forgebe@1.3.0
```

---

## Post-Release Tasks

### 1. Update Documentation Links

- [ ] Update README.md version badges
- [ ] Update quick start section if changed
- [ ] Update docs/cli-reference.md if new commands added

### 2. Announce Release

- [ ] Post di GitHub Discussions
- [ ] Tweet/LinkedIn announcement
- [ ] Share di internal Slack/Teams channel
- [ ] Update team wiki/docs

### 3. Monitoring

```bash
# Monitor download stats
gh release view v1.3.0 --json assets --jq '.assets[] | {name, downloadCount}'

# View releases over time
gh release list --limit 10 --json tagName,publishedAt,author
```

---

## Troubleshooting

### Release tidak ter-trigger di GitHub Actions

Pastikan:
1. Tag dibuat dengan format `v*` (misal: `v1.3.0`, BUKAN `1.3.0`)
2. Workflow file `.github/workflows/release.yml` ada di repo
3. GitHub Token punya permission `contents: write`

### Binary tidak ter-build di GitHub Actions

Cek logs:
```bash
gh run list --workflow=release.yml
gh run view <run-id> --log
```

Pastikan:
1. Go version di workflow match dengan go.mod
2. Dependencies bisa ter-download (no network issues)
3. No hard-coded paths di build script

### Checksum tidak match

Re-download dan verify:
```bash
# Download
curl -O https://github.com/farid-alfernass/forgebe/releases/download/v1.3.0/forgebe_Darwin_x86_64.tar.gz

# Download checksum
curl -O https://github.com/farid-alfernass/forgebe/releases/download/v1.3.0/checksums.txt

# Verify (macOS)
shasum -a 256 -c checksums.txt

# Verify (Linux)
sha256sum -c checksums.txt
```

---

## Release Timeline

Rekomendasi untuk fase rilis:

| Waktu | Action |
|------|--------|
| **T-1 day** | Create release branch, update CHANGELOG |
| **T-0 day** | Run pre-release checklist, create tag, push |
| **T+0** | GitHub Actions builds binaries (5-10 menit) |
| **T+15 min** | Verify binaries available, test download |
| **T+30 min** | Announce release di public channel |

---

## Automated Release (Optional Advanced Setup)

Untuk super-automated release, bisa setup:

```bash
# .github/workflows/auto-release.yml
name: Auto Release
on:
  push:
    branches: [main]

jobs:
  release:
    if: contains(github.event.head_commit.message, '[release]')
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: |
          VERSION=$(grep '^## \[' CHANGELOG.md | head -1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
          git tag "v${VERSION}"
          git push origin "v${VERSION}"
```

Cara pakai:
```bash
git commit -m "feat: add new feature [release]"
git push origin main
# GitHub Actions otomatis detect [release] dan create release
```

---

## Reference

- [GitHub CLI Release Docs](https://cli.github.com/manual/gh_release_create)
- [GoReleaser Docs](https://goreleaser.com/)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
