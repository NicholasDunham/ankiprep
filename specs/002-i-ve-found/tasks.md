# Tasks: Cloze Deletion Colon Exception

**Input**: Design documents from `/specs/002-i-ve-found/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Found: Go 1.21+ CLI tool extending existing ankiprep
   → Extract: tech stack (Go, cobra, unicode/norm), single project structure
2. Load optional design documents:
   → data-model.md: Extract entities → ClozeDeletionBlock, TypographyContext, TypographyResult
   → contracts/: typography-service.md, cli-interface.md → contract test tasks
   → research.md: Extract decisions → cloze detection, integration approach
3. Generate tasks by category:
   → Setup: Go module, dependencies, linting already configured
   → Tests: contract tests, integration tests from quickstart scenarios
   → Core: models, typography service enhancement, CLI integration
   → Integration: French typography pipeline integration
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
   → Typography service enhancement covered ✓
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `internal/`, `cmd/`, `tests/` at repository root
- Go package structure: `internal/models/`, `internal/services/`, `cmd/ankiprep/`
- Test structure: `tests/unit/`, `tests/integration/`, `tests/contract/`

## Phase 3.1: Setup
- [x] T001 [P] Verify Go module dependencies for cloze deletion feature (go.mod already exists)
- [x] T002 [P] Create test directory structure for cloze functionality
- [x] T003 [P] Create test fixtures directory with sample CSV files containing cloze blocks

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T004 [P] Contract test for TypographyService ProcessFrenchText method ✅ COMPLETED
- [x] T005 [P] Contract test for ProcessFrenchTextRequest validation ✅ COMPLETED
- [x] T006 [P] Contract test for AnkiCard model methods ✅ COMPLETED
- [x] T007 [P] Contract test for TypographyContext model methods ✅ COMPLETED
- [x] T008 [P] Integration test with CSV fixture data ✅ COMPLETED
- [x] T009 [P] CLI command integration test ✅ COMPLETED
- [x] T010 [P] End-to-end acceptance test ✅ COMPLETED

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T011 [P] ClozeDeletionBlock model struct and validation in internal/models/cloze.go
- [ ] T012 [P] TypographyContext model struct in internal/models/typography.go
- [ ] T013 [P] TypographyResult model struct in internal/models/typography.go (same file as T012)
- [ ] T014 [P] Cloze detection regex patterns and parsing logic in internal/services/cloze_parser.go
- [ ] T015 ClozeDeletionHandler implementation in internal/services/typography.go (extends existing file)
- [ ] T016 Enhanced ProcessFrenchText method with cloze awareness in internal/services/typography.go
- [ ] T017 Cloze syntax colon masking logic in internal/services/typography.go
- [ ] T018 Error handling and warning logging for malformed cloze blocks in internal/services/typography.go

## Phase 3.4: Integration
- [ ] T019 Integrate cloze detection into existing French typography pipeline in internal/services/typography.go
- [ ] T020 Update CLI interface to support enhanced French typography in cmd/ankiprep/main.go (extends existing functionality)
- [ ] T021 Add progress reporting for cloze block processing in internal/app/processor.go
- [ ] T022 Update existing CSV processing to use enhanced typography service in internal/app/processor.go

## Phase 3.5: Polish
- [ ] T023 [P] Unit tests for ClozeDeletionBlock validation in tests/unit/models/cloze_test.go
- [ ] T024 [P] Unit tests for cloze detection patterns in tests/unit/services/cloze_parser_test.go
- [ ] T025 [P] Unit tests for colon masking logic in tests/unit/services/typography_test.go
- [ ] T026 Performance test with large CSV files containing many cloze blocks in tests/performance/cloze_performance_test.go
- [ ] T027 [P] Update CLI help text to document cloze deletion support in cmd/ankiprep/main.go
- [ ] T028 [P] Update README.md with cloze deletion feature documentation
- [ ] T029 Run comprehensive quickstart validation scenarios from quickstart.md

## Dependencies

### Critical Dependencies
- **Tests (T004-T010) MUST complete before implementation (T011-T018)**
- **Core models (T011-T013) before service implementation (T014-T018)**
- **Service implementation (T014-T018) before integration (T019-T022)**
- **Integration complete before polish (T023-T029)**

### File Dependencies
- T012 and T013 modify same file (internal/models/typography.go) - must be sequential
- T015-T018 all modify internal/services/typography.go - must be sequential
- T020 and T027 modify cmd/ankiprep/main.go - must be sequential
- T021 and T022 modify internal/app/processor.go - must be sequential

### Logical Dependencies
- T011 (ClozeDeletionBlock) blocks T014 (cloze parsing)
- T014 (cloze parsing) blocks T015 (handler implementation)
- T015 (handler) blocks T016 (ProcessFrenchText enhancement)
- T019 (pipeline integration) blocks T020 (CLI integration)

## Parallel Execution Examples

### Phase 3.2: Test Creation (All Parallel)
```bash
# Launch all contract and integration tests simultaneously:
Task: "Contract test for TypographyService ProcessFrenchText method in tests/contract/typography_service_test.go"
Task: "Contract test for CLI cloze processing behavior in tests/contract/cli_interface_test.go"
Task: "Integration test for basic cloze processing in tests/integration/cloze_basic_test.go"
Task: "Integration test for multiple cloze blocks in tests/integration/cloze_multiple_test.go"
Task: "Integration test for nested content colons in tests/integration/cloze_nested_test.go"
Task: "Integration test for malformed cloze handling in tests/integration/cloze_malformed_test.go"
Task: "Integration test for complex cloze with hints in tests/integration/cloze_hints_test.go"
```

### Phase 3.3: Model Creation (Parallel Safe)
```bash
# Launch independent model files:
Task: "ClozeDeletionBlock model struct and validation in internal/models/cloze.go"
Task: "Cloze detection regex patterns and parsing logic in internal/services/cloze_parser.go"
# Note: T012-T013 must be sequential (same file)
```

### Phase 3.5: Unit Tests and Documentation (Parallel Safe)
```bash
# Launch independent polish tasks:
Task: "Unit tests for ClozeDeletionBlock validation in tests/unit/models/cloze_test.go"
Task: "Unit tests for cloze detection patterns in tests/unit/services/cloze_parser_test.go"
Task: "Unit tests for colon masking logic in tests/unit/services/typography_test.go"
Task: "Performance test with large CSV files in tests/performance/cloze_performance_test.go"
Task: "Update README.md with cloze deletion feature documentation"
```

## Test Scenarios Mapping

### From quickstart.md:
- **Scenario 1** → T006: Basic cloze processing test
- **Scenario 2** → T007: Multiple cloze blocks test  
- **Scenario 3** → T008: Nested content colons test
- **Scenario 4** → T009: Malformed cloze handling test
- **Scenario 5** → T010: Complex cloze with hints test

### From contracts/:
- **typography-service.md** → T004: ProcessFrenchText contract test
- **cli-interface.md** → T005: CLI behavior contract test

## Notes
- **[P] tasks**: Different files, no dependencies, can run simultaneously
- **Sequential tasks**: Same file modifications, must run in dependency order
- **TDD critical**: All tests (T004-T010) must fail before implementation starts
- **Backward compatibility**: All existing functionality must remain unchanged
- **Performance target**: Sub-second processing maintained for typical files
- **Error handling**: Graceful degradation for malformed cloze blocks

## Validation Checklist
*GATE: Checked before task execution*

- [x] All contracts have corresponding tests (T004, T005)
- [x] All entities have model tasks (T011, T012, T013)
- [x] All tests come before implementation (T004-T010 → T011-T018)
- [x] Parallel tasks truly independent (different files)
- [x] Each task specifies exact file path
- [x] No [P] task modifies same file as another [P] task
- [x] Quickstart scenarios mapped to integration tests
- [x] Performance and documentation covered in polish phase