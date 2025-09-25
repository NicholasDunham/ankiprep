# CLI Interface Contract

## Command Interface

### Root Command
```
ankiprep [input.csv] [output.csv] [flags]
```

**Flags:**
- `--keep-header, -k`: Preserve the first row of the CSV (default: false)
- `--help, -h`: Show help information
- `--version`: Show version information

**Arguments:**
- `input.csv`: Path to input CSV file (required)
- `output.csv`: Path for output Anki CSV file (optional, defaults to input with .anki.csv extension)

**Exit Codes:**
- `0`: Success
- `1`: General error (file not found, permission denied)
- `2`: Invalid arguments or flags
- `3`: CSV format error

**Examples:**
```bash
# Default behavior - remove header
ankiprep flashcards.csv

# Preserve header row
ankiprep -k flashcards.csv flashcards.anki.csv

# Show help
ankiprep --help
```

## Error Messages

### File Errors
```
Error: cannot read input file 'nonexistent.csv': file not found
Error: cannot write to output file 'readonly.csv': permission denied
```

### Format Errors
```
Error: malformed CSV at line 5: unexpected quote character
Error: inconsistent column count at line 12: expected 3 columns, got 5
```

### Argument Errors
```
Error: input file required
Error: unknown flag '--invalid'
Usage: ankiprep [input.csv] [output.csv] [flags]
```

## Success Output

### Verbose Mode (default)
```
Processing flashcards.csv...
Removed header row: Text,Extra,Grammar_Notes
Processed 142 rows
Output written to flashcards.anki.csv
```

### With --keep-header flag
```
Processing flashcards.csv...
Preserved header row
Processed 143 rows (including header)
Output written to flashcards.anki.csv
```