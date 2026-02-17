## Context

Littlefactory currently depends on `bd` (beads), an external task tracker CLI. This dependency has proven unreliable. The goal is to replace it with a simple, self-contained JSON file system.

Current state:
- Tasks stored in `.beads/beads.db` (SQLite via bd)
- Agent template instructs agent to run `bd` commands
- Project detection via `.beads/` directory
- BeadsClient wraps bd CLI calls: `bd ready`, `bd show`, `bd close`, `bd sync`

## Goals / Non-Goals

**Goals:**
- Remove all bd dependency from littlefactory
- Use local JSON file for task storage (`.littlefactory/tasks.json`)
- Driver owns all task state transitions
- Agent focuses only on implementation, not task management
- Maintain existing TaskSource interface pattern for future extensibility

**Non-Goals:**
- Supporting multiple task sources simultaneously
- Task dependencies or blockers (sequential execution makes this unnecessary)
- Failure recovery or retry logic (out of scope for now)
- Migration tool for existing .beads databases

## Decisions

### Decision 1: JSON file format for tasks

**Choice**: Simple JSON with array of task objects

```json
{
  "tasks": [
    {
      "id": "001",
      "title": "Implement feature",
      "description": "Full description...",
      "status": "todo"
    }
  ]
}
```

**Rationale**:
- Human-readable and editable
- No external dependencies (SQLite, daemon)
- Array order defines execution order (sequential)
- Status enum: `todo` | `in_progress` | `done`

**Alternative considered**: YAML - rejected for no real advantage over JSON, JSON is more universal.

### Decision 2: Driver manages all state transitions

**Choice**: Driver sets `in_progress` before iteration, `done` or back to `todo` after.

```
todo ──claim──> in_progress ──success──> done
                    │
                    └──failure──> todo
```

**Rationale**:
- Agent should focus on implementation, not bookkeeping
- Single source of truth for state
- Simplifies agent template significantly

**Alternative considered**: Agent manages status - rejected because it complicates the agent prompt and creates two sources of truth.

### Decision 3: Factoryfile for project detection

**Choice**: Detect project root by looking for `Factoryfile` instead of `.beads/`.

**Rationale**:
- Factoryfile is already required for configuration
- Removes coupling to any specific task backend
- More intuitive marker for "this is a littlefactory project"

### Decision 4: Modify TaskSource interface

**Choice**:
- Remove `Sync()` - no longer needed (JSON written immediately)
- Keep `Ready()`, `List()`, `Show()`, `Close()`
- Add internal claim logic in implementation (not interface change)

**Rationale**:
- Sync was bd-specific (persist to JSONL)
- Claim is implementation detail of JSON source, not interface contract

### Decision 5: openspec-to-lf skill

**Choice**: Rename/adapt `openspec-to-beads` skill to output JSON instead of calling bd.

**Rationale**:
- Preserves the rich context generation (proposal, design, specs embedded in task description)
- Changes only the output format, not the parsing logic

## Risks / Trade-offs

**Risk: Concurrent access to tasks.json**
- Littlefactory is single-process, single-machine
- Mitigation: Not a concern for current use case

**Risk: Process killed mid-execution leaves task as in_progress**
- Mitigation: Manual edit of JSON to reset, or restart picks up where it left off
- Acceptable for MVP, can add startup recovery later if needed

**Risk: JSON file corruption**
- Mitigation: JSON is simple, unlikely to corrupt
- Future: Could add backup before write

**Trade-off: No task dependencies**
- Sequential execution by array order is simpler
- Loses bd's dependency graph feature
- Acceptable because OpenSpec tasks are already ordered sequentially
