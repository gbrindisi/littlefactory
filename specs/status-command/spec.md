# status-command

## What It Does
Provides a CLI command that displays task progress for the current project or across all worktrees. Supports summary and detailed views, and can target a specific change via the --change flag.

## Requirements
### Requirement: Status command shows task progress
The system SHALL provide a `status` command that displays task progress.

#### Scenario: Status in current directory
- **WHEN** User runs `littlefactory status`
- **THEN** System shows task progress for the current directory's tasks.json

#### Scenario: Status with change flag
- **WHEN** User runs `littlefactory status --change feature-a`
- **THEN** System shows task progress for `openspec/changes/feature-a/tasks.json`

### Requirement: Status shows summary format
The system SHALL display status in a concise summary format.

#### Scenario: Summary format output
- **WHEN** User runs `littlefactory status -c feature-a`
- **THEN** System outputs `feature-a: 3/7 done` format showing completed vs total tasks

#### Scenario: In-progress task shown
- **WHEN** A task has status `in_progress`
- **THEN** System appends `(in_progress: "<task title>")` to the summary

#### Scenario: All tasks complete
- **WHEN** All tasks have status `done`
- **THEN** System shows `[complete]` indicator

### Requirement: Status shows all worktrees
The system SHALL show status for all worktrees when `--all` flag is used.

#### Scenario: Status all flag
- **WHEN** User runs `littlefactory status --all`
- **THEN** System lists status for each worktree that has a tasks.json

#### Scenario: Status all discovers worktrees
- **WHEN** User runs `littlefactory status --all`
- **THEN** System uses `git worktree list` to find all worktrees and checks each for tasks.json

### Requirement: Status detailed view
The system SHALL show detailed task list when viewing a specific change.

#### Scenario: Detailed task list
- **WHEN** User runs `littlefactory status -c feature-a --verbose`
- **THEN** System shows each task with status indicator: `[done]`, `[in_progress]`, or `[todo]`

## Boundaries

## Gotchas
