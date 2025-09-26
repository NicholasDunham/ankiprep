# Tasks: Improve French Typography NNBSP Handling

**Input**: Design documents from `/specs/005-improve-french-typography/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory ✅
   → Tech stack: Go 1.21+, cobra CLI, golang.org/x/text/unicode/norm
   → Structure: Single CLI application (ankiprep command)
2. Load design documents ✅
   → data-model.md: Enhanced TypographyProcessor with improved NNBSP detection
   → contracts/: typography-interface.md, cli-interface.md  
   → research.md: Custom implementation approach, Unicode character handling
3. Generate tasks by category ✅
   → Setup: dependencies, project structure
   → Tests: contract tests, integration scenarios from quickstart
   → Core: enhance existing typography processor methods
   → Integration: CLI integration with existing file processing
   → Polish: validation, compliance, documentation
4. Apply task rules ✅
   → [P] for different files, sequential for same file
   → Tests before implementation (TDD approach)
5. Tasks numbered T001-T025 ✅
6. Dependencies and parallel execution documented ✅
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- All file paths are absolute and specific

## Path Conventions
- **Single project**: `internal/`, `cmd/`, `tests/` at repository root
- Follows existing ankiprep project structure

## Phase 3.1: Setup
- [x] T001 Review existing French typography logic in `internal/models/typography_processor.go`
- [x] T002 Verify French typography processing dependencies in `go.mod` (golang.org/x/text/unicode/norm already present)
- [x] T003 [P] Create typography processor test files structure

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests (Based on contracts/)
- [x] T004 [P] Typography interface contract tests in `tests/contract/typography_processor_test.go`
- [x] T005 [P] CLI interface contract tests in `tests/contract/cli_french_typography_test.go`

### Unit Tests (Based on data-model.md entities)
- [x] T006 [P] Enhanced TypographyProcessor tests in `tests/unit/models/typography_processor_test.go`
- [x] T007 [P] Enhanced applyFrenchTypography method tests in `tests/unit/models/typography_processor_test.go`
- [x] T008 [P] Enhanced applyGuillemetSpacing method tests in `tests/unit/models/typography_processor_test.go`

### Integration Tests (Based on quickstart.md scenarios)
- [x] T009 [P] Basic quote processing test in `tests/integration/french_quotes_test.go`
- [x] T010 [P] Punctuation spacing test in `tests/integration/french_punctuation_test.go`
- [x] T011 [P] NNBSP preservation test in `tests/integration/nnbsp_preservation_test.go`
- [x] T012 [P] Mixed content processing test in `tests/integration/mixed_french_content_test.go`
- [x] T013 [P] Error handling integration test in `tests/integration/typography_errors_test.go`

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Data Models (Based on data-model.md)
- [x] T014 [P] Enhance existing TypographyProcessor in `internal/models/typography_processor.go`

### Enhanced Typography Logic
- [x] T015 Improve applyFrenchTypography method to detect existing NNBSP in `internal/models/typography_processor.go`
- [x] T016 Improve applyGuillemetSpacing method to replace regular spaces with NNBSP in `internal/models/typography_processor.go`

### CLI Integration
- [ ] T017 Verify French typography integration in main command flow in `cmd/ankiprep/main.go`
- [ ] T018 Test enhanced French typography processing through existing applyTypography function

## Phase 3.4: Integration
- [ ] T019 Validate enhanced typography processing with existing file processing pipeline
- [ ] T020 Add error handling and validation for improved French typography processing

## Phase 3.5: Polish
- [ ] T021 [P] Run quickstart.md validation scenarios and update documentation
- [ ] T022 [P] Validate French typography compliance with standard French typography rules (FR-008)

## Dependencies
- Setup (T001-T003) before all other tasks
- All tests (T004-T013) before implementation (T014-T020)
- Model enhancement (T014) before method improvements (T015-T016)
- Method improvements before CLI integration (T017-T018)
- Integration (T019-T020) before polish (T021-T022)

## Parallel Example
```bash
# Launch contract tests together (T004-T005):
Task: "Typography interface contract tests in tests/contract/typography_processor_test.go"
Task: "CLI interface contract tests in tests/contract/cli_french_typography_test.go"

# Launch model unit tests together (T006-T008):
Task: "Enhanced TypographyProcessor tests in tests/unit/models/typography_processor_test.go"
Task: "Enhanced applyFrenchTypography method tests in tests/unit/models/typography_processor_test.go"
Task: "Enhanced applyGuillemetSpacing method tests in tests/unit/models/typography_processor_test.go"

# Launch integration tests together (T009-T013):
Task: "Basic quote processing test in tests/integration/french_quotes_test.go"
Task: "Punctuation spacing test in tests/integration/french_punctuation_test.go"
Task: "NNBSP preservation test in tests/integration/nnbsp_preservation_test.go"
Task: "Mixed content processing test in tests/integration/mixed_french_content_test.go"
Task: "Error handling integration test in tests/integration/typography_errors_test.go"
```

## Key Implementation Notes

### Unicode Character Handling
- **NNBSP**: U+202F (narrow non-breaking space) for French typography output
- **Regular Space**: U+0020 (detect and replace in French contexts)
- **Quotes**: U+00AB `«`, U+00BB `»` (preserve, add spacing)
- **Punctuation**: U+003A `:`, U+003B `;`, U+0021 `!`, U+003F `?`

### Processing Rules (from contracts/)
1. **Quote Processing**: `«text»` → `«{NNBSP}text{NNBSP}»`
2. **Punctuation Processing**: `text:` → `text{NNBSP}:`
3. **NNBSP Preservation**: Detect existing NNBSP, avoid duplication
4. **Mixed Content**: Apply all rules consistently across file content

### Test Scenarios (from quickstart.md)
- Basic quote processing with NNBSP insertion
- Punctuation spacing validation
- Existing NNBSP preservation
- Complex mixed French content
- Error handling for invalid input

### Integration Points
- Enhance existing `internal/models/typography_processor.go` methods
- Improve `applyFrenchTypography` and `applyGuillemetSpacing` functions
- Work through existing `applyTypography` function in `cmd/ankiprep/main.go`
- Maintain existing CLI interface (no command changes)
- Process all text content in TSV/CSV files through existing pipeline

## Notes
- [P] tasks target different files and can run in parallel
- All tests must fail before implementation begins (TDD)
- Follow existing ankiprep code patterns and Go conventions
- Enhance existing functionality rather than creating new services
- Custom implementation approach (no external French typography libraries)