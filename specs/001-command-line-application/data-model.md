# Data Model: CSV to Anki Processor

**Date**: 2025-09-24  
**Feature**: CSV to Anki Processor CLI  
**Branch**: 001-command-line-application

## Core Entities

### InputFile
Represents a source CSV/TSV file to be processed.

**Fields**:
- `Path string` - Absolute file path
- `Separator rune` - Field separator (comma or tab)
- `Headers []string` - Column header names
- `Records [][]string` - Data rows (excluding header)
- `Encoding string` - Character encoding (UTF-8 only)

**Validation Rules**:
- Path must exist and be readable
- File must be valid CSV/TSV format
- Must use comma or tab separator only
- Must contain at least one data row
- Must use UTF-8 encoding

**State Transitions**:
- Created → Validated → Parsed → Ready

### DataEntry
Represents a single row of data with field values.

**Fields**:
- `Values map[string]string` - Column name to value mapping
- `Source string` - Originating file path
- `LineNumber int` - Original line number in source file

**Validation Rules**:
- Values map must not be empty
- All values must be valid UTF-8 strings
- Source must reference valid input file

**Uniqueness Rules**:
- Exact match across all field values determines duplicates
- Case-sensitive comparison
- Empty values participate in duplicate detection

### OutputFile
Represents the final merged and formatted CSV output.

**Fields**:
- `Path string` - Output file path (always .csv extension)
- `Headers []string` - Union of all input file headers
- `Records []DataEntry` - Deduplicated and merged data entries
- `AnkiHeaders []string` - Anki-specific header lines

**Validation Rules**:
- Path must be writable location
- Headers must be union of all input headers
- Records must be deduplicated
- AnkiHeaders must include separator, html, and columns directives

**State Transitions**:
- Created → Headers Merged → Records Merged → Deduplicated → Formatted → Written

### ProcessingReport
Summary of processing actions and statistics.

**Fields**:
- `InputFiles []string` - List of processed input file paths
- `TotalInputRecords int` - Count of records before deduplication
- `DuplicatesRemoved int` - Count of duplicate records removed
- `OutputRecords int` - Final count of records in output
- `ProcessingTime time.Duration` - Total processing time
- `Errors []string` - List of any processing errors

**Validation Rules**:
- TotalInputRecords >= OutputRecords
- DuplicatesRemoved = TotalInputRecords - OutputRecords
- ProcessingTime must be positive duration

### TypographyProcessor
Handles text formatting transformations.

**Fields**:
- `FrenchMode bool` - Whether French typography rules are enabled
- `ConvertSmartQuotes bool` - Whether to convert straight quotes to smart quotes

**Processing Rules**:
- **Smart Quotes**: Convert " to " and ", convert ' to '
- **French Mode**: Add NNBSP before :;?! and inside « »
- **HTML Line Breaks**: Convert embedded newlines to `<br>` tags
- **Preserve HTML**: Maintain existing HTML tags in content

## Relationships

```
InputFile (1) ----contains----> (N) DataEntry
DataEntry (N) ----merges-to----> (N) OutputFile.Records
ProcessingReport (1) ----references----> (N) InputFile
TypographyProcessor (1) ----transforms----> (N) DataEntry.Values
```

## Data Flow

1. **Input Phase**: Multiple InputFile entities are created and validated
2. **Parsing Phase**: Each InputFile generates DataEntry instances
3. **Merging Phase**: Headers are merged using union operation
4. **Deduplication Phase**: Duplicate DataEntry instances are removed
5. **Typography Phase**: TypographyProcessor transforms text content
6. **Output Phase**: OutputFile is generated with Anki headers
7. **Reporting Phase**: ProcessingReport summarizes the operation

## Constraints

- **Memory**: All data must fit in available system memory
- **Encoding**: UTF-8 only for all text processing
- **Separators**: Comma and tab only for input files
- **Output Format**: Always comma-separated CSV with Anki headers
- **File Size**: Progress indicators required for files >10MB

## Error Handling

- **Invalid File Format**: Immediate termination with clear error message
- **Encoding Issues**: Reject files with non-UTF-8 encoding
- **Memory Exhaustion**: Clean shutdown with partial file cleanup
- **Write Failures**: Remove partial output files on error
- **Malformed Data**: Skip empty rows, error on other malformed content

## Next Steps

Ready for contract generation and test scenario creation based on this data model.