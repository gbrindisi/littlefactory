## 1. Remove TUI package

- [x] 1.1 Delete `internal/tui/` directory (tui.go, tasks_panel.go, output_panel.go, status_bar.go, styles.go, and all tests)

## 2. Remove driver event system

- [x] 2.1 Delete `internal/driver/events.go` (all event message types)
- [x] 2.2 Remove `eventChan` field, `emit()` method, and `newOutputWriter` from `internal/driver/driver.go`
- [x] 2.3 Remove event channel parameter from `NewDriver` constructor
- [x] 2.4 Remove all `d.emit(...)` calls from `Run()` and `RunIteration()`
- [x] 2.5 Update driver tests to remove event channel references and mocks

## 3. Add status line output to driver

- [x] 3.1 Add iteration start status line: `[N/MAX] Starting: <title> (<id>)` printed to stdout before agent execution
- [x] 3.2 Add iteration complete status line: `[N/MAX] Completed`, `[N/MAX] Failed: <error>`, or `[N/MAX] Timed out` after agent execution
- [x] 3.3 Add run summary line: `Run complete: <status> (N/MAX iterations)` or `Run cancelled` at end of run
- [x] 3.4 Add tests for status line output

## 4. Simplify main.go run command

- [ ] 4.1 Remove bubbletea/TUI imports and event channel creation from `runRun`
- [ ] 4.2 Replace TUI event loop with synchronous `d.Run(ctx)` call on main goroutine
- [ ] 4.3 Wire signal handling (SIGINT/SIGTERM) to cancel context directly

## 5. Remove dependencies

- [ ] 5.1 Run `go mod tidy` to remove bubbletea, bubbles, lipgloss, fsnotify and transitive deps
- [ ] 5.2 Verify build and tests pass
