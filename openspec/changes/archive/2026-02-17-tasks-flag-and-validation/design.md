## Context

The `--change/-c` flag reads tasks.json directly from `openspec/changes/<name>/tasks.json`, but the schema still instructs writing to both locations. The `openspec-to-lf` skill is now obsolete. Users need explicit control over tasks.json location via `--tasks/-t`, and malformed tasks.json files should fail fast with clear errors.

Current flag resolution:
1. `--change <name>` derives path from convention
2. No flag uses default `.littlefactory/tasks.json`

Proposed flag resolution:
1. `--tasks <path>` uses explicit path (highest priority)
2. `--change <name>` derives path from convention
3. No flags use default `.littlefactory/tasks.json`

## Goals / Non-Goals

**Goals:**
- Add `--tasks/-t` flag for explicit path override
- Validate tasks.json structure and sequential chain on load
- Remove obsolete `openspec-to-lf` skill
- Update schema to write tasks.json only to change directory

**Non-Goals:**
- Support parallel task execution (strict sequential model)
- Add validation for non-JSON task sources
- Change the Task struct fields

## Decisions

### Decision 1: Validation happens once at load time

Validate in `NewJSONTaskSourceWithPath` and `NewJSONTaskSource` constructors, not on every read. Returns error if validation fails.

**Rationale**: Fail fast at startup rather than encountering issues mid-run. Avoids repeated validation overhead.

### Decision 2: Validation is a separate function

Create `ValidateTasks(tasks []Task) error` that returns a multi-error with all validation failures. Called by constructors after parsing JSON.

**Rationale**: Testable in isolation, can report all errors at once rather than stopping at first failure.

### Decision 3: Flag priority is -t > -c > default

```
if tasksPath != "" {
    // Use explicit --tasks path
} else if changeName != "" {
    // Derive from --change convention
} else {
    // Use default state_dir/tasks.json
}
```

**Rationale**: Explicit always wins. Convention-based second. Default last.

### Decision 4: Sequential validation uses blocker chain walking

Build a map of task ID to task. Find the root (empty blockers). Walk the chain following blockers. Verify all tasks are visited exactly once.

**Rationale**: Simple linear walk, O(n) complexity, catches gaps/branches/cycles.

## Risks / Trade-offs

- [Breaking change for users with parallel tasks] -> Validation will reject. Users must restructure to sequential. This is intentional - the system only supports sequential execution.
- [Skill removal affects users who installed it] -> The skill was always redundant with the schema's tasks-littlefactory artifact. No migration needed.
