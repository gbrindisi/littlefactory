## Context

The TUI currently streams live agent output to the right panel, clearing it each iteration. The output buffer (`outputBuf`) receives data via `driver.OutputMsg` events. Users cannot see previous iteration output without checking `progress.md` manually.

The state directory is now configurable via `cfg.StateDir` (default `.littlefactory`), so progress.md lives at `<project-root>/<cfg.StateDir>/progress.md`.

## Goals / Non-Goals

**Goals:**
- Display progress.md content in the right panel instead of live agent output
- Auto-refresh panel when progress.md changes on disk
- Preserve existing scroll behavior (up/down/pgup/pgdn) and auto-follow toggle (f)
- Simplify task list to display-only (no cursor navigation)

**Non-Goals:**
- Displaying live agent output (users can watch progress.md updates instead)
- Per-task output storage or retrieval
- Modifying progress.md format

## Decisions

### Decision 1: Use fsnotify for file watching

Watch progress.md using `github.com/fsnotify/fsnotify`. On write events, re-read the file and update the viewport.

**Alternatives considered:**
- Polling: Simpler but wastes CPU and has latency
- inotify directly: Platform-specific, fsnotify abstracts this

**Rationale:** fsnotify is the standard Go solution for cross-platform file watching. Well-maintained, minimal overhead.

### Decision 2: Pass config to TUI constructor

Update `tui.New()` to accept `*config.Config` so it can construct the progress file path from `cfg.StateDir`.

```go
func New(eventChan <-chan interface{}, cfg *config.Config, projectRoot string) *Model
```

**Rationale:** TUI needs to know where progress.md lives. Config already holds `StateDir`.

### Decision 3: Remove cursor and j/k keybindings

Remove `cursor int` field from Model. Remove j/k key handlers. Task list renders without cursor highlighting.

**Rationale:** Cursor served no functional purpose beyond visual selection. With the right panel showing progress.md (not per-task output), cursor navigation is unnecessary.

### Decision 4: File watcher runs in goroutine, sends tea.Msg

Create a `FileChangedMsg` type. Watcher goroutine sends this message when progress.md changes. Update handler re-reads file and refreshes viewport.

```go
type FileChangedMsg struct{}

func watchFile(path string) tea.Cmd {
    return func() tea.Msg {
        // Setup watcher, block until change, return FileChangedMsg
    }
}
```

**Rationale:** Fits bubbletea's message-based architecture. Non-blocking, clean integration.

### Decision 5: Initial load on startup

On `Init()`, read progress.md content immediately (if file exists). Don't wait for first change event.

**Rationale:** Users should see existing progress immediately, not a blank panel.

## Risks / Trade-offs

- **[Risk] File doesn't exist yet** - Handle gracefully with placeholder message until first write
- **[Risk] fsnotify edge cases** - Some editors write temp files then rename. fsnotify handles this but may fire multiple events. Debounce or accept re-reads.
- **[Trade-off] No live streaming** - Users won't see agent output in real-time. Acceptable since progress.md updates at iteration boundaries and provides a cleaner cumulative view.
