## ADDED Requirements

### Requirement: State directory configuration
The system SHALL support a configurable state directory via Factoryfile.

#### Scenario: Default state directory
- **WHEN** Factoryfile does not specify `state_dir`
- **THEN** System uses `.littlefactory` as the state directory

#### Scenario: Custom state directory
- **WHEN** Factoryfile specifies `state_dir: custom-dir`
- **THEN** System uses `custom-dir` as the state directory for all runtime files

#### Scenario: State directory validation
- **WHEN** Factoryfile specifies empty `state_dir`
- **THEN** System returns error during config validation
