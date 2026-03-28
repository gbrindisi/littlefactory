# skill-installation

## What It Does
Embeds skill files (including lf:* skills) in the littlefactory binary, extracts them to .littlefactory/skills/ during init, and creates symlinks from .claude/skills/ for Claude Code integration.

## Requirements
### Requirement: Skills embedded in binary
The system SHALL embed skill files in the littlefactory binary using Go embed, including the four `/lf:*` skills (lf-explore, lf-formalize, lf-do, lf-verify).

#### Scenario: Binary contains embedded skills
- **WHEN** littlefactory binary is built
- **THEN** binary contains all skill directories and SKILL.md files from embedded/skills/, including lf-explore, lf-formalize, lf-do, and lf-verify

### Requirement: Extract skills to .littlefactory/skills/
The system SHALL extract embedded skills to .littlefactory/skills/ during init.

#### Scenario: Skills extracted on init
- **WHEN** init runs
- **THEN** system creates .littlefactory/skills/ and extracts all embedded skills including lf-explore/, lf-formalize/, lf-do/, lf-verify/ each containing SKILL.md

#### Scenario: Skill directory structure preserved
- **WHEN** embedded skill lf-formalize/SKILL.md exists
- **THEN** system creates .littlefactory/skills/lf-formalize/SKILL.md

### Requirement: Create symlinks in .claude/skills/ when .claude/ exists
The system SHALL create symlinks from .claude/skills/ to .littlefactory/skills/ when .claude/ directory exists.

#### Scenario: Symlinks created for Claude Code integration
- **WHEN** init runs and .claude/ directory exists
- **THEN** system creates .claude/skills/<name> symlink pointing to ../../.littlefactory/skills/<name> for each skill

#### Scenario: No symlinks when .claude/ does not exist
- **WHEN** init runs and .claude/ directory does not exist
- **THEN** system logs "No .claude/ directory found, skipping Claude Code integration" and creates no symlinks

### Requirement: Create .claude/skills/ directory if needed
The system SHALL create .claude/skills/ directory when creating symlinks.

#### Scenario: .claude/ exists but .claude/skills/ does not
- **WHEN** init runs with .claude/ but no .claude/skills/
- **THEN** system creates .claude/skills/ directory before creating symlinks

## Boundaries
- ALWAYS: Update embed_test.go and upgrade_test.go when adding/removing embedded skills

## Gotchas
- Adding a new embedded skill requires updating exactly three files: the SKILL.md itself, `embed_test.go` (extraction count), and `upgrade_test.go` (symlink count). Miss one and tests fail.
  (learned: add-lf-archive-skill, 2026-03-28)
