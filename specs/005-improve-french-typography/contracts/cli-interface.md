# CLI Command Contract

**Command**: ankiprep  
**Feature**: French Typography Processing  
**Version**: 1.0  
**Date**: September 26, 2025

## Command Interface

### Basic Usage

```bash
ankiprep input.tsv output.tsv
```

**Description**: Process TSV file with French typography rules applied automatically

**Input**:
- `input.tsv`: Source TSV file containing French text content
- `output.tsv`: Destination file for processed content

**Processing Behavior**:
- Apply French typography rules to all text content
- Preserve file structure and formatting
- Maintain header row if present
- Process all columns containing text data

### Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | File processed successfully |
| 1 | File Error | Input file not found or output cannot be written |
| 2 | Format Error | Invalid file format or encoding |
| 3 | Processing Error | Typography processing failed |

## Contract Tests

### Test: Basic File Processing
```bash
# Given: Input file with French text requiring typography
echo -e "Front\tBack\nBonjour\t«Salut»" > input.tsv

# When: Process file
ankiprep input.tsv output.tsv

# Then: Output contains proper French typography
# Expected output.tsv content:
# Front	Back
# Bonjour	« Salut »  # Using NNBSP (U+202F)

# Exit code: 0
```

### Test: Quote Processing
```bash
# Given: File with various quote formats
echo -e "French\tEnglish\n«bonjour»\thello\n« au revoir »\tgoodbye" > quotes.tsv

# When: Process file
ankiprep quotes.tsv processed.tsv

# Then: All quotes have proper NNBSP spacing
# Expected processed.tsv:
# French	English
# « bonjour »	hello      # NNBSP added
# « au revoir »	goodbye   # Regular space replaced with NNBSP

# Exit code: 0
```

### Test: Punctuation Processing  
```bash
# Given: File with French punctuation
echo -e "Question\tAnswer\nComment allez-vous?\tHow are you?\nTrès bien!\tVery well!" > punct.tsv

# When: Process file
ankiprep punct.tsv result.tsv

# Then: Punctuation has proper NNBSP spacing
# Expected result.tsv:
# Question	Answer
# Comment allez-vous ?	How are you?  # NNBSP before ?
# Très bien !	Very well!            # NNBSP before !

# Exit code: 0
```

### Test: Mixed Content Processing
```bash
# Given: File with mixed French typography needs
echo -e "Card\tContent\nGreeting\t«Bonjour: comment allez-vous?»" > mixed.tsv

# When: Process file  
ankiprep mixed.tsv final.tsv

# Then: All typography rules applied
# Expected final.tsv:
# Card	Content
# Greeting	« Bonjour : comment allez-vous ? »  # All NNBSP applied

# Exit code: 0
```

### Test: Preserve Existing NNBSP
```bash  
# Given: File with existing NNBSP characters
echo -e "Text\tNote\nAlready correct\t« Bonjour : comment ? »" > existing.tsv
# Note: Input uses NNBSP (U+202F) characters

# When: Process file
ankiprep existing.tsv output.tsv

# Then: Existing NNBSP preserved
# Expected output.tsv (unchanged):
# Text	Note  
# Already correct	« Bonjour : comment ? »  # NNBSP preserved

# Exit code: 0
```

### Test: Error - File Not Found
```bash
# Given: Non-existent input file

# When: Process file
ankiprep nonexistent.tsv output.tsv

# Then: Error reported and exit code 1
# Expected stderr: "Error: Input file 'nonexistent.tsv' not found"
# Exit code: 1
```

### Test: Error - Invalid UTF-8
```bash
# Given: File with invalid UTF-8 content
printf "Front\tBack\n\xff\xfe\tinvalid" > invalid.tsv

# When: Process file
ankiprep invalid.tsv output.tsv

# Then: Format error reported
# Expected stderr: "Error: Invalid UTF-8 encoding in input file"
# Exit code: 2
```

### Test: Error - Write Permission
```bash
# Given: Valid input file but read-only output directory
echo -e "Text\tNote\n«test»\tnote" > input.tsv
mkdir readonly_dir
chmod 444 readonly_dir

# When: Process file to read-only location
ankiprep input.tsv readonly_dir/output.tsv

# Then: Write error reported
# Expected stderr: "Error: Cannot write to output file"
# Exit code: 1

# Cleanup
chmod 755 readonly_dir && rm -rf readonly_dir
```

## Command Line Validation

### Valid Arguments
- Two positional arguments (input and output file paths)
- File extensions: `.tsv`, `.csv` (case insensitive)  
- Relative and absolute paths supported
- Unicode filenames supported

### Invalid Arguments
- Missing arguments: Error with usage message
- More than 2 arguments: Error with usage message
- Input same as output: Error to prevent data loss
- Invalid characters in filenames: System-dependent behavior

## Typography Processing Details

**Characters Handled**:
- **NNBSP**: U+202F (narrow non-breaking space)
- **Regular Space**: U+0020 (replaced with NNBSP in French contexts)
- **Quotes**: U+00AB `«`, U+00BB `»`
- **Punctuation**: U+003A `:`, U+003B `;`, U+0021 `!`, U+003F `?`

**Processing Rules**:
- Apply to all text content in all columns
- Preserve file structure (headers, column count, row count)
- Maintain other formatting and whitespace
- Process each cell independently