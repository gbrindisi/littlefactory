# specs-dir-config

## What It Does
Allows projects to configure a `specs_dir` field in the Factoryfile, defining where long-lived specification files are stored. This directory is used by skills like `/lf-archive` to merge delta specs from changes into the project's canonical spec library.

## Requirements

### Requirement: Factoryfile accepts specs_dir field
The system SHALL support a `specs_dir` field in the Factoryfile that configures the directory for long-lived specification files.

#### Scenario: specs_dir is configured
- **WHEN** the Factoryfile contains `specs_dir: "specs/"`
- **THEN** the system uses `specs/` as the specs directory for reading and writing spec files

#### Scenario: specs_dir is not configured
- **WHEN** the Factoryfile does not contain a `specs_dir` field
- **THEN** the system does not assume a default specs directory and features that depend on it (e.g., `/lf-archive` spec merging) ask the user or skip

#### Scenario: specs_dir resolved relative to project root
- **WHEN** `specs_dir` is set to a relative path
- **THEN** the system resolves it relative to the project root (the directory containing the Factoryfile)

## Boundaries

## Gotchas
