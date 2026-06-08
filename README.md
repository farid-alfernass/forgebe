# ForgeBE

**Open-Source AI-Agnostic Backend Engineering Framework**

ForgeBE is a local-first, AI-agnostic backend engineering framework designed to keep AI-assisted software delivery consistent, safe, and project-aware.

## Quick Start

### macOS (Homebrew)

```bash
brew tap farid-alfernass/tap
brew install forgebe
```

### Linux (Binary)

```bash
curl -fsSL https://github.com/farid-alfernass/forgebe/releases/latest/download/forgebe_Linux_x86_64.tar.gz | tar xz
sudo mv forgebe /usr/local/bin/
```

### From source

```bash
git clone https://github.com/farid-alfernass/forgebe.git
cd forgebe
make build
sudo make install
```

## Usage

```bash
# Initialize a project (guided onboarding)
forgebe init

# Adopt existing project
forgebe adopt

# Scan existing project
forgebe scan .

# Generate boilerplate for new project
forgebe scaffold --template go-api --output ./new-backend

# Verify project against engineering policies
forgebe verify

# Health check for CI/pre-commit
forgebe check

# Brief your AI tool
forgebe brief claude

# Generate context-aware prompts for AI tasks
forgebe prompt implement "Add user authentication"
forgebe prompt review "Check the payment handler"

# Export adapter files
forgebe export adapter cursor

# Sync AI context files into the project root
forgebe sync --dry-run
forgebe sync

# Keep AI context files up-to-date while coding
forgebe watch
```

See [docs/cli-reference.md](docs/cli-reference.md) for the full reference, or [docs/phase10-feature-expansion.md](docs/phase10-feature-expansion.md) for Phase 10 feature guide with detailed examples.

## Project Status

**v1.2.0** — Phase 10 Complete (Feature Expansion). Production-ready with 5 new operational commands. See [CHANGELOG.md](CHANGELOG.md) for release history.

## Docs

| Document | Description |
|----------|-------------|
| [Architecture](docs/architecture.md) | System design, data flow, extension points |
| [Security](docs/security.md) | Threat model, security boundaries, usage guidelines |
| [CLI Reference](docs/cli-reference.md) | Complete command reference and examples |
| [Phase 10 Guide](docs/phase10-feature-expansion.md) | Feature expansion: scaffold, adopt, verify, check, prompt |
| [Contributing](CONTRIBUTING.md) | Development guidelines, quality gates |

## License

MIT
