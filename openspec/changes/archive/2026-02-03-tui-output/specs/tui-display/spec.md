## ADDED Requirements

### Requirement: Split-panel terminal interface
The system SHALL display a two-panel TUI with task list on the left and agent output on the right.

#### Scenario: Initial layout
- **WHEN** TUI starts
- **THEN** system displays left panel (30 cols fixed) with task list and right panel (remaining width) with output viewport

#### Scenario: Window resize
- **WHEN** terminal window is resized
- **THEN** system adjusts right panel width while keeping left panel at 30 cols

### Requirement: Task list panel with status indicators
The system SHALL display all tasks with visual status indicators.

#### Scenario: Task status display
- **WHEN** task list is rendered
- **THEN** each task shows status icon: [x] for closed, [>] for active, [ ] for pending, [!] for blocked

#### Scenario: Active task highlighting
- **WHEN** an iteration starts for a task
- **THEN** that task is marked with [>] indicator and visually highlighted

#### Scenario: Task list scrolling
- **WHEN** task count exceeds panel height
- **THEN** task list is scrollable with j/k keys

### Requirement: Output viewport with real-time streaming
The system SHALL display agent output in real-time with ANSI color preservation.

#### Scenario: Output streaming
- **WHEN** agent produces output
- **THEN** output appears in right panel immediately (not buffered until completion)

#### Scenario: ANSI color preservation
- **WHEN** agent output contains ANSI escape sequences
- **THEN** colors and formatting are rendered correctly in viewport

#### Scenario: Output scrollback
- **WHEN** user scrolls up in output panel
- **THEN** historical output is accessible via Up/Down/PgUp/PgDn keys

#### Scenario: Auto-follow mode
- **WHEN** new output arrives and auto-follow is enabled
- **THEN** viewport scrolls to bottom automatically

#### Scenario: Toggle auto-follow
- **WHEN** user presses 'f' key
- **THEN** auto-follow mode toggles on/off

### Requirement: Status bar displays run progress
The system SHALL display a status bar at the bottom showing run progress.

#### Scenario: Status bar content
- **WHEN** run is in progress
- **THEN** status bar shows iteration count, task counts (done/active/pending), and keyboard hints

### Requirement: Keyboard navigation
The system SHALL support keyboard navigation for tasks and output.

#### Scenario: Task navigation with j/k
- **WHEN** user presses j or k
- **THEN** task list selection moves down or up respectively

#### Scenario: Quit with q
- **WHEN** user presses q or Ctrl+C
- **THEN** TUI exits gracefully

### Requirement: Output buffer cleared per iteration
The system SHALL clear the output buffer when a new iteration starts.

#### Scenario: New iteration clears output
- **WHEN** a new iteration starts
- **THEN** output viewport is cleared and shows output from new iteration only
