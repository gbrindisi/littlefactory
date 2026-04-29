# lf-do-skill

## What It Does
The `/lf-do` skill runs a littlefactory change as a background process, monitors its progress, and reports the outcome. It infers the change name from the project state or prompts the user when ambiguous.

## Requirements

### Requirement: Do skill runs littlefactory in background
The system SHALL provide an embedded `/lf-do` skill that invokes `littlefactory run` as a background process and monitors it to completion.

#### Scenario: Invokes run command
- **WHEN** the user invokes `/lf-do` with a determined change name
- **THEN** the skill runs `littlefactory run -c <name>` in the background

#### Scenario: Monitors via status command
- **WHEN** the background process is running
- **THEN** the skill checks progress via `littlefactory status -c <name>` when the background process completes or when the user asks for an update

#### Scenario: Reports completion
- **WHEN** the background process finishes
- **THEN** the skill reports the outcome (all tasks completed, some failed, or cancelled) and suggests `/lf-verify` on success

### Requirement: Do skill infers change name
The system SHALL have the do skill infer the change name from the project state or prompt for selection if ambiguous.

#### Scenario: Single active change
- **WHEN** the user invokes `/lf-do` and exactly one change directory exists under `.littlefactory/changes/`
- **THEN** the skill uses that change automatically

#### Scenario: Multiple active changes
- **WHEN** the user invokes `/lf-do` and multiple change directories exist
- **THEN** the skill lists the available changes and prompts the user to select one

## Boundaries

## Gotchas
