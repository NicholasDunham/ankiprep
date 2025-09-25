# ankiprep Examples

This directory contains example input files to demonstrate ankiprep's capabilities.

## Files

### basic.csv
Simple vocabulary list with English-French translations.

```bash
ankiprep examples/basic.csv
```

### french.csv
French text with punctuation and special characters.

```bash
# Process with French typography rules
ankiprep --french examples/french.csv

# Process with smart quotes
ankiprep --smart-quotes examples/french.csv

# Apply both French rules and smart quotes
ankiprep --french --smart-quotes examples/french.csv
```

### multilingual.csv
Mixed language content with Unicode characters from various scripts.

```bash
ankiprep examples/multilingual.csv
```

### trivia.csv
Three-column format for trivia questions with hints.

```bash
ankiprep -o trivia_flashcards.csv examples/trivia.csv
```

### spanish_vocabulary.tsv
Tab-separated values file with pronunciation guides.

```bash
# TSV files are automatically detected
ankiprep examples/spanish_vocabulary.tsv
```

## Processing Multiple Files

You can process multiple example files at once:

```bash
# Process all CSV examples
ankiprep examples/basic.csv examples/french.csv examples/multilingual.csv

# Process all files (CSV and TSV)
ankiprep examples/*.csv examples/*.tsv
```

## Expected Output

All examples will generate an `anki_output.csv` file (unless you specify a different output file with `-o`) that is compatible with Anki's import format.

## Performance Testing

For performance testing, you can generate larger files based on these examples:

```bash
# Create a larger test file by repeating basic.csv
for i in {1..1000}; do cat examples/basic.csv | tail -n +2; done >> large_test.csv
# Add header
echo "front,back" | cat - large_test.csv > temp && mv temp large_test.csv
```