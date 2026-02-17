## ADDED Requirements

### Requirement: OpenSpec binary prerequisite check
The system SHALL verify that the `openspec` binary is available in PATH before proceeding with init or upgrade.

#### Scenario: openspec is installed
- **WHEN** user runs `littlefactory init` and `openspec` is found in PATH
- **THEN** system proceeds with the init workflow without error

#### Scenario: openspec is not installed
- **WHEN** user runs `littlefactory init` and `openspec` is not found in PATH
- **THEN** system prints an error message indicating openspec is required and exits with non-zero code without creating any files

#### Scenario: openspec check runs before any other step
- **WHEN** user runs `littlefactory init` without openspec installed
- **THEN** system does not create a Factoryfile, AGENTS.md, .gitignore entries, or any other files

### Requirement: Schema installation
The system SHALL extract the embedded littlefactory schema to `openspec/schemas/littlefactory/` in the project directory.

#### Scenario: Schema extracted to project
- **WHEN** init runs the OpenSpec setup step
- **THEN** system creates `openspec/schemas/littlefactory/schema.yaml` and all template files under `openspec/schemas/littlefactory/templates/`

#### Scenario: Schema directory created if missing
- **WHEN** init runs and `openspec/schemas/` does not exist
- **THEN** system creates the full directory path `openspec/schemas/littlefactory/` and extracts schema files into it

#### Scenario: Schema overwrites existing files
- **WHEN** init runs and `openspec/schemas/littlefactory/` already contains files
- **THEN** system overwrites existing schema files with the embedded versions

### Requirement: OpenSpec config setup
The system SHALL create `openspec/config.yaml` with the default schema set to `littlefactory` if the file does not exist.

#### Scenario: Config created when missing
- **WHEN** init runs and `openspec/config.yaml` does not exist
- **THEN** system creates `openspec/config.yaml` with `schema: littlefactory`

#### Scenario: Config preserved when existing
- **WHEN** init runs and `openspec/config.yaml` already exists
- **THEN** system does not modify the existing config file to preserve user customizations

### Requirement: OpenSpec setup is idempotent
The system SHALL allow the OpenSpec setup step to run multiple times without error.

#### Scenario: Running upgrade after init
- **WHEN** user runs `littlefactory upgrade` after a previous successful init
- **THEN** system completes the OpenSpec setup step without error, updating schema files but preserving existing config

#### Scenario: Running upgrade twice
- **WHEN** user runs `littlefactory upgrade` on a project that was already upgraded
- **THEN** system logs appropriate status for schema (updated/current) and config (preserved) and exits successfully
