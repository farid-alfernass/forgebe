# Security Policy

## Supported Versions

During the initial development phase, only the latest main branch and the latest tagged release are supported.

## Reporting a Vulnerability

Please do not open public issues for sensitive security reports. Instead, contact the maintainer privately and include:

- affected version or commit
- reproduction steps
- impact assessment
- suggested remediation if available

## Security Principles

ForgeBE is designed with the following defaults:

- local-first state under the user's home directory
- no secrets required for core functionality
- no automatic writes into a target repository by default
- atomic local writes for profile persistence
- portable exports treated as user-managed artifacts
