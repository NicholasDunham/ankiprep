# Feature Specification: CSV to Anki Processor

**Feature Branch**: `001-command-line-application`  
**Created**: 2025-09-24  
**Status**: Draft  
**Input**: User description: "Command-line application that accepts one or more CSV files (exported from Google Sheets, Notion, or another application) and cleans, formats, and merges them to prepare for import into Anki."

---

## Clarifications

### Session 2025-09-24

- Q: What should be the duplicate detection strategy for flashcard entries? → A: Exact match on all fields (front, back, tags, etc.)
- Q: What content formatting should be preserved in flashcard data? → A: HTML tags preserved for rich formatting
- Q: What is the maximum file size the system should handle? → A: No explicit limit (memory-dependent)
- Q: When should progress indicators be shown during processing? → A: After 5 seconds of processing time
- Q: What character encoding support is required for input files? → A: UTF-8 only (most common modern standard)
- Q: How should malformed data be handled? → A: Quit with error except for empty rows which are silently dropped
- Q: How should column headers be merged from different files? → A: Union of all headers with empty cells for missing values
- Q: Should the program validate specific field requirements? → A: No assumptions about required fields, accept whatever is in input files
- Q: Should the program support French typography formatting? → A: Yes, via -f/--french flag applying NNBSP rules for quotes and punctuation
- Q: Should the program convert straight quotes to smart quotes? → A: Yes, convert " to " and ", and ' to ' (apostrophes only)
- Q: Should the output use Anki-specific headers? → A: Yes, use #separator:Comma, #html:true, and #columns headers instead of CSV header row
- Q: When processing fails partway through (e.g., out of memory, disk full), what should happen to any partial output file that was being written? → A: Delete any partial output file, leave no traces
- Q: What should the tool do if it encounters CSV files with different field separators (comma vs tab vs semicolon)? → A: Support comma and tab only, reject others
- Q: Should the tool preserve the original file extension in the output filename or always use .csv? → A: Always use .csv extension regardless of input
- Q: How should the tool handle CSV files that contain data with embedded newlines (multiline cell content)? → A: Replace newlines with HTML `<br>` tags for Anki
- Q: When the tool encounters very large files, what specific memory threshold should trigger the progress indicator? → A: 10 MB

## ⚡ Quick Guidelines

- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

## User Scenarios & Testing

### Primary User Story

A student or educator has data in multiple CSV files exported from different sources (Google Sheets, Notion, etc.) with inconsistent formatting, duplicate entries, and varying column structures. They need to combine and clean this data into a single, properly formatted CSV file that can be imported into Anki with the correct headers and formatting.

### Acceptance Scenarios

1. **Given** multiple CSV files with identical column headers, **When** user runs the processor with those files, **Then** the tool produces a single merged CSV maintaining the same column structure
2. **Given** multiple CSV files with different column headers, **When** user runs the processor with those files, **Then** the tool produces a single merged CSV with union of all column headers, filling missing values with empty cells
3. **Given** CSV files containing duplicate flashcard entries, **When** user processes the files, **Then** duplicates are identified and removed using exact match on all fields (front, back, tags, metadata)
4. **Given** CSV files with completely empty rows, **When** user processes the files, **Then** empty rows are silently dropped and processing continues normally
5. **Given** CSV files with malformed data or invalid CSV/TSV format, **When** user processes the files, **Then** the tool produces a clear error message and quits immediately
6. **Given** a single CSV file with formatting issues, **When** user processes the file, **Then** the tool cleans and reformats it to proper CSV format
7. **Given** incompatible input files, **When** user attempts to process them, **Then** the tool provides clear error messages about what's wrong and how to fix it
8. **Given** French flag is specified with valid CSV files, **When** user processes the files, **Then** the tool applies French typographic conventions with Narrow Non-Breaking Spaces (NNBSP) for quotes and punctuation
9. **Given** CSV files with straight quotes and apostrophes, **When** user processes the files, **Then** the tool converts them to smart quotes using English typographical rules
10. **Given** valid CSV files are processed, **When** user runs the processor, **Then** the output file includes Anki-specific headers (#separator, #html, #columns) instead of traditional CSV column headers

### Edge Cases

- What happens when input files have completely different column structures that can't be reconciled?
- How does the system handle very large CSV files that might not fit in memory?
- What happens when input CSV files use different character encodings?
- How does the system handle special characters or formatting in data content that might affect CSV parsing?
- How does the French typography processing handle mixed content with both French and other language conventions?
- How does the smart quotes conversion handle nested quotations or complex punctuation scenarios?
- How does the system handle column names that contain commas or special characters in Anki headers?

## Requirements

### Functional Requirements

- **FR-001**: System MUST accept one or more CSV file paths as command-line arguments
- **FR-002**: System MUST read and parse CSV files exported from Google Sheets, Notion, and other common applications
- **FR-003**: System MUST identify and remove duplicate flashcard entries across all input files using exact match on all fields (front, back, tags, metadata)
- **FR-004**: System MUST silently drop completely empty rows and quit with clear error message for any other malformed data
- **FR-005**: System MUST merge multiple CSV files into a single output file
- **FR-006**: System MUST format output as valid CSV with Anki-specific headers for direct import compatibility
- **FR-007**: System MUST provide progress feedback during processing of large files
- **FR-008**: System MUST generate a summary report of processing actions (duplicates removed, entries cleaned, etc.)
- **FR-009**: System MUST handle different CSV column structures by creating union of all column headers and filling missing values with empty cells
- **FR-010**: System MUST merge files with identical column structures as-is without modification
- **FR-011**: System MUST preserve all column headers and data from input files without validation of specific field requirements
- **FR-012**: System MUST support UTF-8 character encoding for all input and output files
- **FR-013**: System MUST preserve HTML tags in flashcard content for rich formatting support
- **FR-014**: System MUST validate that all input files use comma-separated (CSV) or tab-separated (TSV) format and quit with error if other separators are detected
- **FR-015**: System MUST support French typography flag (-f, --french) to apply French typographic conventions
- **FR-016**: When French flag is enabled, system MUST ensure angled quotes (« and ») are separated from content by Narrow Non-Breaking Space (NNBSP)
- **FR-017**: When French flag is enabled, system MUST ensure double punctuation marks (!, ?, :, ;) are preceded by NNBSP
- **FR-018**: System MUST convert straight double quotes (") to smart quotes (" and ") using English typographical rules
- **FR-019**: System MUST convert straight apostrophes (') to smart apostrophes (') assuming all single quotes represent apostrophes, not quotation marks
- **FR-020**: System MUST include Anki-specific separator header (#separator:Comma) as first line of output file
- **FR-021**: System MUST include Anki-specific HTML header (#html:true) as second line of output file  
- **FR-022**: System MUST include Anki-specific columns header (#columns:...) as third line of output file, listing all column names
- **FR-023**: System MUST exclude traditional CSV column header row from output file, with data starting immediately after Anki headers
- **FR-024**: System MUST delete any partial output file if processing fails partway through (out of memory, disk full, etc.) to leave no traces
- **FR-025**: System MUST always use .csv file extension for output file regardless of input file extensions
- **FR-026**: System MUST replace embedded newlines in cell content with HTML `<br>` tags to maintain Anki compatibility

### Performance Requirements

- **PR-001**: System MUST handle files of any size within available system memory constraints
- **PR-002**: System MUST provide progress indicators for operations taking longer than 5 seconds or when processing files larger than 10 MB

### Usability Requirements

- **UR-001**: System MUST provide help text explaining usage, supported file formats, and French typography flag
- **UR-002**: System MUST output clear error messages when files cannot be processed
- **UR-003**: System MUST support standard CLI conventions (--help, --version, -f/--french flags)

### Key Entities

- **CSV Input File**: Represents a source file containing tabular data, may have varying column structures and data quality
- **Data Entry**: Individual row of data with values corresponding to the file's column headers
- **Processing Report**: Summary of actions taken during processing including statistics on duplicates, cleaning operations, and errors
- **Output File**: Final merged and cleaned CSV file with Anki-specific headers ready for direct import into Anki

## Review & Acceptance Checklist

*GATE: Automated checks run during main() execution*

### Content Quality

- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

### Requirement Completeness

- [x] No [NEEDS CLARIFICATION] markers remain
- [x] Requirements are testable and unambiguous  
- [x] Success criteria are measurable
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

## Execution Status

*Updated by main() during processing*

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [x] Review checklist passed
