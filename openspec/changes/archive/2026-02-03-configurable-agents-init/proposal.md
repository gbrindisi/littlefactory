## Why

The agent executable is hardcoded to `claude`, preventing users from using alternative agents or custom wrappers. Additionally, users must manually create Factoryfile configuration, with no scaffolding command to bootstrap new projects.

## What Changes

- Add `littlefactory init` command that creates a Factoryfile with default settings
- Replace `start` command with `run <agent-name>` command
- Support multiple named agents in configuration via `agents:` map
- Add `default_agent` configuration to specify which agent to use when none specified
- Make agent command fully configurable (executable + arguments)
- **BREAKING**: Remove `start` command, replace with `run`
- **BREAKING**: Change `agent:` config key to `agents:` (plural, map structure)
- Remove session ID management from agent interface (simplification)

## Capabilities

### New Capabilities
- `init-command`: CLI command to scaffold a new Factoryfile with default configuration

### Modified Capabilities
- `config-management`: Add `agents` map and `default_agent` field, remove singular `agent` section
- `agent-interface`: Remove session ID from interface, make command configurable instead of hardcoded

## Impact

- `cmd/littlefactory/main.go`: Add init command, rename start to run, add agent name argument
- `internal/config/config.go`: New config structure with agents map and default_agent
- `internal/agent/agent.go`: Remove sessionID from Agent interface
- `internal/agent/claude.go`: Accept command from config instead of hardcoding, remove session path logic
- `internal/driver/driver.go`: Remove session ID generation and passing
- Existing Factoryfiles will need migration (breaking change)
