# CLI Interface Contract

**Command**: `ankiprep`  
**Feature**: French typography processing with cloze deletion support  
**Version**: Extended from existing CLI

## Command Interface

### Basic Usage
```bash
ankiprep [flags] input.csv [output.csv]
```

### Flags (Extended)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--french` | `-f` | bool | false | Enable French typography rules with cloze deletion support |
| `--smart-quotes` | `-q` | bool | false | Enable smart quote processing |
| `--skip-duplicates` | `-s` | bool | false | Skip duplicate detection |
| `--help` | `-h` | bool | false | Show help information |
| `--version` | `-v` | bool | false | Show version information |

## Behavior Changes

### French Typography Enhancement
When `--french` flag is used, the tool now:

1. **Detects Anki cloze deletion blocks** using pattern `{{c#::...}}`
2. **Preserves cloze syntax colons** (no NNBSP added to `::` within `{{}}`)
3. **Applies French rules to content colons** (NNBSP added to colons within cloze content)
4. **Maintains existing behavior** for all other typography rules
5. **Logs warnings** for malformed cloze blocks to stderr

### Backward Compatibility
- All existing functionality preserved
- Files without cloze blocks process identically to previous version
- Existing command-line flags and behavior unchanged
- Exit codes unchanged

## Input/Output Contract

### Input Requirements
- **File format**: CSV with UTF-8 encoding
- **File size**: Up to 10MB efficiently processed
- **Cloze format**: Anki-compatible `{{c#::content}}` or `{{c#::content::hint}}`

### Output Guarantees  
- **Format preservation**: Same CSV structure and encoding as input
- **Data integrity**: No data loss, only typography modifications
- **Cloze preservation**: All valid cloze block structure maintained exactly

### Error Handling
- **Malformed cloze blocks**: Warning to stderr, continue processing
- **File errors**: Clear error message, non-zero exit code
- **Processing errors**: Graceful degradation, preserve original content

## Examples

### Basic Cloze Processing
```bash
# Input CSV contains: "Text,Extra\n\"Je {{c1::mange}} une pomme : délicieuse !\",\"\""
$ ankiprep --french input.csv output.csv
Processing input.csv with French typography...
Processed 1 rows, found 1 cloze blocks
Written to output.csv

# Output preserves cloze syntax, adds NNBSP to content colon
# "Je {{c1::mange}} une pomme : délicieuse !",""
```

### Multiple Flags
```bash  
$ ankiprep -fq input.csv output.csv
Processing input.csv with French typography and smart quotes...
Processed 100 rows, found 45 cloze blocks  
Written to output.csv
```

### Malformed Cloze Warning
```bash
$ ankiprep --french problematic.csv fixed.csv
Processing problematic.csv with French typography...
Warning: Malformed cloze deletion block at row 5, column 1: "{{c1::incomplete"
Processed 200 rows, found 38 cloze blocks, 1 warning
Written to fixed.csv
```

### Help Output (Enhanced)
```bash
$ ankiprep --help
ankiprep - Process CSV files for Anki with typography improvements

USAGE:
    ankiprep [FLAGS] <input.csv> [output.csv]

FLAGS:
    -f, --french            Enable French typography rules (includes cloze deletion support)  
    -q, --smart-quotes      Enable smart quote processing
    -s, --skip-duplicates   Skip duplicate detection
    -h, --help              Show this help message
    -v, --version           Show version information

EXAMPLES:
    ankiprep --french cards.csv                    # Process with French typography  
    ankiprep -fq input.csv output.csv             # French + smart quotes
    ankiprep --french --skip-duplicates cards.csv  # French with no duplicate checking

CLOZE DELETION SUPPORT:
    When --french is enabled, Anki cloze deletions are processed specially:
    - Cloze syntax colons (::) are preserved without modification
    - Content within cloze blocks follows French typography rules  
    - Malformed cloze blocks generate warnings but continue processing

For more information, visit: https://github.com/user/ankiprep
```

## Performance Contract

- **Processing time**: Sub-second for files up to 1MB
- **Memory usage**: < 100MB peak for largest supported files
- **Progress reporting**: Visual indicator for files > 1000 rows
- **Responsiveness**: User can interrupt with Ctrl+C

## Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | Processing completed successfully |
| 1 | File error | Input file not found or not readable |
| 2 | Format error | Invalid CSV format or encoding |
| 3 | Processing error | Unexpected error during processing |
| 130 | User interrupt | Process interrupted by user (Ctrl+C) |

## Logging

### Standard Output (stdout)
- Processing progress messages
- Summary statistics (rows processed, cloze blocks found)
- Success confirmation

### Standard Error (stderr)  
- Warning messages for malformed cloze blocks
- Error messages for failures
- Debug information (if enabled)

### Example Log Output
```
Processing input.csv with French typography...
Row 50: Warning - Malformed cloze deletion block: "{{c1::incomplete"
Row 75: Warning - Malformed cloze deletion block: "{{c2::missing::close"
Processed 100 rows, found 23 cloze blocks, 2 warnings
Written to output.csv
```