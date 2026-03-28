# init-command

## What It Does
The init command bootstraps a littlefactory project by creating a Factoryfile with default configuration, setting up AGENTS.md, updating .gitignore, installing skills (including /lf:* skills), and creating the `.littlefactory/changes/` directory.

## Requirements

### Requirement: Init command creates Factoryfile
The system SHALL provide an `init` command that creates a Factoryfile with default configuration, sets up AGENTS.md, updates .gitignore, installs skills (including `/lf:*` skills), and creates the `.littlefactory/changes/` directory.

#### Scenario: Successful init in empty directory
- **WHEN** user runs `littlefactory init` in a directory without Factoryfile
- **THEN** system creates Factoryfile, AGENTS.md, updates .gitignore, installs skills (including lf-explore, lf-formalize, lf-do, lf-verify), creates `.littlefactory/changes/` directory, and exits with code 0

#### Scenario: Init fails if Factoryfile exists
- **WHEN** user runs `littlefactory init` in a directory with existing Factoryfile
- **THEN** system prints error message and exits with non-zero code without modifying existing file

#### Scenario: Init fails if Factoryfile.yaml exists
- **WHEN** user runs `littlefactory init` in a directory with existing Factoryfile.yaml
- **THEN** system prints error message and exits with non-zero code without modifying existing file

### Requirement: Default Factoryfile content
The system SHALL generate a Factoryfile with sensible defaults for immediate use.

#### Scenario: Default configuration values
- **WHEN** Factoryfile is created by init command
- **THEN** File contains max_iterations: 10, timeout: 600, default_agent: "claude", and agents map with claude agent configured

#### Scenario: Default claude agent configuration
- **WHEN** Factoryfile is created by init command
- **THEN** agents.claude.command is set to "claude --dangerously-skip-permissions --print"

### Requirement: Init logs all operations
The system SHALL log each step of the init process to stdout with clear progress indicators.

#### Scenario: Init output shows numbered steps
- **WHEN** init runs
- **THEN** system prints numbered steps [1/5] through [5/5] with indented sub-operations (Factoryfile, AGENTS.md, .gitignore, skills, changes directory)

#### Scenario: Init reports completion
- **WHEN** init completes successfully
- **THEN** system prints summary message indicating littlefactory is ready

### Requirement: Init orchestrates all setup steps
The system SHALL execute Factoryfile creation, AGENTS.md setup, gitignore updates, skill installation, and changes directory creation as part of init.

#### Scenario: Init runs all setup steps
- **WHEN** init runs successfully
- **THEN** system creates Factoryfile, sets up AGENTS.md, updates gitignore, installs skills, and creates `.littlefactory/changes/` directory in order

## Boundaries

## Gotchas
