## Why

The current littlefactory output is an append-only log that scrolls past as the agent runs. Users cannot see task progress at a glance or scroll back through agent output. A split-panel TUI would provide real-time visibility into both task status and agent activity.

## What Changes

- Replace append-only stdout logging with a Bubbletea-based TUI
- Add left panel showing all tasks with status indicators (done/active/pending/blocked)
- Add right panel showing real-time streaming agent output with ANSI color preservation
- Use PTY for agent subprocess to preserve terminal features (colors, spinners)
- Add scrollback support for reviewing agent output history
- Add keyboard navigation (j/k for tasks, standard viewport scrolling for output)

## Capabilities

### New Capabilities

- `tui-display`: Split-panel terminal UI with task list and agent output viewport
- `pty-agent-execution`: Agent subprocess runs in PTY for full terminal emulation

### Modified Capabilities

- `agent-interface`: Agent.Run() accepts io.Writer for streaming output instead of buffering
- `task-source-interface`: Add List() method to retrieve all tasks with status
- `loop-driver`: Driver emits events to TUI instead of printing directly to stdout

## Impact

- `internal/tui/` - New package with Bubbletea model and components
- `internal/agent/claude.go` - PTY integration, streaming output
- `internal/agent/agent.go` - Interface change: Run() gains io.Writer parameter
- `internal/tasks/source.go` - Interface change: add List() method
- `internal/tasks/beads.go` - Implement List() via bd list --json
- `internal/driver/driver.go` - Emit events instead of Print* calls
- `internal/driver/output.go` - May be removed or repurposed
- `cmd/littlefactory/main.go` - Initialize TUI as main event loop
- Dependencies: bubbletea, bubbles, lipgloss, creack/pty, stripansi
