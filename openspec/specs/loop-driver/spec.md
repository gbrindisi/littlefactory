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

### Requirement: Driver accepts change name parameter
The system SHALL accept a change name to determine the task source path.

#### Scenario: Driver with change name
- **WHEN** Driver is initialized with change name `feature-a`
- **THEN** Driver uses `openspec/changes/feature-a/tasks.json` as task source

#### Scenario: Driver without change name
- **WHEN** Driver is initialized without change name
- **THEN** Driver uses default `<state_dir>/tasks.json` as task source

### Requirement: Driver supports worktree workspace
The system SHALL support running in a worktree directory.

#### Scenario: Driver runs in worktree
- **WHEN** Driver is configured with a worktree path
- **THEN** Driver changes to worktree directory before executing agent

#### Scenario: Driver creates worktree when requested
- **WHEN** Driver is configured with worktree creation enabled
- **THEN** Driver creates worktree before running agent loop
