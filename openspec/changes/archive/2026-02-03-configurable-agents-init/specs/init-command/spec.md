## ADDED Requirements

### Requirement: Init command creates Factoryfile
The system SHALL provide an `init` command that creates a Factoryfile with default configuration.

#### Scenario: Successful init in empty directory
- **WHEN** User runs `littlefactory init` in a directory without Factoryfile
- **THEN** System creates Factoryfile with default configuration and exits with code 0

#### Scenario: Init fails if Factoryfile exists
- **WHEN** User runs `littlefactory init` in a directory with existing Factoryfile
- **THEN** System prints error message and exits with non-zero code without modifying existing file

#### Scenario: Init fails if Factoryfile.yaml exists
- **WHEN** User runs `littlefactory init` in a directory with existing Factoryfile.yaml
- **THEN** System prints error message and exits with non-zero code without modifying existing file

### Requirement: Default Factoryfile content
The system SHALL generate a Factoryfile with sensible defaults for immediate use.

#### Scenario: Default configuration values
- **WHEN** Factoryfile is created by init command
- **THEN** File contains max_iterations: 10, timeout: 600, default_agent: "claude", and agents map with claude agent configured

#### Scenario: Default claude agent configuration
- **WHEN** Factoryfile is created by init command
- **THEN** agents.claude.command is set to "claude --dangerously-skip-permissions --print"
