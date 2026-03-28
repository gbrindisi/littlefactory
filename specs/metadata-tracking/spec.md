# metadata-tracking

## What It Does
Tracks run-level and iteration-level metadata during littlefactory execution, including timing, status, output metrics, and aggregate statistics. Metadata is serialized to JSON after each iteration for persistence and backward compatibility.

## Requirements
### Requirement: Run metadata structure
The system SHALL track run-level metadata including run ID, timestamps, status, and aggregate statistics.

#### Scenario: Run initialization
- **WHEN** Driver starts
- **THEN** System creates RunMetadata with generated run_id (format: YYYYMMDD-HHMMSS), started_at timestamp, and RUNNING status

#### Scenario: Run finalization
- **WHEN** Run completes with any status
- **THEN** System updates RunMetadata with ended_at, final status, total_duration_seconds, and avg_iteration_duration_seconds

### Requirement: Iteration metadata structure
The system SHALL track per-iteration metadata including timing, status, task info, and output metrics.

#### Scenario: Iteration start tracking
- **WHEN** Iteration begins
- **THEN** System creates IterationMetadata with iteration_number, started_at, RUNNING status, and target task_id/task_title

#### Scenario: Iteration completion tracking
- **WHEN** Iteration ends
- **THEN** System updates IterationMetadata with ended_at, duration_seconds, final status, exit_code, error_message (if any), output_lines, output_bytes, session_id, and session_path

### Requirement: JSON serialization
The system SHALL serialize metadata to JSON at `<state_dir>/run_metadata.json`.

#### Scenario: Metadata save after each iteration
- **WHEN** Each iteration completes
- **THEN** System writes updated RunMetadata to `<state_dir>/run_metadata.json` with proper ISO8601 timestamps

#### Scenario: Backward compatible JSON format
- **WHEN** System serializes metadata
- **THEN** JSON structure matches Python ciccio format exactly for backward compatibility

### Requirement: Aggregate statistics
The system SHALL compute and track aggregate statistics across iterations.

#### Scenario: Success and failure counts
- **WHEN** Iterations complete
- **THEN** RunMetadata tracks total_iterations, successful_iterations, and failed_iterations

#### Scenario: Average duration calculation
- **WHEN** Run finalizes
- **THEN** System computes avg_iteration_duration_seconds from all completed iterations

### Requirement: SaveMetadata receives config
The system SHALL pass full config to SaveMetadata for state directory access.

#### Scenario: SaveMetadata signature
- **WHEN** SaveMetadata is called
- **THEN** Function receives projectRoot, *config.Config, and metadata parameters

## Boundaries

## Gotchas
