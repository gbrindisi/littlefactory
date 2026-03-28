# change-management

## What It Does
Change directories live under `.littlefactory/changes/` and represent units of work with a fixed artifact structure. Each change contains a tasks.json and optionally proposal, specs, and design documents. The changes directory is created during init and preserved across upgrades.

## Requirements

### Requirement: Changes directory exists under state dir
The system SHALL store changes under `.littlefactory/changes/<name>/` where each change is a directory containing artifact files.

#### Scenario: Change directory created by formalize skill
- **WHEN** the `/lf:formalize` skill creates a new change
- **THEN** a directory `.littlefactory/changes/<name>/` is created containing proposal.md, specs/\*/spec.md, design.md (optional), and tasks.json

#### Scenario: Change directory contains tasks.json
- **WHEN** a change exists at `.littlefactory/changes/<name>/`
- **THEN** the directory MUST contain a valid tasks.json file consumable by `littlefactory run -c <name>`

### Requirement: Init creates changes directory
The system SHALL create the `.littlefactory/changes/` directory during init if it does not exist.

#### Scenario: Changes directory created on init
- **WHEN** `littlefactory init` runs
- **THEN** `.littlefactory/changes/` directory exists after completion

#### Scenario: Changes directory preserved on upgrade
- **WHEN** `littlefactory upgrade` runs and `.littlefactory/changes/` already exists with content
- **THEN** existing changes are not modified or deleted

### Requirement: Change artifacts follow fixed structure
The system SHALL expect changes to contain artifacts in a fixed sequence: proposal.md, specs/\*/spec.md, design.md (optional), tasks.json.

#### Scenario: Minimal valid change
- **WHEN** a change directory contains only tasks.json
- **THEN** the change is valid for `littlefactory run -c <name>`

#### Scenario: Full change with all artifacts
- **WHEN** a change is created by `/lf:formalize`
- **THEN** the directory contains proposal.md, at least one specs/\*/spec.md, optionally design.md, and tasks.json

## Boundaries

## Gotchas
