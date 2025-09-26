# Processing Interface Contract

## File Processing Pipeline

### Input Processing
```go
type InputFile struct {
    Path     string
    Format   FileFormat  // CSV or TSV
    Encoding string      // Must be UTF-8
    Content  [][]string  // Parsed CSV data
}
```

### Processing Configuration
```go
type ProcessingConfig struct {
    French         bool   // Add thin spaces before French punctuation
    SmartQuotes    bool   // Convert straight quotes to curly quotes
    SkipDuplicates bool   // Remove duplicate entries
    KeepHeader     bool   // Preserve first row
    Verbose        bool   // Enable verbose logging
    Output         string // Output file path
}
```

## Typography Processing Contract

### French Typography Rules
- Input: Text with standard punctuation
- Processing: Add thin space (U+2009) before `:`, `;`, `!`, `?`
- Output: Text with French typography formatting
- Constraint: Must preserve all other characters unchanged

### Smart Quotes Processing
- Input: Text with straight quotes `"` and `'`
- Processing: Convert to curly quotes `"`, `"`, `'`, `'`
- Output: Text with typographically correct quotes
- Constraint: Must handle nested quotes and quote context correctly

## Duplicate Detection Contract

### Duplicate Criteria
- Comparison: Exact match of all field values after processing
- Behavior: Keep first occurrence, remove subsequent duplicates
- Ordering: Maintain stable sort order for non-duplicates

### Processing Requirements
- Must process duplicates AFTER typography processing
- Must respect SkipDuplicates flag (skip detection when false)
- Must maintain original row ordering for retained entries

## Output Generation Contract

### File Format Requirements
```go
type OutputFile struct {
    Path     string    // Target file path
    Format   string    // Always CSV for Anki compatibility
    Encoding string    // Always UTF-8
    Data     [][]string // Processed rows
}
```

### Anki Compatibility
- CSV format with comma separators
- Proper escaping of commas, quotes, and newlines
- UTF-8 encoding without BOM
- Consistent field ordering (Front, Back, additional fields)

## Validation Contracts

### Pre-processing Validation
- File existence and readability
- Valid CSV/TSV format structure
- UTF-8 encoding verification
- Minimum column count (≥2) validation

### Post-processing Validation
- Output file creation success
- Data integrity verification (row count, field preservation)
- Encoding validation (UTF-8 output)
- Format compliance (Anki-compatible CSV)

## Error Handling Contract

### Error Categories
```go
type ErrorCategory int
const (
    UserError   ErrorCategory = 1  // Invalid input, missing files
    SystemError ErrorCategory = 2  // I/O failures, permission issues
)
```

### Error Response Requirements
- Specific error messages with actionable guidance
- Appropriate exit codes (1 for user errors, 2 for system errors)
- Error output to stderr, not stdout
- Partial file cleanup on failure

## Memory Management Contract

### Large File Handling
- Stream processing for files >10MB
- Memory monitoring with configurable limits
- Graceful degradation or chunked processing
- Progress reporting for long operations

### Resource Cleanup
- Proper file handle closure
- Memory deallocation for large datasets  
- Temporary file cleanup on errors
- Goroutine cleanup (if used)