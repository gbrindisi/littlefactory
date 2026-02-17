## ADDED Requirements

### Requirement: Skills embedded in binary
The system SHALL embed skill files in the littlefactory binary using Go embed.

#### Scenario: Binary contains embedded skills
- **WHEN** littlefactory binary is built
- **THEN** binary contains all skill directories and SKILL.md files from embedded/skills/

### Requirement: Extract skills to .littlefactory/skills/
The system SHALL extract embedded skills to .littlefactory/skills/ during init.

#### Scenario: Skills extracted on init
- **WHEN** init runs
- **THEN** system creates .littlefactory/skills/ and extracts all embedded skills into it

#### Scenario: Skill directory structure preserved
- **WHEN** embedded skill openspec-to-lf/SKILL.md exists
- **THEN** system creates .littlefactory/skills/openspec-to-lf/SKILL.md

### Requirement: Create symlinks in .claude/skills/ when .claude/ exists
The system SHALL create symlinks from .claude/skills/ to .littlefactory/skills/ when .claude/ directory exists.

#### Scenario: Symlinks created for Claude Code integration
- **WHEN** init runs and .claude/ directory exists
- **THEN** system creates .claude/skills/<name> symlink pointing to ../../.littlefactory/skills/<name> for each skill

#### Scenario: No symlinks when .claude/ does not exist
- **WHEN** init runs and .claude/ directory does not exist
- **THEN** system logs "No .claude/ directory found, skipping Claude Code integration" and creates no symlinks

### Requirement: Skip existing skill symlinks
The system SHALL not overwrite existing files or symlinks in .claude/skills/.

#### Scenario: Skill already exists in .claude/skills/
- **WHEN** init or upgrade runs and .claude/skills/<name> already exists (file or symlink)
- **THEN** system logs "skill already exists, skipping" and does not modify it

### Requirement: Create .claude/skills/ directory if needed
The system SHALL create .claude/skills/ directory when creating symlinks.

#### Scenario: .claude/ exists but .claude/skills/ does not
- **WHEN** init runs with .claude/ but no .claude/skills/
- **THEN** system creates .claude/skills/ directory before creating symlinks
