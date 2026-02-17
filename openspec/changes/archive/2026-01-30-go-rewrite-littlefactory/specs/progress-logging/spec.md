## ADDED Requirements

### Requirement: Progress file initialization
The system SHALL initialize or reuse progress file at tasks/progress.txt.

#### Scenario: New progress file creation
- **WHEN** Progress file does not exist at run start
- **THEN** System creates tasks/progress.txt with header "# Ciccio Progress Log", "Started: <timestamp>", and separator "---"

#### Scenario: Existing progress file reuse
- **WHEN** Progress file exists at run start
- **THEN** System appends to existing file without overwriting

### Requirement: Iteration session logging
The system SHALL append session info to progress file after each iteration.

#### Scenario: Session info format
- **WHEN** Iteration completes
- **THEN** System appends block with "## Ciccio Iteration N", task ID, status, session path, and separator "---"

#### Scenario: Missing session path handling
- **WHEN** Iteration has no session_path
- **THEN** System still appends iteration block but omits session line

### Requirement: Append-only semantics
The system SHALL only append to progress file, never truncate or replace.

#### Scenario: Multiple runs preserve history
- **WHEN** Multiple runs execute against same project
- **THEN** Progress file accumulates all iteration logs chronologically
