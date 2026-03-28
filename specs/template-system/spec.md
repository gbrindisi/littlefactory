# template-system

## What It Does
Embeds the CLAUDE.md agent template in the littlefactory binary at compile time. The template instructs agents on how to work within the littlefactory loop: read progress, implement the task, run checks, commit, and update progress -- without needing to manage task status themselves.

## Requirements
### Requirement: Embedded template content
The system SHALL embed CLAUDE.md template in binary at compile time.

#### Scenario: Template available offline
- **WHEN** Binary runs in environment without access to source
- **THEN** Embedded template is available for use

#### Scenario: Template excludes task management commands
- **WHEN** Template is rendered
- **THEN** Template does not include bd commands or task status update instructions

#### Scenario: Template workflow focuses on implementation
- **WHEN** Template is rendered
- **THEN** Template instructs agent to: read progress, implement task, run checks, commit, update progress

#### Scenario: Template indicates automatic completion
- **WHEN** Template is rendered
- **THEN** Template informs agent that task completion is handled automatically by littlefactory

## Boundaries

## Gotchas
