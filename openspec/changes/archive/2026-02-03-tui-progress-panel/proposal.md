## Why

The TUI right panel currently shows live agent output that resets each iteration. Users lose visibility into past work and have no persistent view of progress. Showing `progress.md` instead provides a cumulative log that persists across iterations and runs.

## What Changes

- Right panel displays `progress.md` content instead of live agent output
- File watching detects updates to `progress.md` and refreshes the panel
- Auto-follow scrolls to bottom when file updates (existing `f` toggle preserved)
- Up/down arrows scroll the progress content
- Remove `j/k` keybindings and cursor tracking (task list becomes display-only)
- Read progress file from configured state directory (`cfg.StateDir`)

## Capabilities

### New Capabilities

- `tui-file-watching`: File watching capability to detect changes to progress.md and trigger panel refresh

### Modified Capabilities

- `tui-display`: Right panel shows progress.md instead of live output; remove cursor navigation and j/k keys; task list is display-only

## Impact

- `internal/tui/tui.go` - Replace output buffer with file content, add file watcher, remove cursor field, pass config for state directory
- `internal/tui/tasks_panel.go` - Remove cursor highlighting logic
- `internal/tui/status_bar.go` - Update keyboard hints (remove j/k mention)
- New dependency: `github.com/fsnotify/fsnotify` for file watching
- `cmd/littlefactory/main.go` - Pass config to TUI constructor
