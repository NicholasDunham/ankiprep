# CLI Interface Contract

## Command Structure

### Basic Command
```bash
ankiprep <input-files...> [options]
```

**Required Parameters**:
- `input-files`: One or more CSV/TSV file paths

**Optional Parameters**:
- `-o, --output <path>`: Output file path  
- `-f, --french`: Add thin spaces before French punctuation (:;!?)
- `-q, --smart-quotes`: Convert straight quotes to curly quotes
- `-s, --skip-duplicates`: Remove entries with identical content
- `-k, --keep-header`: Preserve the first row of CSV files
- `-v, --verbose`: Enable verbose output

### Standard Commands
```bash
ankiprep --help     # Show help text
ankiprep --version  # Show version information
```

## Input Validation

### Required Validations
- At least one input file must be provided
- All input files must exist and be readable
- Input files must have .csv or .tsv extensions
- If --output specified, output directory must be writable

### Error Handling
- Invalid file paths → Exit code 1, stderr message
- Unreadable files → Exit code 1, stderr message  
- Invalid flags → Exit code 1, stderr message with usage
- No input files → Exit code 1, stderr message with usage

## Output Contract

### Success Case
- Exit code: 0
- Stdout: Processing status messages (if -v/--verbose)
- Stderr: Empty
- File system: Output file created at specified or default location

### Failure Cases
- Exit code: Non-zero (1 for user errors, 2 for system errors)
- Stdout: Empty or partial progress messages
- Stderr: Error description and guidance
- File system: No output file created or partially created file cleaned up

## Behavior Preservation

### Must Maintain Identical Behavior For
- All command-line flag parsing and validation
- Help text format and content
- Error messages and exit codes  
- File input/output format and encoding (UTF-8)
- CSV/TSV parsing and output generation
- French typography processing with -f flag
- Smart quotes processing with -q flag  
- Duplicate detection and removal with -s flag
- Header handling with -k flag
- Verbose output format with -v flag

### Must Support Same Use Cases
- Single file processing: `ankiprep input.csv`
- Multiple files: `ankiprep file1.csv file2.csv -o combined.csv`
- Flag combinations: `ankiprep -fqvs input1.csv input2.csv -o output.csv`
- Output to specific location: `ankiprep input.csv -o /path/to/output.csv`
- Processing with all formatting options: `ankiprep --french --smart-quotes input.csv`

## Non-Functional Requirements

### Performance
- Sub-second response for typical file sizes (< 10MB)
- Graceful handling of large files with memory monitoring
- No performance regression compared to current implementation

### Reliability  
- Identical output for identical inputs
- Consistent behavior across runs
- Proper error handling and resource cleanup
- No data corruption or loss during processing

### Compatibility
- Cross-platform operation (Linux, macOS, Windows)
- UTF-8 encoding support maintained
- Backward compatibility with existing command usage patterns