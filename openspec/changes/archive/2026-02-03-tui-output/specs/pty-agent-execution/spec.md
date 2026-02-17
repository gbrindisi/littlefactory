## ADDED Requirements

### Requirement: Agent subprocess runs in PTY
The system SHALL execute agent subprocesses in a pseudo-terminal (PTY) to preserve terminal features.

#### Scenario: PTY allocation
- **WHEN** agent Run() is called
- **THEN** system creates a PTY pair and attaches slave to subprocess stdin/stdout/stderr

#### Scenario: Terminal detection by subprocess
- **WHEN** subprocess checks isatty() on stdout
- **THEN** isatty() returns true (subprocess believes it's in a real terminal)

#### Scenario: Output streaming from PTY master
- **WHEN** subprocess writes to stdout/stderr
- **THEN** output is read from PTY master and streamed to provided io.Writer

### Requirement: ANSI escape sequence passthrough
The system SHALL pass through ANSI escape sequences from PTY output unchanged.

#### Scenario: Color codes preserved
- **WHEN** subprocess outputs ANSI color codes
- **THEN** color codes are included in output sent to io.Writer

#### Scenario: Cursor movement codes preserved
- **WHEN** subprocess outputs cursor movement sequences
- **THEN** sequences are included in output (TUI viewport handles rendering)

### Requirement: ANSI stripping for metrics
The system SHALL strip ANSI escape sequences when calculating output metrics.

#### Scenario: Line count excludes escape sequences
- **WHEN** AgentResult.OutputLines is calculated
- **THEN** ANSI escape sequences are stripped before counting newlines

#### Scenario: Byte count uses raw output
- **WHEN** AgentResult.OutputBytes is calculated
- **THEN** raw byte count includes escape sequences (for accurate transfer size)
