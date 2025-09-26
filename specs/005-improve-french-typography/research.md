# Research: French Typography Libraries and Implementation

**Feature**: Improve French Typography NNBSP Handling  
**Date**: September 26, 2025

## Research Tasks

### 1. French Typography Library Evaluation
**Question**: Are there existing Go libraries specifically designed for French typography rules?

**Research Findings**:

#### Go Typography Libraries Survey
- **golang.org/x/text**: Standard library extension with unicode normalization and text processing
  - Provides unicode/norm for normalization
  - Has language-specific text processing capabilities
  - No specific French typography rules built-in
  
- **Third-party Options**:
  - Limited specific French typography libraries in Go ecosystem
  - Most typography libraries focus on general text processing
  - No mature, battle-tested French typography libraries identified

#### Decision: Custom Implementation
- **What was chosen**: Develop custom French typography logic within existing codebase structure
- **Rationale**: No suitable existing library found that specifically handles French NNBSP rules for quotes and punctuation. Custom implementation allows precise control over the specific requirements.
- **Alternatives considered**: 
  - Using general text processing libraries (insufficient for specific French rules)
  - Porting rules from other languages (adds complexity without clear benefit)
  - Using unicode/norm only (doesn't handle French-specific spacing rules)

### 2. Unicode Character Handling
**Question**: What are the specific Unicode characters involved in French typography processing?

**Research Findings**:

#### Character Mapping
- **Regular Space**: U+0020 (what users type with spacebar)
- **Narrow Non-Breaking Space (NNBSP)**: U+202F (what French typography requires)
- **Opening Angle Quote**: U+00AB (« - left-pointing double angle quotation mark)
- **Closing Angle Quote**: U+00BB (» - right-pointing double angle quotation mark)
- **Double Punctuation**: : (U+003A), ; (U+003B), ! (U+0021), ? (U+003F)

#### Implementation Approach
- Use Go's built-in string handling and rune processing
- Leverage existing golang.org/x/text/unicode/norm for consistent text normalization
- Implement pattern matching for quote and punctuation detection

### 3. Integration with Existing Codebase
**Question**: How should French typography processing integrate with AnkiPrep's current architecture?

**Research Findings**:

#### Current Architecture Analysis
- Single-user CLI application with simplified architecture
- Text processing currently handled in main command logic
- Existing typography processing using golang.org/x/text/unicode/norm
- File processing works with CSV/TSV text content

#### Integration Strategy
- **What was chosen**: Extend existing typography processing within main command flow
- **Rationale**: Maintains architectural simplicity, leverages existing text processing pipeline
- **Alternatives considered**:
  - Separate typography service (unnecessary complexity for single-user CLI)
  - External typography processor (violates simplicity principle)
  - Plugin architecture (over-engineering for this scope)

### 4. Testing Strategy
**Question**: How should French typography rules be tested effectively?

**Research Findings**:

#### Testing Approach
- **Unit Tests**: Test individual typography transformations with specific input/output pairs
- **Integration Tests**: Test complete CLI workflows with French text processing
- **Test Data**: Create comprehensive test cases covering:
  - All quote scenarios (no space, regular space, existing NNBSP)
  - All punctuation scenarios (:, ;, !, ? with various spacing)
  - Edge cases (multiple punctuation, nested quotes, line boundaries)

#### Test Organization
- Add to existing test structure under `tests/unit/` for typography logic
- Add to existing integration tests for CLI workflow validation
- Maintain TDD approach: write failing tests first

## Summary

**Primary Technical Decisions**:
1. **Custom Implementation**: No suitable Go library exists, will implement French typography rules directly
2. **Unicode Handling**: Use Go's native string/rune processing with specific character codes
3. **Architecture Integration**: Extend existing typography processing in main command flow
4. **Testing Strategy**: Comprehensive unit and integration tests covering all scenarios

**No Remaining NEEDS CLARIFICATION**: All technical unknowns resolved for implementation planning.