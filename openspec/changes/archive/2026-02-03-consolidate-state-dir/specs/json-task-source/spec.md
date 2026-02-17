## MODIFIED Requirements

### Requirement: JSON file task storage
The system SHALL store tasks in a JSON file at `<state_dir>/tasks.json`.

#### Scenario: Tasks file location
- **WHEN** JSONTaskSource is initialized
- **THEN** System reads tasks from `<project-root>/<state_dir>/tasks.json`

#### Scenario: Tasks file structure
- **WHEN** Tasks file is read
- **THEN** System parses JSON with `tasks` array containing objects with `id`, `title`, `description`, and `status` fields

#### Scenario: Status values
- **WHEN** Task status is read
- **THEN** Status is one of: `todo`, `in_progress`, `done`

### Requirement: JSONTaskSource receives config at construction
The system SHALL pass full config to JSONTaskSource constructor for state directory access.

#### Scenario: NewJSONTaskSource signature
- **WHEN** NewJSONTaskSource is called
- **THEN** Function receives projectRoot and *config.Config parameters

#### Scenario: State directory from config
- **WHEN** JSONTaskSource constructs file path
- **THEN** System uses cfg.StateDir instead of hardcoded ".littlefactory"

### Requirement: Directory creation on write
The system SHALL create state directory if it does not exist when writing tasks file.

#### Scenario: Directory does not exist
- **WHEN** JSONTaskSource writes to tasks file and state directory does not exist
- **THEN** System creates directory before writing file
