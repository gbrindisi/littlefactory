## MODIFIED Requirements

### Requirement: Factoryfile format
The system SHALL support YAML-based Factoryfile configuration.

#### Scenario: Basic configuration keys
- **WHEN** Factoryfile contains max_iterations and timeout keys
- **THEN** System parses and applies integer values

#### Scenario: Missing Factoryfile
- **WHEN** Factoryfile does not exist
- **THEN** System continues with defaults (no error)

#### Scenario: Invalid Factoryfile
- **WHEN** Factoryfile has invalid YAML syntax
- **THEN** System returns error and does not start

#### Scenario: Agent env configuration
- **WHEN** Factoryfile agent has `env` field with string values
- **THEN** System parses as static environment variables

#### Scenario: Agent env with shell configuration
- **WHEN** Factoryfile agent has `env` field with `{ shell: "command" }` values
- **THEN** System parses as dynamic shell-evaluated environment variables
