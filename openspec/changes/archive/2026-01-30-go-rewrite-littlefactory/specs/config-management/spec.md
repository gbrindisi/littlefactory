## ADDED Requirements

### Requirement: Configuration loading hierarchy
The system SHALL load configuration from defaults, Factoryfile, then CLI flags in order of precedence.

#### Scenario: Default values
- **WHEN** No Factoryfile or flags provided
- **THEN** System uses hardcoded defaults (max_iterations: 10, timeout: 600)

#### Scenario: Factoryfile overrides defaults
- **WHEN** Factoryfile exists at project root
- **THEN** System loads values from Factoryfile, overriding defaults

#### Scenario: CLI flags override everything
- **WHEN** CLI flags are provided
- **THEN** System uses flag values, overriding both defaults and Factoryfile

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

### Requirement: Future agent configuration support
The system SHALL define Factoryfile structure allowing future agent type configuration.

#### Scenario: Reserved agent section
- **WHEN** Factoryfile includes agent section (future use)
- **THEN** System parses but does not require it (reserved for future)
