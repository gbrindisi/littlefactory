## ADDED Requirements

### Requirement: Driver claims task before iteration
The system SHALL mark task as in_progress before starting agent execution.

#### Scenario: Task claimed on iteration start
- **WHEN** Driver starts an iteration with a task
- **THEN** System calls TaskSource.Claim(taskID) before running agent

#### Scenario: Claim persists immediately
- **WHEN** Task is claimed
- **THEN** Status change is persisted to storage before agent starts

### Requirement: Driver marks task done on success
The system SHALL mark task as done when agent exits successfully.

#### Scenario: Successful iteration marks done
- **WHEN** Agent completes with exit code 0
- **THEN** System calls TaskSource.Close(taskID, "Completed")

### Requirement: Driver resets task on failure
The system SHALL reset task to todo when agent fails.

#### Scenario: Failed iteration resets task
- **WHEN** Agent completes with non-zero exit code
- **THEN** System calls TaskSource.Reset(taskID)

#### Scenario: Timeout resets task
- **WHEN** Agent exceeds iteration timeout
- **THEN** System calls TaskSource.Reset(taskID)
