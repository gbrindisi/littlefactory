## Context

The `littlefactory init` command currently only creates a Factoryfile. This leaves users to manually configure:
- Agent instruction files (AGENTS.md / CLAUDE.md)
- Gitignore entries for runtime state files
- Skills that littlefactory provides

The goal is agent-agnostic tooling, with AGENTS.md as the canonical instruction file. Claude Code uses CLAUDE.md, so we need to bridge these with symlinks.

Current codebase structure:
- `cmd/littlefactory/main.go`: Contains `runInit` function
- Skills currently live in `.claude/skills/` in the littlefactory repo itself

## Goals / Non-Goals

**Goals:**
- Single `init` command that fully bootstraps a project for littlefactory
- AGENTS.md as source of truth, with CLAUDE.md symlinked for Claude Code compatibility
- Idempotent gitignore updates
- Embedded skills extracted to `.littlefactory/skills/` with agent symlinks
- Separate `upgrade` command for existing projects
- Verbose logging of all init/upgrade operations

**Non-Goals:**
- Support for agents other than Claude Code (future work)
- Interactive prompts during init (keep it non-interactive)
- Conflict resolution UI for AGENTS.md/CLAUDE.md merge (use concatenation)

## Decisions

### Decision 1: AGENTS.md merge strategy
**Choice**: Concatenate with separator when both AGENTS.md and CLAUDE.md exist

**Rationale**: No data loss, simple implementation, user can clean up later. Alternatives considered:
- Prompt user to choose one: Loses content
- Smart section-based merge: Complex, fragile
- Fail and ask user to resolve: Friction

**Implementation**:
```
if both exist:
    agents_content = read(AGENTS.md)
    claude_content = read(CLAUDE.md)
    merged = agents_content + "\n\n---\n<!-- Merged from CLAUDE.md -->\n\n" + claude_content
    write(AGENTS.md, merged)
    rm(CLAUDE.md)
    symlink(CLAUDE.md -> AGENTS.md)
```

### Decision 2: Skills embedded in binary using Go embed
**Choice**: Use `//go:embed` directive to bundle skills

**Rationale**: Self-contained binary, no external dependencies or downloads. Upgrade gets new skills automatically.

**Structure**:
```go
//go:embed embedded/skills/*
var embeddedSkills embed.FS
```

### Decision 3: Symlink direction
**Choice**: `.claude/skills/<name>` -> `../../.littlefactory/skills/<name>`

**Rationale**: littlefactory owns the skills, agent directories just reference them. Keeps ownership clear.

### Decision 4: Init vs Upgrade separation
**Choice**: Separate commands with different preconditions

**Rationale**:
- `init`: Creates new Factoryfile, fails if exists. For new projects.
- `upgrade`: Requires Factoryfile, applies improvements idempotently. For existing projects.

Alternatives considered:
- Single idempotent init: Confusing semantics, harder to understand what will happen
- Auto-detect mode: Magic behavior is harder to reason about

### Decision 5: Logging format
**Choice**: Numbered steps with indented sub-operations

**Format**:
```
[1/4] Creating Factoryfile
      Created Factoryfile with default configuration
[2/4] Setting up AGENTS.md
      Created AGENTS.md with default content
      Created symlink CLAUDE.md -> AGENTS.md
```

**Rationale**: Clear progress indication, scannable output, identifies which step failed.

### Decision 6: Package organization
**Choice**: New `internal/init/` package with sub-packages

**Structure**:
```
internal/init/
├── init.go           # Orchestrates init steps
├── agentsmd/         # AGENTS.md handling
├── gitignore/        # .gitignore management
├── skills/           # Skill extraction and symlinking
└── upgrade.go        # Upgrade command logic
```

**Rationale**: Separation of concerns, testable units, reusable between init and upgrade.

## Risks / Trade-offs

**Risk**: Symlink handling on Windows
**Mitigation**: Go's `os.Symlink` works on Windows with appropriate permissions. Document requirement for developer mode or admin rights.

**Risk**: Binary size increase from embedded skills
**Mitigation**: Skills are small text files (< 100KB total). Acceptable trade-off for self-contained binary.

**Risk**: User has complex CLAUDE.md/AGENTS.md setup
**Mitigation**: Merge concatenates with clear separator. User can manually refine. No data loss.

**Risk**: Gitignore modifications break user's custom setup
**Mitigation**: Only append, never remove. Check if entries already exist before adding.

**Trade-off**: Init is not idempotent
**Acceptance**: Clear separation between init (new projects) and upgrade (existing projects) makes behavior predictable.
