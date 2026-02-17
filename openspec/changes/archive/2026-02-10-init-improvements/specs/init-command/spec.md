## MODIFIED Requirements

### Requirement: Init command creates Factoryfile
The system SHALL provide an `init` command that creates a Factoryfile with default configuration, sets up AGENTS.md, updates .gitignore, and installs skills.

#### Scenario: Successful init in empty directory
- **WHEN** User runs `littlefactory init` in a directory without Factoryfile
- **THEN** System creates Factoryfile, AGENTS.md, updates .gitignore, installs skills, and exits with code 0

#### Scenario: Init fails if Factoryfile exists
- **WHEN** User runs `littlefactory init` in a directory with existing Factoryfile
- **THEN** System prints error message and exits with non-zero code without modifying existing file

#### Scenario: Init fails if Factoryfile.yaml exists
- **WHEN** User runs `littlefactory init` in a directory with existing Factoryfile.yaml
- **THEN** System prints error message and exits with non-zero code without modifying existing file

## ADDED Requirements

### Requirement: Init logs all operations
The system SHALL log each step of the init process to stdout with clear progress indicators.

#### Scenario: Init output shows numbered steps
- **WHEN** init runs
- **THEN** system prints numbered steps (e.g., "[1/4] Creating Factoryfile") with indented sub-operations

#### Scenario: Init reports completion
- **WHEN** init completes successfully
- **THEN** system prints summary message indicating littlefactory is ready

### Requirement: Init orchestrates all setup steps
The system SHALL execute AGENTS.md setup, gitignore updates, and skill installation as part of init.

#### Scenario: Init runs all setup steps
- **WHEN** init runs successfully
- **THEN** system executes Factoryfile creation, AGENTS.md setup, gitignore management, and skill installation in order
