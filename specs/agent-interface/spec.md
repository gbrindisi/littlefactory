# agent-interface

## What It Does
The Agent interface abstracts agent execution, defining the contract for running autonomous agents and capturing their results. A configurable implementation allows any command to be used as an agent, and the run command selects agents by name from configuration.

## Requirements

### Requirement: Agent interface defines execution contract
The system SHALL define an Agent interface that abstracts autonomous agent execution.

#### Scenario: Agent executes with context
- **WHEN** Driver calls Agent.Run() with prompt
- **THEN** Agent executes and returns AgentResult with exit code and output

#### Scenario: Agent respects context timeout
- **WHEN** Context deadline is exceeded during execution
- **THEN** Agent terminates and returns context.DeadlineExceeded error

### Requirement: Configurable agent implementation
The system SHALL provide a configurable agent that executes a user-specified command.

#### Scenario: Agent invocation with configured command
- **WHEN** Agent runs with command "claude --dangerously-skip-permissions --print"
- **THEN** System executes the exact command string with prompt via stdin

#### Scenario: Agent invocation with custom command
- **WHEN** Agent runs with command "/path/to/custom-agent --flag"
- **THEN** System executes the exact command string with prompt via stdin

#### Scenario: Output capture
- **WHEN** Agent completes
- **THEN** AgentResult includes stdout+stderr combined output with line and byte counts

### Requirement: Run command with agent selection
The system SHALL provide a `run` command that executes the autonomous loop with a specified agent.

#### Scenario: Run with explicit agent name
- **WHEN** User runs `littlefactory run claude`
- **THEN** System loads claude agent from config and starts autonomous loop

#### Scenario: Run with default agent
- **WHEN** User runs `littlefactory run` without agent name
- **THEN** System uses default_agent from config

#### Scenario: Run with unknown agent
- **WHEN** User runs `littlefactory run nonexistent`
- **THEN** System prints error that agent is not configured and exits with non-zero code

## Boundaries

## Gotchas
