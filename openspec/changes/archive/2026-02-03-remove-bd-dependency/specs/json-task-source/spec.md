## ADDED Requirements

### Requirement: JSON file task storage
The system SHALL store tasks in a JSON file at `.littlefactory/tasks.json`.

#### Scenario: Tasks file location
- **WHEN** JSONTaskSource is initialized
- **THEN** System reads tasks from `<project-root>/.littlefactory/tasks.json`

#### Scenario: Tasks file structure
- **WHEN** Tasks file is read
- **THEN** System parses JSON with `tasks` array containing objects with `id`, `title`, `description`, and `status` fields

#### Scenario: Status values
- **WHEN** Task status is read
- **THEN** Status is one of: `todo`, `in_progress`, `done`

### Requirement: JSONTaskSource implements TaskSource interface
The system SHALL provide a JSONTaskSource that implements the TaskSource interface.

#### Scenario: Ready tasks retrieval
- **WHEN** JSONTaskSource.Ready() is called
- **THEN** System returns first task with status `todo` (array order determines priority)

#### Scenario: No ready tasks
- **WHEN** JSONTaskSource.Ready() is called and no tasks have status `todo`
- **THEN** System returns empty list

#### Scenario: List all tasks
- **WHEN** JSONTaskSource.List() is called
- **THEN** System returns all tasks from the JSON file

#### Scenario: Show task details
- **WHEN** JSONTaskSource.Show(id) is called
- **THEN** System returns task with matching ID including all fields

#### Scenario: Show unknown task
- **WHEN** JSONTaskSource.Show(id) is called with non-existent ID
- **THEN** System returns error

### Requirement: Task status updates persist to file
The system SHALL write status changes immediately to the JSON file.

#### Scenario: Claim task
- **WHEN** JSONTaskSource.Claim(id) is called
- **THEN** System sets task status to `in_progress` and writes file

#### Scenario: Close task success
- **WHEN** JSONTaskSource.Close(id, reason) is called
- **THEN** System sets task status to `done` and writes file

#### Scenario: Reset task
- **WHEN** JSONTaskSource.Reset(id) is called
- **THEN** System sets task status to `todo` and writes file

### Requirement: Directory creation on write
The system SHALL create `.littlefactory/` directory if it does not exist when writing tasks file.

#### Scenario: Directory does not exist
- **WHEN** JSONTaskSource writes to tasks file and `.littlefactory/` does not exist
- **THEN** System creates directory before writing file
