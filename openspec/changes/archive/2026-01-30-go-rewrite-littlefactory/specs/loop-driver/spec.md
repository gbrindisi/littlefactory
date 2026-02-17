## ADDED Requirements

### Requirement: Sequential iteration execution
The system SHALL execute agent iterations sequentially up to max iterations limit.

#### Scenario: Normal iteration flow
- **WHEN** Driver runs with max_iterations=10
- **THEN** System executes iterations 1 through N sequentially until tasks complete or max reached

#### Scenario: Early completion
- **WHEN** All tasks are complete before max iterations
- **THEN** System exits with RunStatus.COMPLETED

#### Scenario: Max iterations reached
- **WHEN** Max iterations reached with tasks remaining
- **THEN** System exits with RunStatus.MAX_ITERATIONS

#### Scenario: User interruption
- **WHEN** User sends SIGINT (Ctrl+C) during execution
- **THEN** System gracefully stops and exits with RunStatus.ABORTED and exit code 130

### Requirement: Iteration timeout enforcement
The system SHALL enforce per-iteration timeout using context deadline.

#### Scenario: Iteration completes within timeout
- **WHEN** Agent completes before timeout (default 600s)
- **THEN** Iteration marked as COMPLETED

#### Scenario: Iteration exceeds timeout
- **WHEN** Agent exceeds timeout deadline
- **THEN** Iteration marked as TIMEOUT and execution continues to next iteration

### Requirement: Keep-going error strategy
The system SHALL continue to next iteration on failures.

#### Scenario: Failed iteration continues
- **WHEN** Agent returns non-zero exit code
- **THEN** Iteration marked as FAILED and system continues to next iteration

#### Scenario: Exception during iteration continues
- **WHEN** Exception occurs during agent execution
- **THEN** Iteration marked as FAILED with error message and system continues to next iteration

### Requirement: Visual progress output
The system SHALL print visual banners for run start and iterations.

#### Scenario: Run start banner
- **WHEN** Driver starts
- **THEN** System prints "Starting Littlefactory - Max iterations: N" and "Tasks: M ready"

#### Scenario: Iteration banner
- **WHEN** Each iteration starts
- **THEN** System prints banner with iteration number, max iterations, and target task ID/title

#### Scenario: Run summary
- **WHEN** Run completes
- **THEN** System prints summary with status, run ID, iteration counts, timings, and metadata path
