## ADDED Requirements

### Requirement: Add littlefactory runtime files to gitignore
The system SHALL add littlefactory runtime files to .gitignore during init.

#### Scenario: Add run_metadata.json to gitignore
- **WHEN** init runs and .gitignore does not contain .littlefactory/run_metadata.json
- **THEN** system appends .littlefactory/run_metadata.json to .gitignore

#### Scenario: Add tasks.json to gitignore
- **WHEN** init runs and .gitignore does not contain .littlefactory/tasks.json
- **THEN** system appends .littlefactory/tasks.json to .gitignore

### Requirement: Gitignore updates are idempotent
The system SHALL not duplicate entries in .gitignore.

#### Scenario: Entry already exists in gitignore
- **WHEN** init or upgrade runs and .gitignore already contains the entry
- **THEN** system logs "already present" and does not add duplicate

### Requirement: Create gitignore if missing
The system SHALL create .gitignore if it does not exist.

#### Scenario: No gitignore exists
- **WHEN** init runs in a directory without .gitignore
- **THEN** system creates .gitignore with littlefactory runtime entries

### Requirement: Preserve existing gitignore content
The system SHALL append to existing .gitignore without modifying existing entries.

#### Scenario: Existing gitignore has custom entries
- **WHEN** init runs and .gitignore contains user's custom entries
- **THEN** system appends new entries at end without removing or modifying existing content
