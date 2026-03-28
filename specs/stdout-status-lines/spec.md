# stdout-status-lines

## What It Does
Prints human-readable status lines to stdout during driver execution, showing iteration progress, task info, and a final run summary. These lines provide real-time feedback without requiring a TUI.

## Requirements
### Requirement: Iteration start status line
The system SHALL print a status line to stdout when each iteration begins.

#### Scenario: Iteration start output
- **WHEN** an iteration starts for a task
- **THEN** system prints `[N/MAX] Starting: <task-title> (<task-id>)` to stdout

### Requirement: Iteration complete status line
The system SHALL print a status line to stdout when each iteration completes.

#### Scenario: Successful iteration output
- **WHEN** an iteration completes successfully
- **THEN** system prints `[N/MAX] Completed` to stdout

#### Scenario: Failed iteration output
- **WHEN** an iteration fails
- **THEN** system prints `[N/MAX] Failed: <error>` to stdout

#### Scenario: Timed out iteration output
- **WHEN** an iteration times out
- **THEN** system prints `[N/MAX] Timed out` to stdout

### Requirement: Run summary status line
The system SHALL print a summary line to stdout when the run finishes.

#### Scenario: Completed run summary
- **WHEN** the run completes
- **THEN** system prints `Run complete: <status> (<iterations-run>/<max-iterations> iterations)` to stdout

#### Scenario: Cancelled run summary
- **WHEN** the run is cancelled via signal
- **THEN** system prints `Run cancelled` to stdout

## Boundaries

## Gotchas
