# Research: Remove Duplicate CSV Header in Anki Output

## Technology Decisions

### CSV Processing Approach
- **Decision**: Extend existing CSV processing in `internal/app/processor.go`
- **Rationale**: 
  - Leverages established patterns in the codebase
  - Minimal changes required to existing architecture
  - Maintains consistency with current CSV handling
- **Alternatives considered**: 
  - New separate CSV processor module (rejected: adds unnecessary complexity)
  - Modify existing typography service directly (rejected: violates single responsibility)

### Command-Line Flag Implementation
- **Decision**: Use github.com/spf13/cobra flag system with `--keep-header`/`-k`
- **Rationale**:
  - Consistent with existing CLI patterns in the codebase
  - Cobra provides standard POSIX flag handling
  - Boolean flag is intuitive for header preservation control
- **Alternatives considered**:
  - Positional arguments (rejected: less discoverable)
  - Configuration file option (rejected: overkill for single flag)
  - `--no-remove-header` (rejected: double-negative is confusing)

### Header Detection Strategy
- **Decision**: Remove first row by default, preserve with flag
- **Rationale**:
  - Most CSV files exported from applications include headers
  - Default behavior matches common user expectation
  - Flag provides escape hatch for edge cases
- **Alternatives considered**:
  - Auto-detection based on content patterns (rejected: too complex and error-prone)
  - Always require user specification (rejected: poor UX for common case)

### Error Handling Pattern
- **Decision**: Immediate exit with descriptive error messages and non-zero codes
- **Rationale**:
  - Follows POSIX conventions for CLI tools
  - Consistent with existing error handling in codebase
  - Enables proper shell scripting and automation
- **Alternatives considered**:
  - Warning-only approach (rejected: could lead to corrupted output)
  - Interactive prompts (rejected: breaks scriptability)

### Integration Points
- **Decision**: Modify existing processor workflow in `internal/app/processor.go`
- **Rationale**:
  - Single point of CSV processing modification
  - Maintains existing service separation
  - Minimal impact on other components
- **Alternatives considered**:
  - New middleware layer (rejected: adds complexity)
  - Service-level modification (rejected: mixed concerns)

## Implementation Patterns

### Existing Architecture Analysis
- Current CSV processing flows through `internal/app/processor.go`
- Uses `encoding/csv` standard library for file operations
- Typography service handles text transformations separately
- CLI commands defined in `cmd/ankiprep/` using cobra framework

### Flag Integration Pattern
- Add `keepHeader` boolean flag to root command
- Pass flag value through to processor component
- Modify processor to conditionally skip first row based on flag

### Error Propagation Pattern
- Wrap CSV reading errors with context
- Return structured errors with descriptive messages  
- CLI layer converts errors to appropriate exit codes

## Testing Strategy

### Unit Test Coverage
- Test header removal in default case
- Test header preservation with flag
- Test error handling for malformed CSV
- Test edge cases (empty files, single row files)

### Integration Test Coverage
- End-to-end CLI command testing
- Flag parsing and validation
- Output file verification
- Error message validation

### Backward Compatibility
- Existing functionality must remain unchanged
- New flag is additive only
- Default behavior change is the intended fix