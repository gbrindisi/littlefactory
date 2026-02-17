## MODIFIED Requirements

### Requirement: JSON file task storage
The system SHALL store tasks in a JSON file at the configured or specified location.

#### Scenario: Tasks file location with explicit path
- **WHEN** JSONTaskSource is initialized with explicit path via NewJSONTaskSourceWithPath
- **THEN** System reads tasks from the specified path

#### Scenario: Tasks file location with config
- **WHEN** JSONTaskSource is initialized via NewJSONTaskSource
- **THEN** System reads tasks from `<project-root>/<state_dir>/tasks.json`

#### Scenario: Tasks file structure
- **WHEN** Tasks file is read
- **THEN** System parses JSON with `tasks` array containing objects with `id`, `title`, `description`, `status`, `labels`, and `blockers` fields

#### Scenario: Status values
- **WHEN** Task status is read
- **THEN** Status is one of: `todo`, `in_progress`, `done`

#### Scenario: Validation on load
- **WHEN** JSONTaskSource is initialized
- **THEN** System validates tasks.json and returns error if validation fails

#### Scenario: Explicit path file not found
- **WHEN** NewJSONTaskSourceWithPath is called with non-existent path
- **THEN** System returns error indicating file not found
