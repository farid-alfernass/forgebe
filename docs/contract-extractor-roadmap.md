# Contract Extractor Roadmap

ForgeBE Phase 5 introduces the first production-ready foundation for **Contract Extractor** — the capability to read documentation and AI-generated instruction files, then infer engineering rules back into ForgeBE.

## Goal

Allow ForgeBE to understand rules from files such as:
- `README.md`
- `CONTRIBUTING.md`
- `CLAUDE.md`
- `AGENTS.md`
- `.cursorrules`
- `AI_CONTEXT.md`
- `.github/copilot-instructions.md`

## Phase 5 Capabilities

Implemented in this phase:
- Detect whether a file is an AI instruction file.
- Infer testing strategy signals like `tdd` and `tests before implementation`.
- Infer dependency approval rules from human-written or AI-written docs.
- Infer preserve-structure / minimal-change constraints.
- Extract sensitive areas from markdown sections such as `Sensitive Areas`.
- Preserve raw signals for later auditing.

Current implementation entrypoint:
- `internal/contract/extractor.go`

## Current Limitations

Phase 5 extractor is rule-based and pattern-driven. It does **not** yet:
- parse every possible markdown style,
- use embeddings or LLM-based semantic extraction,
- merge multiple files into one normalized contract automatically,
- write extracted signals back into profile without explicit orchestration.

## Planned Next Steps

### Phase 6+
- Multi-file extraction orchestration:
  - scan docs folder,
  - prioritize AI instruction files,
  - merge rules with precedence.
- Confidence scoring per extracted rule.
- Conflict detection between discovered docs and local project profile.
- `forgebe doctor` enhancement to warn when docs and profile drift apart.
- Optional `forgebe extract contract` command.
- LLM-assisted semantic extraction mode for loosely structured docs.

## Extraction Precedence Proposal

Recommended merge precedence:
1. ForgeBE local profile (`~/.forgebe/...`)
2. Explicit ForgeBE engineering contract
3. AI instruction files (`CLAUDE.md`, `AGENTS.md`, `.cursorrules`, Copilot instructions)
4. `CONTRIBUTING.md`
5. `README.md`

## Why This Matters

This is the bridge between:
- local-first ForgeBE profile, and
- existing project rules that already live in docs.

It allows ForgeBE to become progressively smarter on mature codebases without forcing teams to rewrite their conventions first.
