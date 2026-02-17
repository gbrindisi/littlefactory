## Why

Agent commands often need environment variables (API keys, credentials) that must be set dynamically at runtime. Currently, users must wrap commands in `sh -c "..."` with complex quoting to inject env vars, which is error-prone and hard to read.

## What Changes

- Add `env` field to agent configuration in Factoryfile
- Support static env values: `VAR: "value"`
- Support dynamic env values via shell: `VAR: { shell: "command" }`
- Agent inherits parent environment, `env` field adds/overrides

## Capabilities

### New Capabilities
- `agent-env`: Environment variable configuration for agents, supporting static values and dynamic shell-evaluated values

### Modified Capabilities
- `config-management`: Add env field parsing to AgentConfig struct

## Impact

- `internal/config/config.go`: Add Env field to AgentConfig, custom YAML unmarshaling for EnvValue type
- `internal/agent/claude.go`: Accept env config, resolve shell values, set cmd.Env before execution
- `Factoryfile`: New optional `env` section under each agent
