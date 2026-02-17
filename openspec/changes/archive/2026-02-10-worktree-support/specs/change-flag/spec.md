## ADDED Requirements

### Requirement: Change flag on run command
The system SHALL accept `--change` / `-c` flag on the `run` command to specify which openspec change to use as the task source.

#### Scenario: Run with change flag
- **WHEN** User runs `littlefactory run claude --change feature-a`
- **THEN** System uses `openspec/changes/feature-a/tasks.json` as the task source

#### Scenario: Run with short change flag
- **WHEN** User runs `littlefactory run claude -c feature-a`
- **THEN** System uses `openspec/changes/feature-a/tasks.json` as the task source

#### Scenario: Change not found
- **WHEN** User runs `littlefactory run claude --change nonexistent`
- **THEN** System returns error "Change 'nonexistent' not found at openspec/changes/nonexistent/"

#### Scenario: Change has no tasks.json
- **WHEN** User runs `littlefactory run claude --change incomplete-change` and tasks.json does not exist
- **THEN** System returns error "No tasks.json found for change 'incomplete-change'"

### Requirement: Change flag is optional
The system SHALL allow running without `--change` flag, using the default task source.

#### Scenario: Run without change flag
- **WHEN** User runs `littlefactory run claude` without `--change`
- **THEN** System uses the default task source from `<state_dir>/tasks.json`
