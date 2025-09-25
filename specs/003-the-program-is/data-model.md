# Data Model: Remove Duplicate CSV Header in Anki Output

## Core Entities

### ProcessingOptions
Represents the user-configurable options for CSV processing behavior.

**Fields:**
- `KeepHeader` (bool): Whether to preserve the first row of the CSV file
  - Default: `false` (remove header by default)
  - Controlled by `--keep-header`/`-k` CLI flag
- `InputFile` (string): Path to the input CSV file
- `OutputFile` (string): Path for the output Anki-formatted CSV file

**Validation Rules:**
- `InputFile` must exist and be readable
- `OutputFile` directory must be writable
- `InputFile` and `OutputFile` must not be the same path

**State Transitions:**
- Created → Validated → Used for processing

### CSVRecord
Represents a single row from the CSV input file.

**Fields:**
- `Values` ([]string): The column values for this row
- `IsHeader` (bool): Whether this row is identified as a header row
- `LineNumber` (int): Original line number in the source file (for error reporting)

**Validation Rules:**
- `Values` must not be nil
- `LineNumber` must be positive
- Number of values should be consistent across records (for well-formed CSV)

### ProcessingResult
Represents the outcome of CSV processing operation.

**Fields:**
- `ProcessedRows` (int): Number of data rows processed (excluding header if removed)
- `SkippedHeader` (bool): Whether the first row was skipped as a header
- `OutputPath` (string): Path where the processed file was written
- `Error` (error): Any error that occurred during processing

**Validation Rules:**
- `ProcessedRows` must be non-negative
- `OutputPath` must be set if processing succeeded
- `Error` should be nil on successful processing

## Data Flow

### Input Processing
1. Parse CLI flags into `ProcessingOptions`
2. Validate `ProcessingOptions` (file existence, permissions)
3. Open and read CSV file into `[]CSVRecord`
4. Determine header handling based on `KeepHeader` flag

### Header Processing Logic
```
IF KeepHeader == false AND file has rows:
    Mark first CSVRecord as IsHeader = true
    Skip first record in output processing
ELSE:
    Process all records as data
```

### Output Generation
1. Write Anki metadata headers (`#separator:comma`, `#html:true`, `#columns:...`)
2. Write data records (excluding any marked as header)
3. Close output file
4. Return `ProcessingResult`

## Error Handling

### Error Types
- `FileNotFoundError`: Input file doesn't exist
- `FilePermissionError`: Cannot read input or write output
- `CSVFormatError`: Malformed CSV content
- `OutputWriteError`: Cannot write to output file

### Error Propagation
- All errors bubble up to CLI layer with context
- CLI layer converts to appropriate exit codes (non-zero)
- Error messages include file paths and specific failure reason

## Relationships

- `ProcessingOptions` → used by → CSV Processing Logic
- `CSVRecord[]` → transformed by → Header Processing Logic
- Header Processing Logic → produces → Anki Output + `ProcessingResult`
- Errors → propagated through → All layers

## Storage Considerations

- No persistent storage required
- All processing is file-to-file transformation
- Temporary in-memory storage for CSV records during processing
- Memory usage scales linearly with input file size