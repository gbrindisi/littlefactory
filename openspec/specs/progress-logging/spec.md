# progress-logging Specification

## Purpose
TBD - created by archiving change go-rewrite-littlefactory. Update Purpose after archive.
## Requirements
### Requirement: Progress file initialization
The system SHALL initialize or reuse progress file at `<state_dir>/progress.md`.

#### Scenario: New progress file creation
- **WHEN** Progress file does not exist at run start
- **THEN** System creates `<state_dir>/progress.md` with header "# Little Factory Progress Log", "**Started:** <timestamp>", and separator "---"

#### Scenario: Existing progress file reuse
- **WHEN** Progress file exists at run start
- **THEN** System appends to existing file without overwriting

### Requirement: Iteration session logging
The system SHALL append session info to progress file after each iteration in markdown format.

#### Scenario: Session info format
- **WHEN** Iteration completes
- **THEN** System appends markdown block with "## Iteration N", "- **Task:** <task-id>", "- **Status:** <status>", and separator "---"

#### Scenario: Missing session path handling
- **WHEN** Iteration has no session_path
- **THEN** System still appends iteration block but omits session line

### Requirement: Append-only semantics
The system SHALL only append to progress file, never truncate or replace.

#### Scenario: Multiple runs preserve history
- **WHEN** Multiple runs execute against same project
- **THEN** Progress file accumulates all iteration logs chronologically

### Requirement: Progress functions receive config
The system SHALL pass full config to progress functions for state directory access.

#### Scenario: InitProgressFile signature
- **WHEN** InitProgressFile is called
- **THEN** Function receives projectRoot and *config.Config parameters

#### Scenario: AppendSessionToProgress signature
- **WHEN** AppendSessionToProgress is called
- **THEN** Function receives projectRoot, *config.Config, and iteration parameters

#### Scenario: ProgressFilePath signature
- **WHEN** ProgressFilePath is called
- **THEN** Function receives projectRoot and *config.Config parameters

