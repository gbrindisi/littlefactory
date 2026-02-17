## ADDED Requirements

### Requirement: File watcher detects progress.md changes
The system SHALL watch the progress file and notify the TUI when it changes.

#### Scenario: Progress file updated
- **WHEN** progress.md is written to disk
- **THEN** TUI receives a file change notification within 100ms

#### Scenario: Progress file created
- **WHEN** progress.md is created (did not exist before)
- **THEN** TUI receives a file change notification and loads initial content

### Requirement: File watcher uses configured state directory
The system SHALL construct the progress file path from the configured state directory.

#### Scenario: Custom state directory
- **WHEN** config specifies `state_dir: custom-state`
- **THEN** TUI watches `<project-root>/custom-state/progress.md`

#### Scenario: Default state directory
- **WHEN** config uses default state directory
- **THEN** TUI watches `<project-root>/.littlefactory/progress.md`

### Requirement: File watcher handles missing file gracefully
The system SHALL handle the case where progress.md does not yet exist.

#### Scenario: File does not exist at startup
- **WHEN** TUI starts and progress.md does not exist
- **THEN** TUI displays placeholder message and waits for file creation

#### Scenario: File deleted during run
- **WHEN** progress.md is deleted while TUI is running
- **THEN** TUI displays placeholder message and continues watching for recreation
