# ankiprep

A Go CLI tool for processing CSV/TSV files into Anki-compatible format.

## Features

- Process CSV and TSV files
- Smart quotes and French typography formatting
- Duplicate detection and removal
- Memory monitoring for large files

## Installation

```bash
# Build the binary
go build -o ankiprep ./cmd/ankiprep
```

## Usage

```bash
# Basic usage
./ankiprep input.csv

# Multiple files with output path
./ankiprep file1.csv file2.csv -o combined.csv

# French typography and smart quotes
./ankiprep -f -q input.csv

# Skip duplicate detection for faster processing
./ankiprep -s input.csv

# Verbose output
./ankiprep -v input.csv

# Combine multiple short flags
./ankiprep -fqvs input1.csv input2.csv -o output.csv

# Keep original CSV headers in output
./ankiprep --keep-header input.csv
```

### Command Options

- `-o, --output`: Specify output file path
- `-f, --french`: Add thin spaces before French punctuation (:;!?)  
- `-q, --smart-quotes`: Convert straight quotes to curly quotes
- `-s, --skip-duplicates`: Remove entries with identical content
- `-k, --keep-header`: Preserve the first row of CSV files (default: remove header)
- `-v, --verbose`: Enable verbose output

## Input Format

CSV files should have at least two columns with a header row:

```csv
Front,Back
Hello,Bonjour
Goodbye,Au revoir
```

Supports CSV (`.csv`) and TSV (`.tsv`) files with UTF-8 encoding.

## Output

Creates Anki-compatible CSV files with proper escaping and UTF-8 encoding.

## Development

```bash
# Run tests
go test ./...

# Build binary
go build -o ankiprep ./cmd/ankiprep

# Run without building
go run ./cmd/ankiprep --help
