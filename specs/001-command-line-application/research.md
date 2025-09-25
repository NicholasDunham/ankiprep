# Research: CSV to Anki Processor

**Date**: 2025-09-24  
**Feature**: CSV to Anki Processor CLI  
**Branch**: 001-command-line-application

## Research Tasks Completed

### Go Standard Library for CSV Processing

**Decision**: Use Go's built-in `encoding/csv` package for all CSV parsing and writing  
**Rationale**: 
- Handles RFC 4180 CSV standard properly
- Built-in support for different separators (comma, tab)
- Proper handling of quoted fields and embedded newlines
- No external dependencies required
- Battle-tested and performant
- Third-party libraries don't provide significant edge case handling benefits for this use case

**Alternatives considered**: 
- Third-party CSV libraries (gocsv, csvutil) - provide struct mapping but unnecessary for our map-based approach
- Manual CSV parsing - rejected due to complexity and edge case handling

**Re-evaluation**: Standard library remains the best choice as CSV processing is straightforward and well-handled by `encoding/csv`

### Unicode and Typography Processing

**Decision**: Use external library `golang.org/x/text/unicode/norm` plus custom typography processing  
**Rationale**:
- Unicode normalization is complex with many edge cases (NFC, NFD, NFKC, NFKD)
- Professional typography handling requires understanding of language-specific rules
- Custom French typography rules (NNBSP placement) need precise implementation
- Smart quotes conversion has context-dependent rules that benefit from proven algorithms
- The `golang.org/x/text` package provides robust Unicode handling

**Alternatives considered**:
- Pure standard library approach - rejected due to Unicode normalization complexity
- Full typography library (e.g., typography-focused packages) - overkill for specific French/English rules needed

**Implementation approach**:
- Use `golang.org/x/text/unicode/norm` for Unicode normalization
- Implement custom French typography rules for NNBSP placement around :;?! and inside « »
- Implement smart quotes conversion with context awareness
- Use `unicode/utf8` for basic string operations

### Command-Line Interface Design

**Decision**: Use `github.com/spf13/cobra` for advanced CLI functionality  
**Rationale**:
- Provides superior help text generation with examples and usage patterns
- Better error handling and validation for command arguments
- Subcommand support for potential future extensions
- Industry standard for Go CLI applications (used by kubectl, docker, etc.)
- Automatic shell completion generation
- Better POSIX compliance and flag handling edge cases

**Alternatives considered**:
- Standard library `flag` package - adequate but limited help text and validation capabilities
- `github.com/urfave/cli` - good alternative but cobra has better ecosystem support

**Re-evaluation**: The significant benefit in user experience, error handling, and extensibility justifies the external dependency

### File Processing Strategy

**Decision**: In-memory processing with streaming for large files  
**Rationale**:
- Most CSV files for personal use are <10MB and fit comfortably in memory
- Allows for duplicate detection across all files
- Column merging requires full dataset visibility
- Progress indicators can be added for large files

**Alternatives considered**:
- Pure streaming approach - rejected due to duplicate detection requirements
- Database-backed processing - rejected as overkill for use case

### Error Handling Strategy

**Decision**: Fail-fast with descriptive error messages  
**Rationale**:
- Users need clear feedback about malformed files
- Early termination prevents partial/corrupted output
- Aligns with Unix tool philosophy

**Alternatives considered**:
- Best-effort processing - rejected due to data integrity concerns
- Interactive error resolution - rejected as non-scriptable

### Output Format Strategy

**Decision**: Always output comma-separated CSV with Anki headers  
**Rationale**:
- Anki expects comma-separated format for import
- Anki-specific headers (#separator, #html, #columns) ensure proper import
- Consistent output format regardless of input mix

**Alternatives considered**:
- Preserve input separator in output - rejected due to Anki compatibility
- Generic CSV output - rejected due to Anki import requirements

## Technical Decisions Summary

1. **Language**: Go 1.21+ (latest stable version)
2. **Dependencies**: 
   - Core: Standard library (encoding/csv, os, fmt, strings, unicode)
   - CLI: `github.com/spf13/cobra` for enhanced CLI experience
   - Typography: `golang.org/x/text/unicode/norm` for Unicode normalization
3. **Architecture**: Single-binary CLI tool
4. **Processing**: In-memory with progress indicators
5. **Testing**: Go standard testing with table-driven tests
6. **Distribution**: Single static binary per platform

## Performance Considerations

- **Memory usage**: Linear with input file size, suitable for typical use cases
- **CPU usage**: Minimal processing overhead, mostly I/O bound
- **Progress feedback**: 5-second threshold or 10MB file size trigger
- **Error recovery**: Clean shutdown with partial file cleanup

## Next Phase Requirements

All research complete, no outstanding unknowns. Ready for Phase 1 design and contract generation.