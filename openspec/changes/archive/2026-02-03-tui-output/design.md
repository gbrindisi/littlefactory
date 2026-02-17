## Context

Littlefactory currently outputs an append-only log to stdout. The driver calls Print* functions that write banners, agent output, and summaries sequentially. Agent output is buffered until completion, then printed all at once. Users cannot see task progress at a glance or scroll back through output.

The codebase structure:
- `internal/driver/output.go` - Print functions for banners and summaries
- `internal/driver/driver.go` - Main loop, calls PrintAgentOutput after each iteration
- `internal/agent/claude.go` - Buffers stdout/stderr, returns complete output
- `internal/tasks/beads.go` - Only has Ready(), no way to list all tasks

## Goals / Non-Goals

**Goals:**
- Replace stdout logging with split-panel TUI (tasks left, output right)
- Stream agent output in real-time with ANSI color preservation
- Enable scrollback through agent output history
- Show all tasks with status indicators
- Keyboard navigation (j/k for tasks, viewport scrolling for output)

**Non-Goals:**
- Interactive task management (create/edit/close from TUI)
- Multiple simultaneous agent outputs (parallel execution)
- Persistent output history across runs
- Mouse interaction beyond basic scroll
- Customizable layouts or themes

## Decisions

### Decision 1: Bubbletea for TUI framework

**Choice**: Use charmbracelet/bubbletea with bubbles components.

**Alternatives considered**:
- tview: Has widgets out of box but less flexible for custom behavior
- tcell: Too low-level, more work to build panels

**Rationale**: Bubbletea's Elm architecture fits well with event-driven design. The bubbles library provides viewport (with HighPerformanceRendering for ANSI) and list components. Strong ecosystem and active maintenance.

### Decision 2: PTY for agent subprocess

**Choice**: Run agent subprocess in a PTY (pseudo-terminal) using creack/pty.

**Alternatives considered**:
- io.Writer callback: Simpler but loses terminal features
- Channel-based streaming: More complex, same limitations

**Rationale**: PTY makes the subprocess think it's in a real terminal, so isatty() returns true. This preserves colors, spinners, and progress indicators from claude and other tools. The viewport's HighPerformanceRendering handles ANSI escape sequences.

### Decision 3: Fixed-width left panel

**Choice**: Task panel is fixed at 30 columns.

**Alternatives considered**:
- Percentage-based (25%): More responsive but adds complexity
- User-configurable: Overkill for current needs

**Rationale**: Fixed width is simpler to implement and predictable. 30 cols is enough for task IDs and truncated titles.

### Decision 4: TUI owns the event loop

**Choice**: Bubbletea Program is the main event loop. Driver runs in a goroutine and sends messages via channel.

**Alternatives considered**:
- Driver owns loop, TUI as renderer: Less idiomatic for bubbletea

**Rationale**: Bubbletea is designed around its own event loop. The driver becomes a background worker that emits events (TasksRefreshed, OutputReceived, IterationStarted, etc.).

### Decision 5: Strip ANSI for metadata logging

**Choice**: Raw PTY output goes to TUI, ANSI-stripped output used for metrics.

**Rationale**: Logs should be readable in plain text. Use stripansi library to remove escape codes when calculating OutputLines/OutputBytes for metadata.

### Decision 6: Refresh task list after each iteration

**Choice**: Poll task source after each iteration completes.

**Alternatives considered**:
- Continuous polling: More responsive but adds load
- Watch mode (bd list --watch): Adds complexity

**Rationale**: Simple and sufficient. Task status changes primarily when iterations complete.

## Risks / Trade-offs

**[Risk] PTY complexity on Windows** -> Mitigation: creack/pty handles most cross-platform concerns. Windows support may be limited to Win10+ with ConPTY. Document as known limitation.

**[Risk] Large output buffer memory usage** -> Mitigation: Output buffer grows unbounded during a run. Could add optional max buffer size with oldest content truncation, but defer until it's a real problem.

**[Risk] ANSI escape sequence edge cases** -> Mitigation: viewport's HighPerformanceRendering is designed for this. Known issue with wide runes exists but is rare.

**[Trade-off] Lost stdout compatibility** -> Accept: Previous behavior of logging to stdout (useful for piping/redirection) is lost. Users needing logs should capture agent output separately.

## Component Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         main.go                                          │
│                                                                          │
│  model := tui.New(driver, taskSource)                                   │
│  tea.NewProgram(model).Run()                                            │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Bubbletea Model (tui/tui.go)                        │
│                                                                          │
│  Fields:                          Messages:                              │
│  - tasks []Task                   - TasksRefreshedMsg                    │
│  - activeTaskID string            - OutputMsg{Data []byte}              │
│  - viewport viewport.Model        - IterationStartedMsg                  │
│  - outputBuf bytes.Buffer         - IterationCompleteMsg                 │
│  - autoFollow bool                - RunCompleteMsg                       │
│  - list list.Model                - tea.KeyMsg                           │
│                                   - tea.WindowSizeMsg                    │
└─────────────────────────────────────────────────────────────────────────┘
         │                                             ▲
         │ renders                                     │ events
         ▼                                             │
┌─────────────────────────────────────────────────────────────────────────┐
│                      Driver (runs in goroutine)                          │
│                                                                          │
│  - Iterates through tasks                                                │
│  - Creates PTY for each agent run                                        │
│  - Streams PTY output to channel                                         │
│  - Emits events: started, output, complete                               │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Agent (with PTY)                                    │
│                                                                          │
│  pty, tty := pty.Open()                                                 │
│  cmd.Stdin = tty                                                        │
│  cmd.Stdout = tty                                                       │
│  cmd.Stderr = tty                                                       │
│  go io.Copy(outputWriter, pty)   // streams to TUI                      │
└─────────────────────────────────────────────────────────────────────────┘
```

## File Structure

```
internal/
├── tui/
│   ├── tui.go          # Model, Init, Update, View, message types
│   ├── tasks_panel.go  # Left panel: list component wrapper
│   ├── output_panel.go # Right panel: viewport wrapper
│   ├── status_bar.go   # Bottom bar rendering
│   └── styles.go       # lipgloss style definitions
├── driver/
│   ├── driver.go       # Modified: emit events instead of Print*
│   ├── events.go       # Event types sent to TUI
│   └── output.go       # May be removed or kept for non-TUI mode
├── agent/
│   ├── agent.go        # Interface change: Run gains io.Writer
│   └── claude.go       # PTY integration
└── tasks/
    ├── source.go       # Interface change: add List()
    └── beads.go        # Implement List() via bd list --json
```

## Open Questions

- Should there be a fallback non-TUI mode (--no-tui) for CI/headless environments?
- Should the output panel show iteration separators when a new iteration starts?
