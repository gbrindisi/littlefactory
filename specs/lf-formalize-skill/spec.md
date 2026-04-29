# lf-formalize-skill

## What It Does
The `/lf-formalize` skill turns a conversation into a structured change by deriving a change name from context and sequentially generating all artifacts (proposal.md, delta specs, optional design.md, tasks.json) in a single invocation. Each artifact builds on the previous ones, and tasks carry fat context so agents can execute them independently.

## Requirements

### Requirement: Formalize derives change from conversation
The system SHALL provide an embedded `/lf-formalize` skill that derives the change name and content entirely from conversation context, requiring no arguments.

#### Scenario: No arguments required
- **WHEN** the user invokes `/lf-formalize` without arguments and the conversation contains sufficient context
- **THEN** the skill derives a kebab-case change name and generates all artifacts from the conversation

#### Scenario: Formalize creates change directory
- **WHEN** the skill derives a change name
- **THEN** it creates `.littlefactory/changes/<name>/` with subdirectories for specs

### Requirement: Formalize generates artifacts sequentially
The system SHALL generate artifacts in a dependency chain where each artifact reads the previous ones before being created.

#### Scenario: Artifact dependency chain
- **WHEN** the skill generates artifacts
- **THEN** it creates them in order: proposal.md, then specs (one per capability from the proposal), then design.md (conditional), then tasks.json, reading each prior artifact before generating the next

#### Scenario: Design skipped for trivial changes
- **WHEN** the change does not involve cross-cutting concerns, new dependencies, security/performance complexity, or architectural ambiguity
- **THEN** the skill skips design.md generation and proceeds directly to tasks.json

### Requirement: Formalize generates tasks.json with fat context
The system SHALL have the formalize skill generate tasks.json where each task has a self-contained description including context, scope, checklist, implementation plan, acceptance criteria, and key references. The formalize skill SHALL NOT generate a housekeeping spec merge task -- spec merging is handled interactively by `/lf-archive`.

#### Scenario: Fat context task description
- **WHEN** tasks.json is generated
- **THEN** each task description contains ## Context (from proposal+design), ## Scope, ## Checklist (items for this group), ## Implementation plan, ## Acceptance criteria, and ## Key references

#### Scenario: Tasks form linear blocker chain
- **WHEN** tasks.json is generated
- **THEN** tasks have sequential blockers (task N+1 blocked by task N) matching the existing littlefactory task validation rules

#### Scenario: No housekeeping merge task generated
- **WHEN** tasks.json is generated regardless of specs_dir configuration
- **THEN** no housekeeping spec merge task is appended -- spec merging is handled by `/lf-archive`

## Boundaries

## Gotchas
- The formalize skill previously generated a conditional housekeeping merge task when `specs_dir` was configured. This was removed because autonomous agents lack the conversation context needed for meaningful spec enrichment. If the housekeeping task pattern resurfaces, it's a sign archive isn't being used.
  (learned: add-lf-archive-skill, 2026-03-28)
