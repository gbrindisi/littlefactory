## ADDED Requirements

### Requirement: Project root detection via .beads directory
The system SHALL locate project root by searching for .beads directory.

#### Scenario: .beads in current directory
- **WHEN** Current working directory contains .beads/ subdirectory
- **THEN** System uses current directory as project root

#### Scenario: .beads in parent directory
- **WHEN** Current directory does not have .beads/ but parent directory does
- **THEN** System walks up directory tree and uses first directory containing .beads/ as project root

#### Scenario: No .beads found
- **WHEN** No .beads/ directory found in current or parent directories
- **THEN** System uses current directory as project root (default fallback)

### Requirement: Tasks directory convention
The system SHALL use tasks/ subdirectory within project root for all artifacts.

#### Scenario: Progress file location
- **WHEN** System writes progress file
- **THEN** File is created at <project-root>/tasks/progress.txt

#### Scenario: Metadata file location
- **WHEN** System writes metadata
- **THEN** File is created at <project-root>/tasks/run_metadata.json

#### Scenario: Template override location
- **WHEN** System checks for template override
- **THEN** System looks at <project-root>/tasks/CLAUDE.md

#### Scenario: Tasks directory creation
- **WHEN** tasks/ directory does not exist
- **THEN** System creates it (mkdir -p behavior)
