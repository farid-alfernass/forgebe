# ForgeBE

**Open-Source AI-Agnostic Backend Engineering Framework**

ForgeBE is a local-first, AI-agnostic backend engineering framework designed to keep AI-assisted software delivery consistent, safe, and project-aware across any codebase — whether new or existing.

[![CI](https://github.com/farid-alfernass/forgebe/actions/workflows/ci.yml/badge.svg)](https://github.com/farid-alfernass/forgebe/actions/workflows/ci.yml)

---

## Why ForgeBE?

AI tools like Claude, Cursor, Copilot, GPT, Gemini, Codex, Hermes, and local LLMs are powerful. But without explicit guardrails, they produce inconsistent code, ignore project conventions, introduce vulnerabilities, and drift from your engineering standards.

ForgeBE bridges this gap by acting as a **context orchestrator and adapter generator**:

- Defines a canonical **Project Profile** once
- Exports **tool-specific adapters** for any AI assistant
- Validates consistency between your profile and actual project state
- Bundles portable context for team handoff or cross-device sync

---

## 🔧 Installation

### macOS (Homebrew)
```bash
# Coming soon
brew tap farid-alfernass/forgebe
brew install forgebe
```

### Pre-built binary
Download the latest binary from the [Releases](https://github.com/farid-alfernass/forgebe/releases) page.

```bash
# Linux / macOS
chmod +x forgebe-linux-amd64
sudo mv forgebe-linux-amd64 /usr/local/bin/forgebe
```

### From source
```bash
git clone https://github.com/farid-alfernass/forgebe.git
cd forgebe
make build
sudo mv bin/forgebe /usr/local/bin/
```

---

## 🚀 Quickstart

### 1. Initialize a project
```bash
cd /path/to/your/project
forgebe init --non-interactive .
```
ForgeBE scans your codebase (detects Go, Node, Python, Java, Rust), infers architecture, testing framework, and CI setup — then stores a canonical profile locally.

### 2. View your project profile
```bash
forgebe profile show
```

### 3. Validate health
```bash
forgebe doctor
```

### 4. Generate a brief for any AI tool
```bash
forgebe brief claude
forgebe brief cursor
forgebe brief copilot
```

### 5. Export tool-specific adapter files
```bash
forgebe export adapter claude      # → CLAUDE.md
forgebe export adapter cursor      # → .cursorrules
forgebe export adapter copilot     # → .github/copilot-instructions.md
```

### 6. Share context with a teammate
```bash
forgebe export bundle
# Transfer the .forgebe.zip file to another machine
forgebe import bundle ~/Downloads/project.forgebe.zip
```

---

## 📖 Command Reference

| Command | Description |
|---------|-------------|
| `forgebe init` | Interactive or auto-discovery onboarding |
| `forgebe scan [path]` | Discover language, architecture, testing, CI |
| `forgebe doctor [id]` | Validate project profile and local state |
| `forgebe profile show` | Display the active project profile |
| `forgebe brief <tool>` | Compact prompt ready to paste into AI chat |
| `forgebe export summary [id]` | Full markdown project summary |
| `forgebe export adapter <tool> [id]` | Write tool-specific adapter file |
| `forgebe export bundle [id] [path]` | Portable .forgebe.zip bundle |
| `forgebe import <bundle> [dir]` | Restore a bundle into local storage |
| `forgebe version` | Show version |

### Supported AI adapters
- **claude** → `CLAUDE.md`
- **cursor** → `.cursorrules`
- **copilot** → `.github/copilot-instructions.md`
- **hermes** → `AGENTS.md`
- **generic** → `AI_CONTEXT.md`

---

## 🏗 Architecture

```
~/.forgebe/
  └── projects/
      └── <project-id>/
          ├── project-profile.yaml     # Canonical project context
          ├── discovery-report.yaml    # Scanner output
          ├── engineering-contract.md  # Human-readable rules
          └── summary.md              # Compact markdown summary
  └── exports/                        # Adapter files and bundles
  └── cache/                          # Temporary scan cache
```

ForgeBE is **local-first** by design. No data leaves your machine unless you choose to share a bundle.

**Git-safe**: ForgeBE never writes to your project's git-tracked directory by default. All state stays in `~/.forgebe/`. Adapter exports are written to your local exports directory, not the repo.

---

## 🧩 Contract Extractor

ForgeBE can read existing documentation and AI instruction files to extract engineering rules:

- Reads `CLAUDE.md`, `AGENTS.md`, `.cursorrules`, `AI_CONTEXT.md`, `CONTRIBUTING.md`, `README.md`
- Infers testing strategy, dependency policies, and sensitive areas
- Detects whether a file is an AI instruction file vs. regular documentation
- Foundation for future multi-file orchestration and conflict detection

See [Contract Extractor Roadmap](docs/contract-extractor-roadmap.md) for planned enhancements.

---

## 🔒 Security

- **ZipSlip protection**: Bundle import validates all paths and blocks traversal attempts
- **Symlink hardening**: Import resolves and validates symlinks to prevent escape
- **No secrets in code**: Use environment variables or credential files
- **Local-first**: State never leaves your machine by default
- **Atomic writes**: Profile integrity maintained during concurrent access

Run `forgebe doctor` regularly to validate your project setup.

---

## 🛠 Development

```bash
git clone https://github.com/farid-alfernass/forgebe.git
cd forgebe

make build      # Build the binary
make test       # Full test suite
make fmt        # Format code
make vet        # Static analysis
make lint       # golangci-lint (if installed)
```

### Requirements
- Go 1.22+
- No runtime dependencies beyond stdlib + Cobra/Viper/YAML

---

## 📋 Project Status

- **v0.1.0** — Proof of Concept: Core models, storage, CLI skeleton
- **v0.5.0** — Usable MVP: Scan, init, doctor, profile, adapters, briefs, bundles
- **v1.0.0** (planned) — Production Ready: Full security audit, comprehensive docs, CI/CD integration, v1 profile schema

---

## 🤝 Contributing

Contributions are welcome! Please see `CONTRIBUTING.md` for guidelines.

### Development principles
- **Local-first by default**: State stays in `~/.forgebe/`
- **Git-safe**: No tracked files in target repos without explicit opt-in
- **Model-agnostic**: Not tied to any AI vendor or tool
- **Project-aware**: Adapts to existing codebase conventions
- **Portable context**: Export/import for team handoff

---

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🙏 Acknowledgements

ForgeBE was built with feedback from real backend engineers who use AI tools daily. The framework is designed to solve the consistency problem they face: every AI tool needs the same engineering context, but none of them share it.
