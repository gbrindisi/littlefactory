# tui-display Specification

## Purpose
Terminal user interface for displaying task progress and status during littlefactory runs.

## Requirements
### Requirement: Split-panel terminal interface
The system SHALL display a two-panel TUI with task list on the left and progress log on the right.

#### Scenario: Initial layout
- **WHEN** TUI starts
- **THEN** system displays left panel (30 cols fixed) with task list and right panel (remaining width) with progress.md content

#### Scenario: Window resize
- **WHEN** terminal window is resized
- **THEN** system adjusts right panel width while keeping left panel at 30 cols

### Requirement: Task list panel with status indicators
The system SHALL display all tasks with visual status indicators (display-only, no cursor).

#### Scenario: Task status display
- **WHEN** task list is rendered
- **THEN** each task shows status icon: [x] for closed, [>] for active, [ ] for pending, [!] for blocked

#### Scenario: Active task highlighting
- **WHEN** an iteration starts for a task
- **THEN** that task is marked with [>] indicator and visually highlighted

### Requirement: Progress viewport with file content
The system SHALL display progress.md content in the right panel with auto-follow on file updates.

#### Scenario: Progress content display
- **WHEN** progress.md exists
- **THEN** right panel displays full file content in scrollable viewport

#### Scenario: Auto-follow on file update
- **WHEN** progress.md is updated and auto-follow is enabled
- **THEN** viewport scrolls to bottom automatically

#### Scenario: Manual scrolling disables auto-follow
- **WHEN** user scrolls up in progress panel
- **THEN** auto-follow is disabled until re-enabled with 'f' key

#### Scenario: Toggle auto-follow
- **WHEN** user presses 'f' key
- **THEN** auto-follow mode toggles on/off

### Requirement: Status bar displays run progress
The system SHALL display a status bar at the bottom showing run progress.

#### Scenario: Status bar content
- **WHEN** run is in progress
- **THEN** status bar shows iteration count, task counts (done/active/pending), and keyboard hints

### Requirement: Keyboard navigation for progress panel
The system SHALL support keyboard navigation for scrolling progress content.

#### Scenario: Scroll with arrow keys
- **WHEN** user presses up or down arrow
- **THEN** progress panel scrolls one line up or down respectively

#### Scenario: Page scroll
- **WHEN** user presses PgUp or PgDn
- **THEN** progress panel scrolls one page up or down

#### Scenario: Quit with q
- **WHEN** user presses q or Ctrl+C
- **THEN** TUI exits gracefully

