## ADDED Requirements

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
