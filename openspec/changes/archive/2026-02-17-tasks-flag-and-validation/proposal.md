## Why

The `--change` flag already reads tasks.json from the change directory, making the copy to `.littlefactory/tasks.json` redundant. We need a `--tasks/-t` flag for explicit path override, and validation to catch malformed tasks.json files early with clear error messages.

## What Changes

- Add `--tasks/-t` flag to `run` command for explicit tasks.json path specification
- Add tasks.json validation on load: required fields, valid status values, unique IDs, and strict sequential blocker chain
- Update littlefactory schema to remove duplicate write to `.littlefactory/tasks.json`
- Remove obsolete `openspec-to-lf` skill from embedded skills

## Capabilities

### New Capabilities

- `tasks-validation`: Validation of tasks.json structure and sequential blocker chain on load

### Modified Capabilities

- `change-flag`: Add `--tasks/-t` flag with priority over `--change` and default resolution
- `json-task-source`: Add validation on load, explicit path constructor validation

## Impact

- `cmd/littlefactory/main.go`: Add `--tasks` flag, update flag priority logic
- `internal/tasks/json.go`: Add validation logic, update constructors
- `internal/init/skills/embedded/skills/openspec-to-lf/`: Remove directory
- `openspec/schemas/littlefactory/schema.yaml`: Update tasks-littlefactory artifact instruction
- `.littlefactory/skills/openspec-to-lf/`: Remove from installed skills
