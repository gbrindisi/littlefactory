## 1. Dependencies and Setup

- [x] 1.1 Add bubbletea dependency (github.com/charmbracelet/bubbletea)
- [x] 1.2 Add bubbles dependency (github.com/charmbracelet/bubbles)
- [x] 1.3 Add lipgloss dependency (github.com/charmbracelet/lipgloss)
- [x] 1.4 Add creack/pty dependency (github.com/creack/pty)
- [x] 1.5 Add stripansi dependency (github.com/acarl005/stripansi)
- [x] 1.6 Create internal/tui/ package directory structure

## 2. TaskSource Interface Extension

- [x] 2.1 Add List() method to TaskSource interface in source.go
- [x] 2.2 Implement List() in BeadsClient using bd list --json -n 0 --all
- [x] 2.3 Add unit tests for BeadsClient.List()

## 3. Agent Interface and PTY Integration

- [x] 3.1 Update Agent interface: Run() takes io.Writer parameter
- [x] 3.2 Implement PTY creation in ConfigurableAgent.Run()
- [x] 3.3 Stream PTY output to provided io.Writer
- [x] 3.4 Implement ANSI stripping for OutputLines calculation
- [x] 3.5 Update agent tests for new interface signature

## 4. Driver Event System

- [x] 4.1 Create events.go with message types (RunStarted, IterationStarted, OutputMsg, etc.)
- [x] 4.2 Add event channel to Driver struct
- [x] 4.3 Modify Driver.Run() to emit events instead of calling Print* functions
- [x] 4.4 Modify Driver.RunIteration() to pass io.Writer to agent and emit output events
- [x] 4.5 Update driver tests for event-based behavior

## 5. TUI Core Model

- [x] 5.1 Create tui/tui.go with Model struct (tasks, viewport, outputBuf, autoFollow)
- [x] 5.2 Implement Init() to subscribe to driver event channel
- [x] 5.3 Implement Update() for all message types (TasksRefreshed, Output, IterationStarted, etc.)
- [x] 5.4 Implement View() with two-panel layout using lipgloss.JoinHorizontal

## 6. TUI Components

- [x] 6.1 Create tui/styles.go with lipgloss style definitions
- [x] 6.2 Create tui/tasks_panel.go with list component wrapper and status icons
- [x] 6.3 Create tui/output_panel.go with viewport wrapper (HighPerformanceRendering enabled)
- [x] 6.4 Create tui/status_bar.go with iteration count and keyboard hints

## 7. Keyboard Navigation

- [x] 7.1 Implement j/k navigation for task list
- [x] 7.2 Implement viewport scrolling (up/down/pgup/pgdn forwarded to viewport)
- [x] 7.3 Implement f key to toggle auto-follow mode
- [x] 7.4 Implement q/Ctrl+C to quit gracefully

## 8. Main Integration

- [x] 8.1 Update cmd/littlefactory/main.go to create TUI model
- [x] 8.2 Start driver in goroutine with event channel
- [x] 8.3 Run tea.NewProgram(model) as main event loop
- [x] 8.4 Handle graceful shutdown on TUI exit

## 9. Cleanup

- [x] 9.1 Remove or deprecate Print* functions in output.go
- [x] 9.2 Update any remaining references to old output functions
- [x] 9.3 Verify all tests pass with new TUI-based architecture
