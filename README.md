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

# Scan existing project
forgebe scan .

# Brief your AI tool
forgebe brief claude

# Export adapter files
forgebe export adapter cursor

# Sync AI context files into the project root
forgebe sync --dry-run
forgebe sync

# Keep AI context files up-to-date while coding
forgebe watch
```

See [docs/cli-reference.md](docs/cli-reference.md) for the full reference.

## Project Status

**v1.0.0** — Production-ready. See [CHANGELOG.md](CHANGELOG.md) for release history.

## Docs

| Document | Description |
|----------|-------------|
| [Architecture](docs/architecture.md) | System design, data flow, extension points |
| [Security](docs/security.md) | Threat model, security boundaries, usage guidelines |
| [CLI Reference](docs/cli-reference.md) | Complete command reference and examples |
| [Contributing](CONTRIBUTING.md) | Development guidelines, quality gates |

## License

MIT
