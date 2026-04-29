# lf-verify-skill

## What It Does
The `/lf-verify` skill validates that a change's implementation matches its artifacts (specs, tasks, design) across three dimensions: completeness, correctness, and coherence. It produces an actionable report with prioritized issues and suggests archiving when the implementation passes.

## Requirements

### Requirement: Verify checks three dimensions
The system SHALL have the verify skill check implementation against change artifacts across completeness, correctness, and coherence.

#### Scenario: Completeness check
- **WHEN** verify runs
- **THEN** it checks that all tasks are marked done and all spec requirements have corresponding implementation evidence in the codebase

#### Scenario: Correctness check
- **WHEN** verify runs
- **THEN** it checks that implementations match requirement intent and that scenarios from delta specs are covered by code or tests

#### Scenario: Coherence check
- **WHEN** verify runs
- **THEN** it checks that implementation follows design.md decisions (if present) and that new code is consistent with project patterns

### Requirement: Verify produces actionable report
The system SHALL have the verify skill produce a report with CRITICAL, WARNING, and SUGGESTION issues, each with specific file references and recommendations. When all checks pass or only non-critical issues remain, the final assessment SHALL suggest running `/lf-archive`.

#### Scenario: Report with critical issues
- **WHEN** verify finds critical issues (incomplete tasks, missing requirement implementations)
- **THEN** the report lists each critical issue with a specific, actionable recommendation and file references

#### Scenario: Report with no critical issues suggests archive
- **WHEN** verify finds no critical issues
- **THEN** the final assessment includes "Ready for archive. Run `/lf-archive` to merge specs and capture learnings."

#### Scenario: Report with all checks passed suggests archive
- **WHEN** verify finds all tasks done, requirements implemented, and design followed
- **THEN** report states "All checks passed. Run `/lf-archive` to merge specs and capture learnings."

### Requirement: Verify degrades gracefully
The system SHALL have the verify skill adapt its checks based on which artifacts are available, rather than failing when artifacts are missing.

#### Scenario: Tasks-only change
- **WHEN** verify runs on a change that only has tasks.json (no specs or design)
- **THEN** the skill verifies task completion only, skips spec and design checks, and notes which checks were skipped

## Boundaries

## Gotchas
- The verify skill checks task completion via `"status"` fields in tasks.json, not markdown checkboxes. The earlier openspec verify parsed `- [ ]`/`- [x]` in tasks.md -- that format no longer exists in littlefactory.
  (learned: add-lf-archive-skill, 2026-03-28)
