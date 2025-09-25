# Processing Interface Contract

## Processor Interface

### ProcessCSV Function
```go
func ProcessCSV(inputPath, outputPath string, options ProcessingOptions) (*ProcessingResult, error)
```

**Input Parameters:**
- `inputPath` (string): Absolute path to input CSV file
- `outputPath` (string): Absolute path for output Anki CSV file  
- `options` (ProcessingOptions): Configuration for processing behavior

**Return Values:**
- `*ProcessingResult`: Success information or nil on error
- `error`: Detailed error information or nil on success

**Behavior Contract:**
- MUST validate input file exists and is readable
- MUST validate output directory is writable
- MUST read entire CSV file into memory for processing
- MUST apply header removal logic based on `options.KeepHeader`
- MUST generate Anki metadata headers in output
- MUST preserve original data row content and order
- MUST return descriptive errors for all failure modes
- MUST NOT modify input file
- MUST create output file atomically (temp file + rename)

**Error Conditions:**
- Returns `FileNotFoundError` if input file missing
- Returns `FilePermissionError` if cannot read input or write output
- Returns `CSVFormatError` if input is malformed CSV
- Returns `OutputWriteError` if cannot complete output file creation

### ProcessingOptions Structure
```go
type ProcessingOptions struct {
    KeepHeader bool // Default: false
}
```

**Field Contracts:**
- `KeepHeader=false`: Remove first row of CSV input
- `KeepHeader=true`: Preserve first row as data

### ProcessingResult Structure  
```go
type ProcessingResult struct {
    ProcessedRows int
    SkippedHeader bool
    OutputPath    string
}
```

**Field Contracts:**
- `ProcessedRows`: Count of data rows written (excluding header if removed)
- `SkippedHeader`: True if first row was identified and skipped as header
- `OutputPath`: Absolute path where output file was created

## CSV Format Contract

### Input Requirements
- MUST be valid CSV according to RFC 4180
- MUST use comma (`,`) as field separator  
- MUST handle quoted fields containing commas and newlines
- MUST support UTF-8 encoding
- MAY have inconsistent column counts (handled gracefully)

### Output Format
- MUST start with Anki metadata headers:
  ```
  #separator:comma
  #html:true  
  #columns:Text,Extra,Grammar_Notes
  ```
- MUST follow with data rows in CSV format
- MUST preserve original field quoting and escaping
- MUST NOT include original CSV header row (unless `KeepHeader=true`)

### Header Detection Logic
```
IF KeepHeader == false:
    Skip first row in output
    Set SkippedHeader = true
ELSE:
    Include all rows in output  
    Set SkippedHeader = false
```

## Error Handling Contract

### Error Types
All errors MUST implement Go error interface and provide descriptive messages.

### Error Messages
- MUST include file paths in error messages
- MUST include specific failure reason
- MUST be suitable for display to end users
- MUST follow format: "Error: {action} '{path}': {reason}"

### Error Propagation
- Low-level errors (file I/O, CSV parsing) MUST be wrapped with context
- Wrapped errors MUST preserve original error information
- Top-level errors MUST be actionable by users