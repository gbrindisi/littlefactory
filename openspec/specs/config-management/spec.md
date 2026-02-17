## MODIFIED Requirements

### Requirement: Factoryfile format
The system SHALL support YAML-based Factoryfile configuration with named agents.

#### Scenario: Basic configuration keys
- **WHEN** Factoryfile contains max_iterations, timeout, default_agent, and agents keys
- **THEN** System parses and applies values correctly

#### Scenario: Missing Factoryfile
- **WHEN** Factoryfile does not exist
- **THEN** System returns error (Factoryfile now required for run command)

#### Scenario: Invalid Factoryfile
- **WHEN** Factoryfile has invalid YAML syntax
- **THEN** System returns error and does not start

#### Scenario: Agents map structure
- **WHEN** Factoryfile contains agents map with named entries
- **THEN** System parses each agent with its command field

#### Scenario: Default agent validation
- **WHEN** Factoryfile specifies default_agent
- **THEN** System validates that the name exists in agents map

#### Scenario: Missing default_agent
- **WHEN** Factoryfile does not specify default_agent
- **THEN** System returns error during config validation

#### Scenario: Empty agents map
- **WHEN** Factoryfile has empty agents map
- **THEN** System returns error during config validation

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

### Requirement: Worktrees directory configuration
The system SHALL support a `worktrees_dir` option in Factoryfile.

#### Scenario: Parse worktrees_dir
- **WHEN** Factoryfile contains `worktrees_dir: ../worktrees`
- **THEN** System parses and stores the worktrees directory path

#### Scenario: Missing worktrees_dir uses default
- **WHEN** Factoryfile does not specify `worktrees_dir`
- **THEN** System uses `..` (sibling to repo) as default

## REMOVED Requirements

### Requirement: Future agent configuration support
**Reason**: Replaced by full agents map implementation - no longer "future", now implemented.
**Migration**: Change `agent:` section to `agents:` map with named entries.
