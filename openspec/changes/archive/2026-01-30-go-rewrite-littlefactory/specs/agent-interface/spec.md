## ADDED Requirements

### Requirement: Agent interface defines execution contract
The system SHALL define an Agent interface that abstracts autonomous agent execution.

#### Scenario: Agent executes with context
- **WHEN** Driver calls Agent.Run() with prompt and session ID
- **THEN** Agent executes and returns AgentResult with exit code and output

#### Scenario: Agent respects context timeout
- **WHEN** Context deadline is exceeded during execution
- **THEN** Agent terminates and returns context.DeadlineExceeded error

### Requirement: Claude Code agent implementation
The system SHALL provide a Claude Code implementation of the Agent interface.

#### Scenario: Claude invocation with correct flags
- **WHEN** Claude agent runs
- **THEN** System executes `claude --dangerously-skip-permissions --print --session-id <uuid>` with prompt via stdin

#### Scenario: Session path computation
- **WHEN** Claude agent runs with session ID
- **THEN** AgentResult includes computed session path at ~/.claude/projects/<encoded-root>/<session-id>.jsonl

#### Scenario: Output capture
- **WHEN** Claude agent completes
- **THEN** AgentResult includes stdout+stderr combined output with line and byte counts
