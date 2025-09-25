# Quickstart Guide: CSV to Anki Processor

**Date**: 2025-09-24  
**Feature**: CSV to Anki Processor CLI  
**Branch**: 001-command-line-application

## Prerequisites

- Go 1.21+ installed
- Git for source control
- Test CSV files for validation

## Installation

```bash
# Clone repository
git clone <repository-url>
cd ankiprep

# Build the tool
go build -o ankiprep ./cmd/ankiprep

# Verify installation
./ankiprep --version
./ankiprep --help
```

## Quick Test

### 1. Prepare Test Data

Create test CSV file `test-input.csv`:
```csv
Front,Back,Tags
What is Go?,A programming language,programming
What is CSV?,Comma-separated values,data
```

### 2. Basic Processing

```bash
# Process single file
./ankiprep test-input.csv

# Expected output file: merged.csv
# Expected console output:
# Processing 1 input files...
# Merging headers: found 3 unique columns  
# Processing records: 2 total entries
# Removing duplicates: found 0 duplicates
# Writing output to merged.csv
# Done. Processed 2 unique entries in <time>
```

### 3. Verify Output

Expected `merged.csv` content:
```csv
#separator:Comma
#html:true
#columns:Front,Back,Tags
What is Go?,A programming language,programming
What is CSV?,Comma-separated values,data
```

### 4. Test French Typography

Create French test file `french-test.csv`:
```csv
Front,Back
Question : "Que dis-tu ?",Réponse : « Bonjour ! »
```

```bash
# Process with French typography
./ankiprep -f -o french-output.csv french-test.csv

# Verify NNBSP added before punctuation and inside quotes
```

### 5. Test Multiple Files

Create second file `test-input2.csv`:
```csv
Front,Back,Category
What is Anki?,Spaced repetition software,tools
What is Go?,A programming language,programming
```

```bash
# Process multiple files (note duplicate detection)
./ankiprep -o combined.csv test-input.csv test-input2.csv

# Expected: 3 unique entries (1 duplicate removed)
```

## Validation Tests

### Test 1: Single File Processing
- **Input**: One CSV file with 2-3 rows
- **Expected**: Valid Anki-formatted output with all rows
- **Verify**: Anki headers present, data intact

### Test 2: Multiple File Merging  
- **Input**: Two CSV files with overlapping columns
- **Expected**: Union of headers, all unique data
- **Verify**: Column merging, duplicate detection

### Test 3: French Typography
- **Input**: CSV with French text and punctuation
- **Expected**: NNBSP applied correctly
- **Verify**: Proper spacing around : ; ! ? and « »

### Test 4: Smart Quotes Conversion
- **Input**: CSV with straight quotes "text" and 'apostrophes'  
- **Expected**: Converted to "text" and 'apostrophes'
- **Verify**: Typography transformation

### Test 5: Error Handling
- **Input**: Non-existent file or malformed CSV
- **Expected**: Clear error message, non-zero exit code
- **Verify**: Proper error reporting

### Test 6: Large File Processing
- **Input**: CSV file >10MB or processing >5 seconds
- **Expected**: Progress indicator displayed
- **Verify**: User feedback during processing

### Test 7: Mixed Separators
- **Input**: CSV and TSV files together
- **Expected**: Both processed correctly
- **Verify**: Auto-detection of separators

### Test 8: Multiline Content
- **Input**: CSV with embedded newlines in cells
- **Expected**: Newlines converted to `<br>` tags
- **Verify**: HTML formatting preserved

## Development Validation

Run the complete test suite:

```bash
# Unit tests
go test ./...

# Integration tests  
go test ./tests/integration/...

# CLI contract tests
go test ./tests/contract/...

# Performance tests
go test -bench=. ./tests/performance/...
```

## Success Criteria

✅ All test scenarios pass  
✅ No errors or warnings in test output  
✅ Generated CSV files import successfully into Anki  
✅ Performance meets targets (sub-second for typical files)  
✅ Memory usage remains reasonable for large files  
✅ Cross-platform compatibility verified  

## Next Steps

After quickstart validation:
1. Review generated code structure
2. Run full test suite
3. Perform manual integration testing with real Anki import
4. Verify cross-platform builds (Linux, macOS, Windows)
5. Document any additional usage patterns discovered

## Troubleshooting

### Common Issues

- **"File not found"**: Check file paths and permissions
- **"Invalid format"**: Ensure files are valid CSV/TSV with UTF-8 encoding  
- **"Permission denied"**: Check write permissions for output directory
- **Memory issues**: Use smaller batch sizes for very large files

### Debug Mode

```bash
# Enable verbose output
./ankiprep --verbose input.csv

# Verbose with French typography
./ankiprep -v -f input.csv
```