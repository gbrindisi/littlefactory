## MODIFIED Requirements

### Requirement: Agent interface defines execution contract
The system SHALL define an Agent interface that abstracts autonomous agent execution.

#### Scenario: Agent executes with streaming output
- **WHEN** Driver calls Agent.Run() with prompt and io.Writer
- **THEN** Agent executes, streams output to writer, and returns AgentResult with exit code and metrics

#### Scenario: Agent respects context timeout
- **WHEN** Context deadline is exceeded during execution
- **THEN** Agent terminates and returns context.DeadlineExceeded error

#### Scenario: Output streaming during execution
- **WHEN** Agent subprocess produces output
- **THEN** Output is written to provided io.Writer in real-time (not buffered)

## ADDED Requirements

### Requirement: ConfigurableAgent implementation
The system SHALL provide a ConfigurableAgent implementation of the Agent interface.

#### Scenario: Command execution with PTY
- **WHEN** ConfigurableAgent runs
- **THEN** System executes configured command in PTY with prompt via stdin

#### Scenario: Output streaming to writer
- **WHEN** ConfigurableAgent executes
- **THEN** PTY output is copied to provided io.Writer in real-time

#### Scenario: Metrics calculation
- **WHEN** ConfigurableAgent completes
- **THEN** AgentResult includes OutputLines (ANSI-stripped) and OutputBytes (raw)

## REMOVED Requirements

### Requirement: Claude Code agent implementation
**Reason**: Replaced by ConfigurableAgent which can execute any command including claude
**Migration**: Use ConfigurableAgent with command "claude --dangerously-skip-permissions --print"
