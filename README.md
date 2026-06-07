# ForgeBE

**ForgeBE** is an open-source, local-first, AI-agnostic backend engineering framework that helps developers keep AI-assisted delivery consistent, safe, and project-aware across existing and new codebases.

## Status

ForgeBE is under active development. Current milestone: **Day 1 foundation** toward a production-ready v1.0.

## Vision

ForgeBE provides a canonical project profile and engineering contract, then renders AI-specific context adapters for tools such as Claude Code, Codex, Hermes Agent, Cursor, Copilot, and generic LLM workflows.

## Principles

- **Local-first**: project profiles are stored locally under `~/.forgebe/`
- **Git-safe by default**: ForgeBE does not write tracked files into a target repo unless explicitly asked
- **AI-agnostic**: one engineering contract, many AI adapters
- **Project-aware**: guided onboarding + auto-discovery fallback
- **Portable**: export compact summaries and bundle archives for device migration

## Planned v1 command surface

```bash
forgebe init
forgebe scan
forgebe doctor
forgebe profile show
forgebe export summary
forgebe export bundle
forgebe import bundle <file>
forgebe export adapter <tool>
forgebe brief <tool>
```

## Day 1 foundation

Implemented today:
- Go module and repo skeleton
- canonical profile model
- local storage paths
- atomic file writes
- project id generation
- root CLI skeleton
- Makefile
- initial documentation

## Build

```bash
make build
```

## Test

```bash
make test
```

## License

MIT
