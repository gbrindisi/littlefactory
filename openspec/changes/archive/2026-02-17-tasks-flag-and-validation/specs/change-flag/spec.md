## ADDED Requirements

### Requirement: Tasks flag on run command
The system SHALL accept `--tasks` / `-t` flag on the `run` command to specify an explicit path to tasks.json.

#### Scenario: Run with tasks flag
- **WHEN** User runs `littlefactory run claude --tasks path/to/tasks.json`
- **THEN** System uses the specified file as the task source

#### Scenario: Run with short tasks flag
- **WHEN** User runs `littlefactory run claude -t path/to/tasks.json`
- **THEN** System uses the specified file as the task source

#### Scenario: Tasks file not found
- **WHEN** User runs `littlefactory run claude --tasks nonexistent.json`
- **THEN** System returns error "tasks file not found: nonexistent.json"

#### Scenario: Tasks flag accepts relative path
- **WHEN** User runs `littlefactory run claude -t ./my-tasks.json`
- **THEN** System resolves path relative to current working directory

#### Scenario: Tasks flag accepts absolute path
- **WHEN** User runs `littlefactory run claude -t /absolute/path/tasks.json`
- **THEN** System uses the absolute path directly

### Requirement: Flag priority resolution
The system SHALL resolve task source with priority: --tasks > --change > default.

#### Scenario: Tasks flag overrides change flag
- **WHEN** User runs `littlefactory run claude --tasks custom.json --change feature-a`
- **THEN** System uses custom.json (ignores --change)

#### Scenario: Tasks flag overrides default
- **WHEN** User runs `littlefactory run claude --tasks custom.json`
- **THEN** System uses custom.json (does not use default state_dir/tasks.json)

#### Scenario: Change flag used when no tasks flag
- **WHEN** User runs `littlefactory run claude --change feature-a` without --tasks
- **THEN** System uses openspec/changes/feature-a/tasks.json

#### Scenario: Default used when no flags
- **WHEN** User runs `littlefactory run claude` without --tasks or --change
- **THEN** System uses state_dir/tasks.json
