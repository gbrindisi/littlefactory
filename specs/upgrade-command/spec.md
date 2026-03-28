# upgrade-command

## What It Does
The upgrade command updates an existing littlefactory project to the latest version, re-applying AGENTS.md setup, gitignore rules, embedded skills, and the changes directory structure. It is idempotent and logs all operations performed.

## Requirements

### Requirement: Upgrade command for existing projects
The system SHALL provide an `upgrade` command that updates an existing littlefactory project to the latest version. The command MUST require a valid Factoryfile in the project root.

#### Scenario: Factoryfile is present
- **WHEN** the user runs `littlefactory upgrade` in a directory containing a Factoryfile
- **THEN** the system applies AGENTS.md setup, gitignore management, skill installation, and changes directory structure

#### Scenario: Factoryfile is missing
- **WHEN** the user runs `littlefactory upgrade` in a directory without a Factoryfile
- **THEN** the system exits with an error indicating the project has not been initialized

### Requirement: Upgrade is idempotent
The system SHALL produce the same result when the upgrade command is run multiple times in succession.

#### Scenario: Repeated upgrade produces same result
- **WHEN** the user runs `littlefactory upgrade` twice consecutively without changes
- **THEN** the project state after the second run is identical to the state after the first run

### Requirement: Upgrade installs new skills
The system SHALL install any new embedded skills that were added since the last upgrade, while leaving existing matching skills unchanged.

#### Scenario: New skill available in latest version
- **WHEN** the upgrade runs and a new skill exists in the embedded skill set that is not present in the project
- **THEN** the system installs the new skill into `.littlefactory/skills/`

#### Scenario: Existing skill matches embedded version
- **WHEN** the upgrade runs and an installed skill already matches the embedded version
- **THEN** the system leaves the existing skill unchanged

### Requirement: Upgrade logs all operations
The system SHALL log every operation performed during upgrade so the user can see what changed.

#### Scenario: Operations are logged
- **WHEN** the upgrade command completes
- **THEN** the system outputs a log of all operations performed (e.g., files created, skills installed, files updated)

## Boundaries

## Gotchas
