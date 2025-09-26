<!--
Sync Impact Report:
- Version change: template → 1.0.0 (initial creation)
- Modified principles: All principles defined for Go CLI application
- Added sections: Core Principles, Development Standards, Quality Gates, Governance
- Removed sections: None (initial creation)
- Templates requiring updates: ✅ aligned with plan-template.md, spec-template.md, tasks-template.md
- Follow-up TODOs: None
-->

# AnkiPrep Constitution

## Core Principles

### I. POSIX Compliance

All command-line interfaces MUST follow POSIX conventions and Linux/Unix best practices.
Commands MUST use standard exit codes (0 for success, non-zero for errors), support
standard flags (--help, --version), and follow conventional argument patterns.
Text processing MUST handle UTF-8 encoding and respect locale settings.

**Rationale**: Ensures the CLI integrates seamlessly with existing Unix toolchains
and user expectations.

### II. Test-First Development

Testing is NON-NEGOTIABLE. Every feature MUST have tests written before implementation.
Follow strict Red-Green-Refactor cycles: write failing test → user approval →
implement minimum code to pass → refactor. Unit tests for all Go packages,
integration tests for CLI command flows.

**Rationale**: Ensures reliability and maintainability in a tool that users depend on
for their study workflows.

### III. Clean CLI Interface

All CLI commands MUST provide clear, consistent interfaces with proper help text,
examples, and error messages. Support both human-readable and machine-parseable output
formats (JSON when appropriate). Input via stdin/args → output to stdout,
errors to stderr.

**Rationale**: Users need intuitive, scriptable tools that integrate well with
their existing workflows.

### IV. Go Conventions

Code MUST follow standard Go conventions: gofmt formatting, golint compliance,
effective Go patterns, proper error handling with wrapped errors, and idiomatic
package structure. Use Go modules for dependency management, follow semantic
import versioning.

Local linting tools (golangci-lint, go vet, gofmt) are RECOMMENDED for development
but not required for single-user hobby projects. Focus on code clarity and Go idioms
over comprehensive linting infrastructure.

**Rationale**: Maintains code quality and leverages the Go ecosystem's established
best practices while avoiding unnecessary tooling overhead for personal projects.

### V. Performance & Reliability

The application MUST be fast and reliable for typical study preparation workloads.
Performance targets: sub-second response for typical operations, graceful handling
of large datasets, minimal memory footprint. Proper logging for debugging,
no silent failures.

**Rationale**: Study preparation tools must be responsive and dependable to support
effective learning workflows.

## Development Standards

All code MUST be written in Go using the latest stable version. External dependencies
MUST be justified and minimal - prefer standard library solutions where possible.
All public APIs MUST be documented with Go doc comments. Configuration MUST follow
XDG Base Directory specification for file placement.

For single-user hobby projects, code quality is maintained through Go conventions
and testing rather than extensive tooling infrastructure. Focus on readable,
maintainable code over comprehensive automation.

## Quality Gates

Before any release:

- All tests MUST pass (unit and integration)
- Code coverage SHOULD be above 80% (measured locally, not enforced by CI)
- Basic go vet checks SHOULD pass (run locally as needed)
- CLI help text and examples MUST be accurate
- Cross-platform compatibility MUST be verified (Linux, macOS, Windows)
- Performance regression tests MUST pass

Manual testing of common user workflows is required for major releases.
Local development tools (linting, formatting) are encouraged but not mandated
for single-user projects.

## Governance

This constitution supersedes all other development practices. Amendments require
documentation of rationale, impact assessment, and migration plan. All development
decisions MUST align with these principles.

Complexity that violates these principles MUST be justified with clear business
rationale. When in doubt, choose the simpler, more conventional approach.

**Version**: 1.0.0 | **Ratified**: 2025-09-24 | **Last Amended**: 2025-09-24
