## ADDED Requirements

### Requirement: Upgrade command for existing projects
The system SHALL provide an `upgrade` command that applies init improvements to existing projects.

#### Scenario: Upgrade requires existing Factoryfile
- **WHEN** user runs `littlefactory upgrade` in a directory without Factoryfile
- **THEN** system prints error "No Factoryfile found. Run 'littlefactory init' first." and exits with non-zero code

#### Scenario: Successful upgrade
- **WHEN** user runs `littlefactory upgrade` in a directory with Factoryfile
- **THEN** system applies AGENTS.md setup, gitignore updates, and skill installation

### Requirement: Upgrade is idempotent
The system SHALL allow upgrade to run multiple times without adverse effects.

#### Scenario: Upgrade on already-upgraded project
- **WHEN** user runs `littlefactory upgrade` on a project that was already upgraded
- **THEN** system logs "already configured" or "already present" for each step and makes no changes

### Requirement: Upgrade installs new skills
The system SHALL install skills that were added in newer versions of littlefactory.

#### Scenario: New skill available in upgraded binary
- **WHEN** upgrade runs and binary contains skill not present in .littlefactory/skills/
- **THEN** system extracts new skill and creates symlink if .claude/ exists

#### Scenario: Existing skill matches embedded version
- **WHEN** upgrade runs and skill already exists in .littlefactory/skills/
- **THEN** system logs "up to date" and does not overwrite existing skill

### Requirement: Upgrade logs all operations
The system SHALL log each step of the upgrade process to stdout.

#### Scenario: Upgrade output shows progress
- **WHEN** upgrade runs
- **THEN** system prints numbered steps with indented sub-operations for each change made or skipped
