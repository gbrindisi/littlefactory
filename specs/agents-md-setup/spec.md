# agents-md-setup

## What It Does
AGENTS.md serves as the single source of truth for agent instructions. During init, the system creates AGENTS.md, migrates content from CLAUDE.md if present, and ensures CLAUDE.md is a symlink to AGENTS.md for backward compatibility.

## Requirements

### Requirement: Create AGENTS.md as source of truth
The system SHALL create AGENTS.md as the canonical agent instruction file when it does not exist.

#### Scenario: AGENTS.md created in empty directory
- **WHEN** init runs in a directory without AGENTS.md or CLAUDE.md
- **THEN** system creates AGENTS.md with default content

### Requirement: Migrate CLAUDE.md to AGENTS.md
The system SHALL migrate existing CLAUDE.md content to AGENTS.md and replace CLAUDE.md with a symlink.

#### Scenario: CLAUDE.md exists without AGENTS.md
- **WHEN** init runs in a directory with CLAUDE.md but no AGENTS.md
- **THEN** system renames CLAUDE.md to AGENTS.md and creates symlink CLAUDE.md -> AGENTS.md

### Requirement: Merge AGENTS.md and CLAUDE.md when both exist
The system SHALL merge content from both files when both AGENTS.md and CLAUDE.md exist.

#### Scenario: Both AGENTS.md and CLAUDE.md exist
- **WHEN** init runs in a directory with both AGENTS.md and CLAUDE.md
- **THEN** system appends CLAUDE.md content to AGENTS.md with separator, removes CLAUDE.md, and creates symlink CLAUDE.md -> AGENTS.md

#### Scenario: Merge preserves content from both files
- **WHEN** AGENTS.md contains "content A" and CLAUDE.md contains "content B"
- **THEN** merged AGENTS.md contains both "content A" and "content B" separated by a marker

### Requirement: Skip AGENTS.md setup when already configured
The system SHALL skip AGENTS.md setup if CLAUDE.md is already a symlink to AGENTS.md.

#### Scenario: CLAUDE.md already symlinked to AGENTS.md
- **WHEN** init or upgrade runs and CLAUDE.md is a symlink pointing to AGENTS.md
- **THEN** system logs "already configured" and makes no changes

### Requirement: Default AGENTS.md content
The system SHALL use a minimal default content for AGENTS.md when creating a new file.

#### Scenario: Default content structure
- **WHEN** AGENTS.md is created without migrating from CLAUDE.md
- **THEN** file contains header "Agent Instructions" and reference to littlefactory task management

## Boundaries

## Gotchas
