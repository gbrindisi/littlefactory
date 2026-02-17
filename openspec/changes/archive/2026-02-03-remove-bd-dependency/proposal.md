## Why

The `bd` (beads) task tracker is an external dependency that has proven unreliable. Littlefactory should use a simple, local JSON-based task system that is self-contained and predictable.

## What Changes

- Replace BeadsClient task source with a JSON file-based implementation
- Tasks stored in `.littlefactory/tasks.json` instead of `.beads/beads.db`
- Driver manages task state transitions (claim on start, done/reset on completion/failure)
- Agent template no longer includes bd commands - agent focuses on implementation only
- Project root detection via `Factoryfile` instead of `.beads/` directory
- New Claude skill `openspec-to-lf` replaces `openspec-to-beads` for task generation
- **BREAKING**: Remove bd CLI dependency and all bd-related code

## Capabilities

### New Capabilities

- `json-task-source`: JSON file-based TaskSource implementation that reads/writes `.littlefactory/tasks.json`

### Modified Capabilities

- `task-source-interface`: Remove Sync() method (no longer needed), add Claim() for status transition
- `loop-driver`: Driver claims task before iteration, marks done on success, resets to todo on failure
- `project-detection`: Detect project root via Factoryfile instead of .beads directory
- `template-system`: Remove bd commands from embedded template, simplify agent workflow

## Impact

- **Code**: `internal/tasks/beads.go` removed, new `internal/tasks/json.go` added
- **Code**: `internal/driver/driver.go` updated for claim/done/reset logic
- **Code**: `internal/config/project.go` updated for Factoryfile detection
- **Code**: `internal/template/embedded/CLAUDE.md` simplified
- **Code**: `cmd/littlefactory/main.go` removes bd CLI check, uses JSON source
- **Skills**: `.claude/skills/openspec-to-beads/` renamed to `openspec-to-lf/`
- **Dependencies**: No longer requires `bd` binary in PATH
- **Data**: New `.littlefactory/` directory replaces `.beads/`
