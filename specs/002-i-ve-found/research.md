# Research: Cloze Deletion Colon Exception

**Date**: 2025-09-24  
**Feature**: Fix French typography processing for Anki cloze deletion blocks

## Research Questions Addressed

### 1. Cloze Deletion Pattern Recognition

**Decision**: Use regex pattern `\{\{c\d+::.*?\}\}` for basic detection, with additional parsing for nested content
**Rationale**: 
- Anki cloze syntax is well-defined: `{{c#::text}}` or `{{c#::text::hint}}`  
- Regex provides efficient initial detection
- Additional parsing handles nested French typography within content
**Alternatives Considered**: 
- Simple string matching (rejected: cannot handle nested content)
- Full parser (rejected: overkill for well-defined syntax)

### 2. French Typography Integration

**Decision**: Extend existing French typography service with pre-processing step to mask cloze syntax colons
**Rationale**:
- Preserves existing French typography logic and patterns
- Allows selective masking of only syntax colons while processing content colons normally  
- Maintains backward compatibility with existing functionality
**Alternatives Considered**:
- Separate cloze-aware typography service (rejected: code duplication)
- Post-processing approach (rejected: harder to distinguish syntax vs content colons)

### 3. Error Handling for Malformed Cloze Blocks

**Decision**: Log warning and fall back to standard French typography processing for the entire field
**Rationale**: 
- Graceful degradation prevents data loss
- Warning provides visibility for user correction
- Consistent with existing error handling patterns in the codebase
**Alternatives Considered**:
- Skip typography processing entirely (rejected: loses valuable typography improvements)
- Attempt partial processing (rejected: complexity risk for edge cases)

### 4. Multiple Cloze Block Processing

**Decision**: Process each valid cloze block independently using iterative pattern matching
**Rationale**:
- Anki supports multiple cloze blocks in single field
- Independent processing ensures consistent behavior per block
- Aligns with user expectation from clarification session
**Alternatives Considered**:
- Process only first block (rejected: limits user functionality)
- Batch processing (rejected: doesn't handle different block patterns cleanly)

### 5. Implementation Architecture

**Decision**: Add `ClozeDeletionHandler` to existing French typography service pipeline
**Rationale**:
- Follows existing service-oriented architecture
- Enables unit testing of cloze detection logic
- Integrates cleanly with existing typography flow
**Alternatives Considered**:
- Inline implementation in main typography function (rejected: reduces testability)
- Standalone utility (rejected: breaks service cohesion)

## Technical Specifications

### Cloze Block Structure
- **Basic pattern**: `{{c1::word}}` 
- **With hint**: `{{c1::word::hint}}`
- **Multiple numbers**: `{{c1::text}} {{c2::more}}`
- **Nested content**: `{{c1::phrase « with : punctuation »}}`

### Colon Handling Rules
- **Cloze syntax colons** (`::` between `{{c#` and content, between content and hint): NO NNBSP
- **Content colons** (within text or hint content): Apply normal French typography (ADD NNBSP)
- **Outside colons** (not within any `{{}}` block): Apply normal French typography

### Integration Points
- Extend `internal/services/typography.go` (inferred from existing project structure)
- Add unit tests in `tests/unit/services/`
- Add integration tests using CSV fixtures with cloze blocks
- Update existing French typography integration tests to verify no regression

## Implementation Approach

1. **Phase 1**: Add cloze detection and parsing logic with comprehensive unit tests
2. **Phase 2**: Integrate with existing French typography pipeline 
3. **Phase 3**: Add integration tests with real CSV fixtures
4. **Phase 4**: Performance testing to ensure sub-second processing maintained

**Estimated Complexity**: Medium - well-defined requirements, existing architecture to extend, clear test scenarios