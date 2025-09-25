# CLI Interface Contract

**Tool**: ankiprep  
**Version**: 1.0.0  
**Date**: 2025-09-24

## Command Signature

```bash
ankiprep [options] <input-files...>
```

## Options

| Flag | Long Form | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `-h` | `--help` | none | - | Show help message and exit |
| `-V` | `--version` | none | - | Show version and exit |
| `-f` | `--french` | none | false | Apply French typography rules |
| `-o` | `--output` | string | `merged.csv` | Output file path |
| `-v` | `--verbose` | none | false | Enable verbose output |

## Arguments

- `input-files`: One or more paths to CSV/TSV files (required)

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | Invalid arguments or options |
| 2 | Input file not found or unreadable |
| 3 | Invalid file format (not CSV/TSV or wrong separator) |
| 4 | Processing error (malformed data, encoding issues) |
| 5 | Output error (write permission, disk space) |

## Standard Output Format

### Success Case
```
Processing 3 input files...
Merging headers: found 8 unique columns
Processing records: 1,247 total entries
Removing duplicates: found 23 duplicates
Applying typography formatting...
Writing output to merged.csv
Done. Processed 1,224 unique entries in 0.15 seconds
```

### Progress Indicator (large files)
```
Processing 3 input files...
[=====>              ] 25% (245/1000 records)
```

## Standard Error Format

### Error Messages
```bash
Error: File 'input.csv' not found
Error: Invalid file format in 'data.txt' - expected CSV or TSV
Error: Malformed CSV data at line 42 in 'file.csv'
Error: Cannot write to 'output.csv' - permission denied
```

### Warnings (non-fatal)
```bash
Warning: Empty rows detected in 'file.csv' - skipping
```

## Examples

### Basic Usage
```bash
# Process single file
ankiprep input.csv

# Process multiple files
ankiprep file1.csv file2.tsv file3.csv

# Specify output file
ankiprep -o flashcards.csv input1.csv input2.csv
```

### French Typography
```bash
# Apply French typography rules
ankiprep -f french_content.csv

# French mode with custom output
ankiprep --french --output french_cards.csv content.csv
```

## Output File Format

### Anki Headers
```csv
#separator:Comma
#html:true
#columns:Front,Back,Tags,Notes
```

### Data Format
```csv
What is the capital of France?,Paris,geography,European capitals
How do you say 'hello' in Spanish?,"Hola (pronounced: /ˈola/)",language,Spanish basics
```

## Validation Rules

### Input Validation
- All input files must exist and be readable
- Files must be valid CSV (comma-separated) or TSV (tab-separated)
- Files must use UTF-8 encoding
- Files must contain at least one header row and one data row

### Output Validation
- Output directory must be writable
- Output filename will always have .csv extension (auto-appended if missing)
- Partial output files are cleaned up on error

## Behavioral Contracts

1. **Deterministic Output**: Same inputs always produce same output
2. **Memory Bounded**: Tool will not exceed available system memory
3. **Atomic Output**: Output file is only written if processing succeeds completely
4. **Progress Feedback**: Long operations (>5s or >10MB) show progress
5. **Error Recovery**: Clean shutdown with descriptive error messages