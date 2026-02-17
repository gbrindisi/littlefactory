## ADDED Requirements

### Requirement: Agent supports environment variable configuration
The system SHALL allow agents to be configured with environment variables via the `env` field in Factoryfile.

#### Scenario: Static environment variable
- **WHEN** Agent config has `env: { VAR: "value" }`
- **THEN** Agent process receives VAR=value in its environment

#### Scenario: Dynamic environment variable via shell
- **WHEN** Agent config has `env: { VAR: { shell: "echo hello" } }`
- **THEN** System executes shell command and sets VAR to stdout (trimmed)

#### Scenario: Multiple environment variables
- **WHEN** Agent config has multiple env entries (static and dynamic mixed)
- **THEN** All variables are set in agent process environment

### Requirement: Agent inherits parent environment
The system SHALL pass the parent process environment to agents, with `env` config adding or overriding values.

#### Scenario: Parent env inherited
- **WHEN** Agent runs without `env` config
- **THEN** Agent receives all environment variables from parent process

#### Scenario: Env config overrides parent
- **WHEN** Parent has VAR=old and agent config has `env: { VAR: "new" }`
- **THEN** Agent receives VAR=new

### Requirement: Shell command failures are fatal
The system SHALL fail agent startup if any shell command in `env` config fails.

#### Scenario: Shell command returns non-zero
- **WHEN** Shell command in env config exits with non-zero status
- **THEN** Agent fails to start with error indicating which variable failed

#### Scenario: Shell command not found
- **WHEN** Shell command references non-existent command
- **THEN** Agent fails to start with error
