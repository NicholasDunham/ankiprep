# Tasks: CSV to Anki Processor

**Input**: Design documents from `/specs/001-command-line-application/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory ✓
   → Extract: Go 1.21+, cobra CLI, unicode/norm, CSV processing
2. Load design documents: ✓
   → data-model.md: InputFile, DataEntry, OutputFile, ProcessingReport, TypographyProcessor
   → contracts/: CLI interface with options and exit codes
   → research.md: External libraries decisions
3. Generate tasks by category: ✓
   → Setup: Go project init, dependencies (cobra, x/text), linting
   → Tests: CLI contract tests, integration scenarios
   → Core: models, CSV processor, typography, CLI commands
   → Integration: file I/O, error handling, progress reporting
   → Polish: unit tests, performance validation, documentation
4. Apply task rules: ✓
   → Different files = mark [P] for parallel
   → CLI implementation = sequential (cobra root command)
   → Tests before implementation (TDD)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Single project structure:
- `cmd/ankiprep/` - Main application entry point
- `internal/` - Internal packages (models, services)
- `pkg/` - Reusable packages
- `tests/` - All test files

## Phase 3.1: Setup
- [ ] T001 Create Go module and project structure (cmd/, internal/, pkg/, tests/)
- [ ] T002 Initialize go.mod with dependencies: github.com/spf13/cobra, golang.org/x/text/unicode/norm
- [ ] T003 [P] Configure linting tools (golangci-lint config in .golangci.yml)
- [ ] T004 [P] Setup GitHub Actions CI workflow in .github/workflows/ci.yml

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T005 [P] CLI contract test for help/version flags in tests/contract/test_cli_basic.go
- [ ] T006 [P] CLI contract test for file processing in tests/contract/test_cli_processing.go
- [ ] T007 [P] CLI contract test for French typography flag in tests/contract/test_cli_french.go
- [ ] T008 [P] CLI contract test for error handling in tests/contract/test_cli_errors.go
- [ ] T009 [P] Integration test single file processing in tests/integration/test_single_file.go
- [ ] T010 [P] Integration test multiple file merging in tests/integration/test_multiple_files.go
- [ ] T011 [P] Integration test duplicate detection in tests/integration/test_duplicates.go
- [ ] T012 [P] Integration test French typography in tests/integration/test_french_typography.go
- [ ] T013 [P] Integration test smart quotes conversion in tests/integration/test_smart_quotes.go
- [ ] T014 [P] Integration test multiline content handling in tests/integration/test_multiline_content.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T015 [P] InputFile model in internal/models/input_file.go
- [ ] T016 [P] DataEntry model in internal/models/data_entry.go
- [ ] T017 [P] OutputFile model in internal/models/output_file.go
- [ ] T018 [P] ProcessingReport model in internal/models/processing_report.go
- [ ] T019 [P] TypographyProcessor model in internal/models/typography_processor.go
- [ ] T020 [P] CSV parser service in internal/services/csv_parser.go
- [ ] T021 [P] File validator service in internal/services/file_validator.go
- [ ] T022 [P] Duplicate detector service in internal/services/duplicate_detector.go
- [ ] T023 [P] Column merger service in internal/services/column_merger.go
- [ ] T024 [P] Typography service in internal/services/typography_service.go
- [ ] T025 [P] Anki formatter service in internal/services/anki_formatter.go
- [ ] T026 [P] Progress reporter service in internal/services/progress_reporter.go
- [ ] T027 Main application processor in internal/app/processor.go
- [ ] T028 CLI root command setup with cobra in cmd/ankiprep/root.go
- [ ] T029 CLI process command implementation in cmd/ankiprep/process.go
- [ ] T030 Main entry point in cmd/ankiprep/main.go

## Phase 3.4: Integration
- [ ] T031 File I/O error handling in internal/services/file_service.go
- [ ] T032 Progress indicator integration for large files (>10MB or >5s processing)
- [ ] T033 Partial output file cleanup on processing failure
- [ ] T034 Exit code handling and error message formatting
- [ ] T035 Memory usage monitoring and optimization for large datasets

## Phase 3.5: Polish
- [ ] T036 [P] Unit tests for InputFile model in tests/unit/test_input_file.go
- [ ] T037 [P] Unit tests for DataEntry model in tests/unit/test_data_entry.go
- [ ] T038 [P] Unit tests for CSV parser in tests/unit/test_csv_parser.go
- [ ] T039 [P] Unit tests for duplicate detection in tests/unit/test_duplicate_detector.go
- [ ] T040 [P] Unit tests for typography processing in tests/unit/test_typography_processor.go
- [ ] T041 [P] Performance tests for large file processing in tests/performance/test_large_files.go
- [ ] T042 [P] Memory usage tests in tests/performance/test_memory_usage.go
- [ ] T043 [P] Cross-platform compatibility tests in tests/integration/test_cross_platform.go
- [ ] T044 [P] Update README.md with installation and usage instructions
- [ ] T045 [P] Add example CSV files in examples/ directory
- [ ] T046 Execute manual testing scenarios from quickstart.md
- [ ] T047 Build and test cross-platform binaries (Linux, macOS, Windows)

## Dependencies
- Setup (T001-T004) before everything
- All tests (T005-T014) before implementation (T015-T030)
- Models (T015-T019) before services (T020-T026)
- Services before app processor (T027)
- App processor before CLI commands (T028-T030)
- Core implementation before integration (T031-T035)
- Everything before polish (T036-T047)

## Parallel Example
```bash
# Launch setup tasks together:
Task: "Configure linting tools (golangci-lint config in .golangci.yml)"
Task: "Setup GitHub Actions CI workflow in .github/workflows/ci.yml"

# Launch model tests together:
Task: "CLI contract test for help/version flags in tests/contract/test_cli_basic.go"
Task: "CLI contract test for file processing in tests/contract/test_cli_processing.go"
Task: "CLI contract test for French typography flag in tests/contract/test_cli_french.go"

# Launch model implementation together:
Task: "InputFile model in internal/models/input_file.go"
Task: "DataEntry model in internal/models/data_entry.go"
Task: "OutputFile model in internal/models/output_file.go"
```

## Key Implementation Notes
- **CSV Processing**: Use standard library `encoding/csv` for parsing and writing
- **CLI Framework**: Use `github.com/spf13/cobra` for command structure and flag handling
- **Typography**: Use `golang.org/x/text/unicode/norm` for Unicode normalization
- **File Structure**: Single binary CLI application with clean internal package organization
- **Testing**: Follow TDD approach - all tests must fail before implementation
- **Error Handling**: Fail-fast with clear error messages and proper exit codes
- **Memory Management**: Process files in memory with progress indicators for large files
- **Output Format**: Always generate Anki-compatible CSV with proper headers

## Validation Checklist
*GATE: Checked before execution*

- [x] All CLI interface contracts have corresponding tests (T005-T008)
- [x] All entities have model implementation tasks (T015-T019)  
- [x] All integration scenarios have test tasks (T009-T014)
- [x] All tests come before implementation (Phase 3.2 → 3.3)
- [x] Parallel tasks are truly independent (different files)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Dependencies clearly defined and enforceable