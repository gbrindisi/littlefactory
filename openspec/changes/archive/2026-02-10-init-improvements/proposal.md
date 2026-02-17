## Why

The `littlefactory init` command currently only creates a Factoryfile. It misses key setup steps: agent instruction files (AGENTS.md), gitignore entries for runtime files, and skill installation. Users must manually configure these, creating friction and inconsistency across projects.

## What Changes

- Init command creates AGENTS.md as the source of truth for agent instructions
- Init command handles CLAUDE.md migration (copy content, symlink to AGENTS.md)
- Init command merges AGENTS.md and CLAUDE.md when both exist
- Init command updates .gitignore with `.littlefactory/run_metadata.json` and `.littlefactory/tasks.json`
- Init command creates `.littlefactory/skills/` and installs embedded skills
- Init command symlinks skills to `.claude/skills/` when `.claude/` directory exists
- New `upgrade` command for existing projects to apply new init improvements
- All init/upgrade steps are logged to stdout with clear progress indicators

## Capabilities

### New Capabilities
- `agents-md-setup`: Manages AGENTS.md as source of truth with CLAUDE.md symlink and merge handling
- `gitignore-management`: Idempotent .gitignore updates for littlefactory runtime files
- `skill-installation`: Embeds skills in binary and extracts to `.littlefactory/skills/` with agent directory symlinks
- `upgrade-command`: Upgrades existing projects with new init improvements (idempotent)

### Modified Capabilities
- `init-command`: Extended to orchestrate AGENTS.md setup, gitignore updates, and skill installation

## Impact

- `cmd/littlefactory/main.go`: Extended init command, new upgrade command
- `internal/`: New packages for agents-md, gitignore, and skill management
- Binary size: Increases due to embedded skill files
- `.gitignore`: Will be modified in user projects
- `AGENTS.md` / `CLAUDE.md`: Will be created/modified in user projects
- `.littlefactory/skills/`: New directory created in user projects
- `.claude/skills/`: Symlinks created when `.claude/` exists
