## MODIFIED Requirements

### Requirement: Init command creates Factoryfile
The system SHALL provide an `init` command that checks for the `openspec` prerequisite, creates a Factoryfile with default configuration, sets up AGENTS.md, updates .gitignore, installs skills, and configures OpenSpec.

#### Scenario: Successful init in empty directory
- **WHEN** user runs `littlefactory init` in a directory without Factoryfile and with `openspec` in PATH
- **THEN** system creates Factoryfile, AGENTS.md, updates .gitignore, installs skills, sets up OpenSpec schema and config, and exits with code 0

#### Scenario: Init fails if Factoryfile exists
- **WHEN** user runs `littlefactory init` in a directory with existing Factoryfile
- **THEN** system prints error message and exits with non-zero code without modifying existing file

#### Scenario: Init fails if Factoryfile.yaml exists
- **WHEN** user runs `littlefactory init` in a directory with existing Factoryfile.yaml
- **THEN** system prints error message and exits with non-zero code without modifying existing file

#### Scenario: Init fails if openspec not installed
- **WHEN** user runs `littlefactory init` and `openspec` is not found in PATH
- **THEN** system prints error message and exits with non-zero code without creating any files

### Requirement: Init logs all operations
The system SHALL log each step of the init process to stdout with clear progress indicators.

#### Scenario: Init output shows numbered steps
- **WHEN** init runs
- **THEN** system prints numbered steps [1/5] through [5/5] with indented sub-operations, including the OpenSpec setup step

#### Scenario: Init reports completion
- **WHEN** init completes successfully
- **THEN** system prints summary message indicating littlefactory is ready

### Requirement: Init orchestrates all setup steps
The system SHALL execute the openspec prerequisite check, Factoryfile creation, AGENTS.md setup, gitignore updates, skill installation, and OpenSpec setup as part of init.

#### Scenario: Init runs all setup steps
- **WHEN** init runs successfully
- **THEN** system checks openspec is installed, creates Factoryfile, sets up AGENTS.md, updates gitignore, installs skills, and sets up OpenSpec schema and config in order
