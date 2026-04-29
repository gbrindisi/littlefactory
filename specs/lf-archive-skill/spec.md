# lf-archive-skill

## What It Does
The `/lf-archive` skill merges delta specs from a completed change into the project's long-lived specs directory, enriches specs with boundaries and gotchas through an interactive conversation, and optionally cleans up the change directory. It is the final step in the change lifecycle, turning implementation learnings into persistent project knowledge.

## Requirements

### Requirement: Archive skill merges delta specs into specs_dir
The system SHALL provide an embedded `/lf-archive` skill that reads delta specs from a change directory and merges them into the project's `specs_dir`.

#### Scenario: Archive with specs_dir configured
- **WHEN** the user invokes `/lf-archive` and `specs_dir` is configured in Factoryfile
- **THEN** the skill reads delta specs from `.littlefactory/changes/<name>/specs/` and merges them into the configured `specs_dir`

#### Scenario: Archive without specs_dir configured
- **WHEN** the user invokes `/lf-archive` and `specs_dir` is not configured in Factoryfile
- **THEN** the skill asks the user where to write specs, or skips the merge step entirely if the user declines

#### Scenario: ADDED specs are copied as new files
- **WHEN** a delta spec contains only `## ADDED Requirements` and no corresponding spec exists in `specs_dir`
- **THEN** the skill creates a new spec file in `specs_dir/<capability>/spec.md` using the enriched format

#### Scenario: MODIFIED specs are merged into existing files
- **WHEN** a delta spec contains `## MODIFIED Requirements` and a corresponding spec exists in `specs_dir`
- **THEN** the skill replaces the matching requirement blocks in the existing spec with the updated content

#### Scenario: REMOVED specs are deleted from existing files
- **WHEN** a delta spec contains `## REMOVED Requirements`
- **THEN** the skill removes the matching requirement blocks from the existing spec and captures the removal reason as a gotcha

### Requirement: Archive enriches specs with boundaries and gotchas
The system SHALL have the archive skill conduct an interactive gotcha-mining conversation to enrich specs beyond the lean delta format.

#### Scenario: Gotcha mining from verify findings
- **WHEN** a verify report was generated in the current conversation
- **THEN** the skill surfaces verify warnings and suggestions as candidate gotchas and asks the user which to persist

#### Scenario: Gotcha mining from implementation experience
- **WHEN** the archive skill has conversation context from the change's implementation
- **THEN** the skill proposes gotchas based on edge cases, workarounds, and surprises encountered during implementation

#### Scenario: Boundary discovery from implementation
- **WHEN** the archive skill reviews the change's design decisions and implementation
- **THEN** the skill proposes boundaries (always/ask/never rules) for the affected capabilities

#### Scenario: User prunes proposed enrichments
- **WHEN** the skill proposes gotchas or boundaries
- **THEN** the user can accept, reject, or edit each proposed enrichment before it is written

### Requirement: Enriched spec format
The system SHALL write enriched specs using a canonical format with four sections: What It Does, Requirements, Boundaries, Gotchas.

#### Scenario: New enriched spec created
- **WHEN** a new spec is created during archive
- **THEN** the file contains `# <capability-name>`, `## What It Does` (one paragraph), `## Requirements` (SHALL/MUST + WHEN/THEN scenarios), `## Boundaries`, and `## Gotchas`

#### Scenario: Existing spec enriched with new sections
- **WHEN** an existing spec that only has Requirements is merged during archive
- **THEN** the skill adds Boundaries and Gotchas sections while preserving existing Requirements content

### Requirement: Archive skill infers change name
The system SHALL have the archive skill infer the change name from conversation context or prompt for selection if ambiguous.

#### Scenario: Single active change
- **WHEN** the user invokes `/lf-archive` and exactly one change exists under `.littlefactory/changes/`
- **THEN** the skill uses that change automatically

#### Scenario: Multiple active changes
- **WHEN** the user invokes `/lf-archive` and multiple changes exist
- **THEN** the skill prompts the user to select which change to archive

### Requirement: Archive optionally cleans up change directory
The system SHALL offer to delete the change directory after successful spec merge.

#### Scenario: User accepts cleanup
- **WHEN** spec merge is complete and the user confirms cleanup
- **THEN** the skill deletes `.littlefactory/changes/<name>/`

#### Scenario: User declines cleanup
- **WHEN** spec merge is complete and the user declines cleanup
- **THEN** the change directory is left intact

## Boundaries

## Gotchas
- The archive skill is the only skill in the workflow that writes to `specs_dir` -- all other skills write to `.littlefactory/changes/`. If `specs_dir` is not configured, the user must be prompted rather than failing silently.
  (learned: add-lf-archive-skill, 2026-03-28)
