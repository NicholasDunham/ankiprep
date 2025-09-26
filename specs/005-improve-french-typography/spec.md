# Feature Specification: Improve French Typography NNBSP Handling

**Feature Branch**: `005-improve-french-typography`  
**Created**: September 26, 2025  
**Status**: Draft  
**Input**: User description: "Currently, if the vocabulary cards already use narrow non-breaking spaces (NNBSP) inside angled quotes, AnkiPrep still adds them again. This is unnecessary. Improve French typography handling to avoid duplicate NNBSP insertion for quotes and punctuation marks."

## Execution Flow (main)
```
1. Parse user description from Input
   → Feature clearly defined: Fix duplicate NNBSP insertion
2. Extract key concepts from description
   → Actors: AnkiPrep users with French vocabulary cards
   → Actions: Typography processing, NNBSP insertion/replacement
   → Data: Text content with French punctuation and quotes
   → Constraints: Follow French typography rules, avoid duplicates
3. For each unclear aspect:
   → All scenarios clearly specified by user
4. Fill User Scenarios & Testing section
   → User flow: Process French text with proper spacing
5. Generate Functional Requirements
   → Each requirement testable and specific
6. Identify Key Entities
   → Text content, punctuation marks, spacing characters
7. Run Review Checklist
   → No [NEEDS CLARIFICATION] markers
8. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines
- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

---

## Clarifications

### Session 2025-09-26
- Q: Should AnkiPrep fall back to custom implementation or skip French typography processing when no suitable library exists? → A: Fall back to custom implementation maintaining current behavior (decision made at development time, not runtime)
- Q: Which specific Unicode characters should be treated as full-width spaces versus standard spaces? → A: Both refer to regular space (U+0020)
- Q: Should typography rules apply to entire text or only to specifically identified French portions? → A: Apply French typography rules to entire text content

---

## User Scenarios & Testing

### Primary User Story
As a French language learner using AnkiPrep to process vocabulary cards, I want the typography processing to intelligently handle existing narrow non-breaking spaces (NNBSP) so that my cards follow proper French punctuation rules without unnecessary duplication of spacing characters.

### Acceptance Scenarios

#### Angled Quotation Marks (« »)
1. **Given** text with angled quotes containing no spaces (e.g., «bonjour»), **When** AnkiPrep processes the text, **Then** it adds NNBSPs after opening and before closing quotes (« bonjour »)

2. **Given** text with angled quotes containing regular spaces (e.g., « bonjour »), **When** AnkiPrep processes the text, **Then** it replaces regular spaces with NNBSPs (« bonjour »)

3. **Given** text with angled quotes already containing NNBSPs (e.g., « bonjour »), **When** AnkiPrep processes the text, **Then** it leaves the spacing unchanged

#### Double Punctuation Marks (:, ;, !, ?)
1. **Given** text with double punctuation marks with no preceding space (e.g., "Bonjour:"), **When** AnkiPrep processes the text, **Then** it adds a NNBSP before the punctuation mark ("Bonjour :")

2. **Given** text with double punctuation marks with regular spaces (e.g., "Bonjour :"), **When** AnkiPrep processes the text, **Then** it replaces the regular space with a NNBSP ("Bonjour :")

3. **Given** text with double punctuation marks already having NNBSPs (e.g., "Bonjour :"), **When** AnkiPrep processes the text, **Then** it leaves the spacing unchanged

### Edge Cases
- What happens when text contains mixed spacing types within the same document?
- How does the system handle nested quotations with different spacing?
- What occurs when punctuation marks appear at the beginning or end of lines?

## Requirements

### Functional Requirements
- **FR-001**: System MUST detect existing narrow non-breaking spaces (NNBSP) in French text before applying typography rules
- **FR-002**: System MUST add NNBSPs after opening angled quotes («) and before closing angled quotes (») when no spacing exists
- **FR-003**: System MUST replace regular spaces with NNBSPs around angled quotation marks
- **FR-004**: System MUST preserve existing NNBSPs around angled quotation marks without duplication
- **FR-005**: System MUST add NNBSPs before French double punctuation marks (:, ;, !, ?) when no space exists
- **FR-006**: System MUST replace regular spaces with NNBSPs before French double punctuation marks
- **FR-007**: System MUST preserve existing NNBSPs before French double punctuation marks without duplication
- **FR-008**: System MUST maintain compliance with standard French typography and punctuation rules
- **FR-009**: System MUST process text without changing the overall structure or complexity of the codebase
- **FR-010**: System MUST use existing or battle-tested typography libraries when available; if none exist, MUST fall back to custom implementation that maintains current behavior (library selection decision made at development time)
- **FR-011**: System MUST apply French typography rules to all text content (language-specific processing for mixed content to be addressed in future feature)

### Key Entities
- **Text Content**: The input text containing French vocabulary and phrases that require typography processing
- **Punctuation Marks**: French double punctuation characters (:, ;, !, ?) that require preceding NNBSP
- **Quotation Marks**: Angled French quotes (« ») that require internal NNBSP spacing
- **Spacing Characters**: Various space types including regular spaces (U+0020) and narrow non-breaking spaces (NNBSP/U+202F)

---

## Review & Acceptance Checklist

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

---

## Execution Status

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [x] Review checklist passed
