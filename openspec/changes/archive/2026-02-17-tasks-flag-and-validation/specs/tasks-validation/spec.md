## ADDED Requirements

### Requirement: Tasks.json validation on load
The system SHALL validate tasks.json structure and content when loading from any source.

#### Scenario: Valid JSON required
- **WHEN** tasks.json contains invalid JSON syntax
- **THEN** System returns error with file path and parse error details

#### Scenario: Tasks array required
- **WHEN** tasks.json is valid JSON but missing "tasks" array
- **THEN** System returns error "missing required field 'tasks'"

#### Scenario: Empty tasks array allowed
- **WHEN** tasks.json contains empty tasks array
- **THEN** System loads successfully with no tasks

### Requirement: Required task fields
The system SHALL validate that each task has required fields with non-empty values.

#### Scenario: Missing id field
- **WHEN** Task has empty or missing "id" field
- **THEN** System returns error 'task at index N: missing required field "id"'

#### Scenario: Missing title field
- **WHEN** Task has empty or missing "title" field
- **THEN** System returns error 'task "ID" (index N): missing required field "title"'

#### Scenario: Missing status field
- **WHEN** Task has empty or missing "status" field
- **THEN** System returns error 'task "ID" (index N): missing required field "status"'

#### Scenario: Empty description allowed
- **WHEN** Task has empty description field
- **THEN** System loads successfully (description is optional)

### Requirement: Valid status values
The system SHALL validate that task status is one of the allowed values.

#### Scenario: Valid status values
- **WHEN** Task has status "todo", "in_progress", or "done"
- **THEN** System loads successfully

#### Scenario: Invalid status value
- **WHEN** Task has status value not in allowed list (e.g., "pending", "complete")
- **THEN** System returns error 'task "ID" (index N): invalid status "VALUE" (must be: todo, in_progress, done)'

### Requirement: Unique task IDs
The system SHALL validate that all task IDs are unique within the file.

#### Scenario: Duplicate task ID
- **WHEN** Two or more tasks have the same ID
- **THEN** System returns error 'task "ID" (index N): duplicate id (first seen at index M)'

### Requirement: Sequential blocker chain
The system SHALL validate that tasks form a strict linear sequence via blockers.

#### Scenario: Exactly one root task
- **WHEN** Tasks are validated
- **THEN** Exactly one task has empty blockers array

#### Scenario: Multiple root tasks
- **WHEN** More than one task has empty blockers array
- **THEN** System returns error 'invalid sequence: multiple root tasks (tasks with no blockers): "ID1", "ID2"'

#### Scenario: No root task
- **WHEN** All tasks have non-empty blockers
- **THEN** System returns error 'invalid sequence: no root task (all tasks have blockers)'

#### Scenario: Single blocker per task
- **WHEN** Non-root task is validated
- **THEN** Task has exactly one blocker

#### Scenario: Multiple blockers on task
- **WHEN** Task has more than one blocker
- **THEN** System returns error 'invalid sequence: task "ID" has N blockers, expected 0 or 1'

#### Scenario: Blocker references existing task
- **WHEN** Task has blocker that references non-existent task ID
- **THEN** System returns error 'task "ID" (index N): blocker "BLOCKER_ID" does not exist'

#### Scenario: Complete chain coverage
- **WHEN** Tasks are validated
- **THEN** All tasks are reachable by walking from root through blocker references

#### Scenario: Orphaned tasks
- **WHEN** Task is not reachable from root via blocker chain
- **THEN** System returns error 'invalid sequence: tasks not reachable from root: "ID1", "ID2"'

### Requirement: Multi-error reporting
The system SHALL report all validation errors at once rather than stopping at first failure.

#### Scenario: Multiple validation errors
- **WHEN** tasks.json has multiple validation issues
- **THEN** System returns error containing all issues, one per line

#### Scenario: Error message header
- **WHEN** Validation fails
- **THEN** Error message starts with 'Error loading tasks from PATH:'
