## ADDED Requirements

### Requirement: Template loading with local override
The system SHALL load templates from local override first, falling back to embedded template.

#### Scenario: Local override exists
- **WHEN** File exists at tasks/CLAUDE.md in project root
- **THEN** System loads template content from local file

#### Scenario: No local override
- **WHEN** File does not exist at tasks/CLAUDE.md
- **THEN** System loads embedded template from binary (via go:embed)

### Requirement: Task detail injection
The system SHALL inject task details into template using string replacement.

#### Scenario: Task details injection
- **WHEN** System renders template with task
- **THEN** System replaces {task_id}, {task_title}, and {task_description} placeholders with task values

#### Scenario: No task available
- **WHEN** System renders template without task (no ready tasks)
- **THEN** System returns template unchanged with placeholders intact

### Requirement: Embedded template content
The system SHALL embed CLAUDE.md template in binary at compile time.

#### Scenario: Template available offline
- **WHEN** Binary runs in environment without access to source
- **THEN** Embedded template is available for use

### Requirement: Template preserves ciccio format
The system SHALL use template content matching ciccio's original CLAUDE.md format.

#### Scenario: Agent instructions compatibility
- **WHEN** Template is rendered
- **THEN** Output matches ciccio template structure with task injection points, workflow steps, and quality requirements
