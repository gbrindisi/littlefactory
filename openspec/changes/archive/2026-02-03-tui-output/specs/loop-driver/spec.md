## MODIFIED Requirements

### Requirement: Visual progress output
The system SHALL emit events to TUI instead of printing directly to stdout.

#### Scenario: Run start event
- **WHEN** Driver starts
- **THEN** System emits RunStartedMsg with max iterations and ready task count

#### Scenario: Iteration start event
- **WHEN** Each iteration starts
- **THEN** System emits IterationStartedMsg with iteration number, task ID, and task title

#### Scenario: Iteration complete event
- **WHEN** Each iteration completes
- **THEN** System emits IterationCompleteMsg with status (completed/failed/timeout)

#### Scenario: Run complete event
- **WHEN** Run completes
- **THEN** System emits RunCompleteMsg with final status and summary metadata

#### Scenario: Tasks refresh event
- **WHEN** Iteration completes
- **THEN** System emits TasksRefreshedMsg with updated task list from TaskSource.List()

## ADDED Requirements

### Requirement: Driver runs as background worker
The system SHALL run Driver in a goroutine, communicating with TUI via event channel.

#### Scenario: Event channel communication
- **WHEN** Driver is started
- **THEN** Driver sends events to provided channel, TUI receives and processes them

#### Scenario: Output writer for agent
- **WHEN** Driver runs an iteration
- **THEN** Driver passes TUI-provided io.Writer to Agent.Run() for output streaming

