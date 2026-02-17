## Why

The `tasks/` directory was legacy from the Python version of littlefactory. Tasks have already moved to `.littlefactory/tasks.json`, but `progress.txt` and `run_metadata.json` still write to `tasks/`. This creates an inconsistent split where runtime state lives in two places. Consolidating all state files into `.littlefactory/` provides a single, configurable location for all littlefactory artifacts.

## What Changes

- Move `progress.txt` and `run_metadata.json` from `tasks/` to `.littlefactory/`
- Rename `progress.txt` to `progress.md` with proper markdown formatting
- Rebrand progress log from "Ciccio" to "Little Factory"
- Add configurable `state_dir` option to Factoryfile (defaults to `.littlefactory`)
- Pass full `*config.Config` to all functions that need state directory path
- Eliminate the `tasks/` directory entirely

## Capabilities

### New Capabilities

(none - this change modifies existing capabilities)

### Modified Capabilities

- `config-management`: Add `state_dir` configuration option with default value `.littlefactory`
- `progress-logging`: Change file location to `<state_dir>/progress.md`, update format to proper markdown, rebrand to "Little Factory"
- `metadata-tracking`: Change file location from `tasks/run_metadata.json` to `<state_dir>/run_metadata.json`
- `json-task-source`: Use `state_dir` from config instead of hardcoded `.littlefactory`

## Impact

- **Code**: `internal/config/config.go`, `internal/driver/progress.go`, `internal/driver/metadata.go`, `internal/tasks/json.go`
- **Tests**: `internal/driver/progress_test.go`, `internal/driver/metadata_test.go`
- **Function signatures**: Functions that take only `projectRoot` will now also take `*config.Config`
- **Breaking**: Existing `tasks/progress.txt` and `tasks/run_metadata.json` files will no longer be updated (new location)
