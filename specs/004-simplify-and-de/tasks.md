# Tasks: Simplify and De-engineer Codebase

**Input**: Design documents from `/specs/004-simplify-and-de/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Found: Go 1.21+ CLI application with simplification requirements
   → Extract: cobra CLI, standard library focus, file system storage
2. Load optional design documents:
   → research.md: Identified over-engineering patterns to remove
   → contracts/: CLI interface and processing contracts to preserve
   → quickstart.md: Validation scenarios for regression testing
3. Generate tasks by category:
   → Setup: baseline establishment, complexity measurement
   → Validation: continuous regression testing
   → Simplification: targeted removal of over-engineering
   → Integration: consolidation and cleanup
   → Polish: final validation and documentation
4. Apply task rules:
   → Validation tests before each simplification step
   → Independent simplification areas marked [P]
   → Sequential dependency for shared components
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph focusing on incremental safety
7. Create parallel execution examples for independent areas
8. Validate task completeness for comprehensive simplification
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files/areas, no dependencies)
- Include exact file paths and validation steps

## Path Conventions
- **Single project**: Repository root with existing structure
- Current structure: `cmd/`, `internal/`, `pkg/`, `tests/`
- Preserve functional structure while simplifying implementation

## Phase 3.1: Setup and Baseline
- [x] T001 Establish pre-refactoring baseline measurements and outputs
- [x] T002 [P] Measure current code complexity metrics (LOC, cyclomatic, packages)
- [x] T003 [P] Create comprehensive test output baselines for all CLI scenarios
- [x] T004 Document current architecture and identify simplification targets

## Phase 3.2: Configuration Simplification ⚠️ VALIDATE AFTER EACH CHANGE
**CRITICAL: Run validation tests after each simplification to ensure no regression**
- [x] T005 [P] Remove unnecessary linting configuration (.golangci.yml) for single-user tool
- [x] T006 [P] Audit go.mod dependencies - remove only if truly unused or provide no complexity benefit over standard library
- [x] T007 Validate CLI functionality unchanged after configuration removal

## Phase 3.3: Package Structure Flattening (ONLY after T007 passes)
- [x] T008 [P] Analyze current package structure in internal/ for over-abstraction
- [x] T009 [P] Flatten internal/services/ - merge simple wrapper services
- [ ] T010 [P] Consolidate internal/models/ - eliminate unnecessary abstractions
- [ ] T011 Merge internal/app/ functionality into main command structure
- [x] T012 Validate all CLI operations work identically after package changes

## Phase 3.4: Interface Simplification
- [x] T013 [P] Remove excessive interfaces in internal/services/ where concrete types suffice
- [x] T014 [P] Simplify error handling - eliminate complex error wrapper chains
- [x] T015 [P] Remove unnecessary dependency injection patterns in cmd/ankiprep/main.go
- [x] T016 Consolidate file handling logic across internal/services/
- [x] T017 Validate CLI contract compliance after interface changes

## Phase 3.5: Code Consolidation
- [x] T018 [P] Merge duplicate CSV parsing logic across services
- [x] T019 [P] Consolidate typography processing into single focused module
- [x] T020 [P] Remove redundant validation logic in multiple layers
- [x] T021 Simplify memory monitoring - remove enterprise patterns for single-user tool
- [x] T022 Validate processing accuracy with all CLI flag combinations

## Phase 3.6: Test Suite Optimization
**NOTE: Test simplification focuses on reducing complexity while maintaining coverage requirements and constitutional compliance. Remove test abstraction overhead, not test capability.**
- [x] T023 [P] Remove excessive test abstraction layers in tests/unit/
- [x] T024 [P] Consolidate integration tests - eliminate redundant test patterns
- [x] T025 [P] Simplify contract tests to focus on actual CLI behavior
- [x] T026 Keep performance tests but simplify measurement approach
- [x] T027 Validate entire test suite passes with simplified structure

## Phase 3.7: Final Integration and Documentation

- [x] T028 [P] Validate complete functionality preservation across all features
- [x] T029 Calculate and document complexity metrics showing quantified improvement
- [x] T030 Update documentation to reflect simplified architecture
- [x] T031 Comprehensive final validation and project completion summary

## Dependencies
- Baseline (T001-T004) before any simplification
- T007 validation before package changes (T008-T012)
- T012 validation before interface changes (T013-T017)  
- T017 validation before consolidation (T018-T022)
- T022 validation before test optimization (T023-T027)
- T027 validation before final integration (T028-T033)
- All validation tasks must pass before proceeding to next phase

## Parallel Example
```
# Launch T002-T003 together during setup:
Task: "Measure current code complexity metrics in complexity-before.txt"
Task: "Create comprehensive test output baselines using examples/spanish_vocabulary.tsv"

# Launch T005-T006 together for configuration:
Task: "Remove .golangci.yml and validate it's not required for build"
Task: "Audit go.mod dependencies - only remove if unused or no benefit vs stdlib"

# Launch T009-T010 together for package flattening:
Task: "Merge simple wrapper services in internal/services/"
Task: "Eliminate abstraction layers in internal/models/"
```

## Validation Commands
**Critical: Run after each simplification phase**
```bash
# After each change, validate CLI behavior unchanged:
go build -o ankiprep-current ./cmd/ankiprep
./ankiprep-current examples/spanish_vocabulary.tsv -o test.csv
diff baseline-simple.csv test.csv

# Validate all flag combinations:
./ankiprep-current examples/spanish_vocabulary.tsv -fqs -v -o test-full.csv  
diff baseline-full.csv test-full.csv

# Ensure tests pass:
go test ./... -v
```

## Success Criteria
- All baseline output comparisons show identical results
- Test suite passes without modification to test expectations
- Measurable reduction in code complexity (≥10% LOC reduction)
- Reduced package count and abstraction layers
- Preserved CLI functionality per contracts/cli-interface.md
- No performance regression in processing operations

## Risk Mitigation
- Small incremental changes with validation after each
- Preserve all existing test expectations
- Maintain git commits for easy rollback
- Focus on internal implementation, not external behavior
- Prioritize simplification that reduces maintenance burden
- Keep dependencies that provide complexity/functionality benefits - only remove if truly unused or redundant

## Notes
- [P] tasks = independent areas that don't modify same files
- All CLI contracts must be preserved exactly
- Validation is mandatory after each simplification phase
- Rollback any change that fails validation
- Document complexity reduction achieved in final metrics