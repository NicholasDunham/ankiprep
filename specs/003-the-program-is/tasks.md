# Tasks: Remove Duplicate CSV Header in Anki Output

**Input**: Design documents from `/specs/003-the-program-is/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Found: Go 1.21+ CLI application with cobra framework
   → Extract: Single project structure, existing modular codebase
2. Load design documents:
   → data-model.md: ProcessingOptions, CSVRecord, ProcessingResult entities
   → contracts/: CLI interface and processing interface contracts
   → research.md: Header removal strategy, flag integration patterns
   → quickstart.md: 5 test scenarios for validation
3. Generate tasks by category:
   → Setup: flag integration, dependencies
   → Tests: contract tests, integration tests (TDD approach)
   → Core: data models, processing logic, CLI command updates
   → Integration: error handling, file operations
   → Polish: unit tests, performance validation, documentation
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests ✓
   → All entities have models ✓
   → All CLI scenarios covered ✓
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: Existing structure with `internal/`, `cmd/`, `tests/`
- Paths based on current ankiprep project structure

## Phase 3.1: Setup
- [x] T001 Add --keep-header/-k flag to cobra root command in `cmd/ankiprep/main.go`
- [x] T002 [P] Update go.mod dependencies if needed (already has cobra and required libs)
- [x] T003 [P] Configure additional linting rules for new code

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T004 [P] Contract test for CLI --keep-header flag parsing in `tests/unit/cmd/flag_test.go`
- [ ] T005 [P] Contract test for ProcessCSV function signature in `tests/unit/services/processor_test.go`
- [ ] T006 [P] Integration test default header removal in `tests/integration/header_removal_test.go`
- [ ] T007 [P] Integration test --keep-header preservation in `tests/integration/keep_header_test.go`
- [ ] T008 [P] Integration test error handling scenarios in `tests/integration/error_handling_test.go`
- [ ] T009 [P] Integration test CLI exit codes in `tests/integration/exit_codes_test.go`
- [ ] T010 [P] Integration test backward compatibility in `tests/integration/backward_compatibility_test.go`

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T011 [P] ProcessingOptions struct in `internal/models/processing_options.go`
- [ ] T012 [P] CSVRecord struct in `internal/models/csv_record.go`
- [ ] T013 [P] ProcessingResult struct in `internal/models/processing_result.go`
- [ ] T014 Extend ProcessCSV function signature in `internal/app/processor.go`
- [ ] T015 Implement header detection logic in `internal/app/processor.go`
- [ ] T016 Add --keep-header flag handling in CLI command chain
- [ ] T017 Update processor to skip first row based on options
- [ ] T018 Add validation for ProcessingOptions

## Phase 3.4: Integration
- [ ] T019 Update error handling to match contract specifications
- [ ] T020 Add file path validation (input exists, output writable)
- [ ] T021 Implement atomic file writing (temp file + rename)
- [ ] T022 Add descriptive error messages with context
- [ ] T023 Update CLI help text and usage examples

## Phase 3.5: Polish
- [ ] T024 [P] Unit tests for ProcessingOptions validation in `tests/unit/models/processing_options_test.go`
- [ ] T025 [P] Unit tests for CSVRecord operations in `tests/unit/models/csv_record_test.go`
- [ ] T026 [P] Unit tests for ProcessingResult handling in `tests/unit/models/processing_result_test.go`
- [ ] T027 [P] Unit tests for header detection logic in `tests/unit/app/header_logic_test.go`
- [ ] T028 Performance tests with large CSV files (1000+ rows)
- [ ] T029 [P] Update README.md with new flag documentation
- [ ] T030 [P] Update CLI help text and man page
- [ ] T031 Execute quickstart.md test scenarios for final validation
- [ ] T032 Remove any code duplication from implementation

## Dependencies
- Tests (T004-T010) MUST complete before implementation (T011-T023)
- T011, T012, T013 must complete before T014 (function signature update)
- T014 blocks T015, T017 (processor modifications)
- T016 blocks T023 (CLI help updates)
- T018, T019, T020, T021, T022 can run in parallel after core implementation
- Implementation (T011-T023) before polish (T024-T032)
- T031 requires all implementation tasks complete

## Parallel Example
```bash
# Launch T004-T010 together (all different test files):
Task: "Contract test for CLI --keep-header flag parsing in tests/unit/cmd/flag_test.go"
Task: "Contract test for ProcessCSV function signature in tests/unit/services/processor_test.go"  
Task: "Integration test default header removal in tests/integration/header_removal_test.go"
Task: "Integration test --keep-header preservation in tests/integration/keep_header_test.go"
Task: "Integration test error handling scenarios in tests/integration/error_handling_test.go"
```

```bash
# Launch T011-T013 together (all different model files):
Task: "ProcessingOptions struct in internal/models/processing_options.go"
Task: "CSVRecord struct in internal/models/csv_record.go"
Task: "ProcessingResult struct in internal/models/processing_result.go"
```

```bash
# Launch T024-T027 together (all different test files):
Task: "Unit tests for ProcessingOptions validation in tests/unit/models/processing_options_test.go"
Task: "Unit tests for CSVRecord operations in tests/unit/models/csv_record_test.go"
Task: "Unit tests for ProcessingResult handling in tests/unit/models/processing_result_test.go"
Task: "Unit tests for header detection logic in tests/unit/app/header_logic_test.go"
```

## Notes
- [P] tasks target different files and have no dependencies between them
- All tests MUST fail initially to follow TDD principles
- Focus on extending existing `internal/app/processor.go` rather than creating new processors
- Maintain backward compatibility - existing functionality unchanged
- CLI flag follows POSIX conventions: long form `--keep-header`, short form `-k`
- Error messages must be descriptive and user-friendly
- Atomic file operations prevent partial writes

## Validation Checklist
- [ ] All contracts from `contracts/` directory have corresponding tests
- [ ] All entities from `data-model.md` have struct implementations  
- [ ] All quickstart scenarios have integration tests
- [ ] CLI interface contract fully implemented
- [ ] Processing interface contract fully implemented
- [ ] Constitution principles maintained (POSIX, TDD, Go conventions)
- [ ] No regressions in existing functionality

## Task Generation Rules
- Contract test per interface specification
- Model implementation per data entity
- Integration test per user scenario
- Unit test per implementation file
- Sequential for same-file modifications
- Parallel for independent files