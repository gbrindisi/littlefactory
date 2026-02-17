## ADDED Requirements

### Requirement: TaskSource interface defines task operations
The system SHALL define a TaskSource interface that abstracts task management operations.

#### Scenario: Ready tasks retrieval
- **WHEN** Driver calls TaskSource.Ready()
- **THEN** System returns list of tasks with no blockers

#### Scenario: Task details retrieval
- **WHEN** Driver calls TaskSource.Show(id)
- **THEN** System returns full task details including ID, title, description, and status

#### Scenario: Task completion
- **WHEN** Driver calls TaskSource.Close(id, reason)
- **THEN** System marks task as complete with given reason

#### Scenario: State persistence
- **WHEN** Driver calls TaskSource.Sync()
- **THEN** System persists pending state changes to storage

### Requirement: Beads client implementation
The system SHALL provide a beads implementation of the TaskSource interface.

#### Scenario: Beads CLI validation at startup
- **WHEN** System initializes BeadsClient
- **THEN** System verifies `bd` binary exists in PATH and returns error if not found

#### Scenario: Ready tasks via bd ready
- **WHEN** BeadsClient.Ready() is called
- **THEN** System executes `bd ready --json` and parses JSON array response

#### Scenario: Task details via bd show
- **WHEN** BeadsClient.Show(id) is called
- **THEN** System executes `bd show <id> --json` and extracts first element from array response

#### Scenario: Task closure via bd close
- **WHEN** BeadsClient.Close(id, reason) is called
- **THEN** System executes `bd close <id> --reason <reason>`

#### Scenario: Sync via bd sync
- **WHEN** BeadsClient.Sync() is called
- **THEN** System executes `bd sync` to persist database to JSONL
