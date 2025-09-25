# Feature Specification: Remove Duplicate CSV Header in Anki Output

**Feature Branch**: `003-the-program-is`  
**Created**: September 24, 2025  
**Status**: Draft  
**Input**: User description: "The program is correctly adding the three Anki-specific headers, but it's also retaining the typical CSV header. The original header line should be removed once the #columns header line is in place."

## Execution Flow (main)
```
1. Parse user description from Input
   → User reports duplicate header issue in Anki CSV output
2. Extract key concepts from description
   → Actors: ankiprep CLI tool, users processing CSV files
   → Actions: CSV processing, header management, Anki format conversion
   → Data: CSV files with headers, Anki-specific metadata headers
   → Constraints: Anki format requires specific headers, original CSV headers become redundant
3. For each unclear aspect:
   → All aspects are clear from user description and provided file examples
4. Fill User Scenarios & Testing section
   → Clear user flow: process CSV file, expect clean Anki output without duplicate headers
5. Generate Functional Requirements
   → Each requirement is testable with input/output file comparison
6. Identify Key Entities
   → CSV files, headers, Anki metadata
7. Run Review Checklist
   → No clarifications needed, implementation details avoided
8. Return: SUCCESS (spec ready for planning)
```

---

## Clarifications

### Session 2025-09-24
- Q: When the system encounters a CSV file where the first row contains data that happens to match typical column names (like "Text,Extra,Grammar_Notes"), how should it determine whether this is a header row to remove or actual flashcard data to preserve? → A: Remove first line by default, flag to keep
- Q: What should happen when the system encounters an error during CSV processing (e.g., malformed CSV, file not found, or write permission issues)? → A: Exit immediately with error code and descriptive message
- Q: What should be the name and format of the command-line flag to preserve the first row? → A: --keep-header / -k (short form)

---

## User Scenarios & Testing *(mandatory)*

### Primary User Story
A user processes a CSV file containing flashcard data using ankiprep. The tool correctly adds the required Anki-specific headers (`#separator:comma`, `#html:true`, `#columns:Text,Extra,Grammar_Notes`) but incorrectly retains the original CSV header row (`Text,Extra,Grammar_Notes`). This creates redundant information in the output file, where the column names appear twice - once in the `#columns` metadata header and again as a data row.

### Acceptance Scenarios
1. **Given** a CSV file with header row "Text,Extra,Grammar_Notes" and data rows, **When** processing with ankiprep, **Then** the output should contain only the Anki metadata headers and data rows, with the original CSV header row removed
2. **Given** any CSV file with a header row, **When** ankiprep processes it and adds `#columns` metadata, **Then** the original header row should be omitted from the output
3. **Given** a properly processed Anki CSV file, **When** imported into Anki, **Then** the import should work correctly without any duplicate or malformed header issues

### Edge Cases
- What happens when the CSV has no header row? (System should handle gracefully)
- What happens when the `#columns` header doesn't match the original CSV headers? (System should use the detected structure)
- How does the system distinguish between header rows and data rows that happen to contain the same text?

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: System MUST remove the first row of CSV input by default when generating Anki-formatted output
- **FR-002**: System MUST provide a `--keep-header` (or `-k`) command-line flag to preserve the first row when it contains data rather than headers
- **FR-003**: System MUST retain all Anki-specific metadata headers (`#separator`, `#html`, `#columns`)
- **FR-004**: System MUST preserve all data rows from the original CSV file (excluding removed header row)
- **FR-005**: System MUST ensure the `#columns` header accurately reflects the column structure without duplication in the data
- **FR-006**: System MUST exit immediately with non-zero error code and descriptive message when encountering processing errors

### Key Entities *(include if feature involves data)*
- **CSV Input File**: Original file with potential header row and data rows containing flashcard information
- **Anki Metadata Headers**: Special comment-style headers required by Anki (`#separator`, `#html`, `#columns`)
- **CSV Header Row**: First row of original CSV that defines column names, should be removed in output
- **Data Rows**: Actual flashcard content rows that must be preserved in the output

---

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

---

## Execution Status
*Updated by main() during processing*

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [x] Review checklist passed
