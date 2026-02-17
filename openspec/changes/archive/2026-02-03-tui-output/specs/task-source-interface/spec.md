## ADDED Requirements

### Requirement: Task listing returns all tasks
The system SHALL provide a method to list all tasks with their current status.

#### Scenario: List all tasks
- **WHEN** Driver calls TaskSource.List()
- **THEN** System returns all tasks including closed, blocked, and pending

#### Scenario: Task status included
- **WHEN** List() returns tasks
- **THEN** Each task includes status field (open, in_progress, blocked, closed)

### Requirement: Beads client implements List
The system SHALL implement List() in BeadsClient using bd CLI.

#### Scenario: List via bd list
- **WHEN** BeadsClient.List() is called
- **THEN** System executes `bd list --json -n 0 --all` and parses JSON array response

#### Scenario: List includes dependencies
- **WHEN** BeadsClient.List() returns tasks
- **THEN** Each task includes blockers array for dependency information
