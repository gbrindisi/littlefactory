## REMOVED Requirements

### Requirement: Beads client implementation
**Reason**: Replaced by JSONTaskSource - bd dependency is being removed entirely
**Migration**: Use JSONTaskSource with `.littlefactory/tasks.json`

### Requirement: Beads client implements List
**Reason**: Replaced by JSONTaskSource - bd dependency is being removed entirely
**Migration**: Use JSONTaskSource with `.littlefactory/tasks.json`

## MODIFIED Requirements

### Requirement: TaskSource interface defines task operations
The system SHALL define a TaskSource interface that abstracts task management operations.

#### Scenario: Ready tasks retrieval
- **WHEN** Driver calls TaskSource.Ready()
- **THEN** System returns list of tasks available for work (status `todo`)

#### Scenario: Task details retrieval
- **WHEN** Driver calls TaskSource.Show(id)
- **THEN** System returns full task details including ID, title, description, and status

#### Scenario: Task completion
- **WHEN** Driver calls TaskSource.Close(id, reason)
- **THEN** System marks task status as `done`

#### Scenario: Task claim
- **WHEN** Driver calls TaskSource.Claim(id)
- **THEN** System marks task status as `in_progress`

#### Scenario: Task reset
- **WHEN** Driver calls TaskSource.Reset(id)
- **THEN** System marks task status as `todo`

### Requirement: Task listing returns all tasks
The system SHALL provide a method to list all tasks with their current status.

#### Scenario: List all tasks
- **WHEN** Driver calls TaskSource.List()
- **THEN** System returns all tasks regardless of status

#### Scenario: Task status included
- **WHEN** List() returns tasks
- **THEN** Each task includes status field (`todo`, `in_progress`, `done`)
