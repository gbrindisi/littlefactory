## 1. Config Structure Changes

- [x] 1.1 Update Config struct: replace AgentConfig with AgentsConfig map and add DefaultAgent field
- [x] 1.2 Create AgentConfig struct with Command field
- [x] 1.3 Update config validation: require agents map non-empty, default_agent must exist in map
- [x] 1.4 Update config loading to parse new structure
- [x] 1.5 Add tests for new config validation rules

## 2. Agent Interface Changes

- [x] 2.1 Remove sessionID parameter from Agent interface Run method
- [x] 2.2 Create ConfigurableAgent that takes command string and executes it
- [x] 2.3 Remove ClaudeAgent (replaced by ConfigurableAgent)
- [x] 2.4 Remove session path computation logic
- [x] 2.5 Update agent tests for new interface

## 3. Driver Changes

- [x] 3.1 Remove session ID generation from driver
- [x] 3.2 Update driver to call Agent.Run without sessionID
- [x] 3.3 Update driver tests

## 4. CLI Changes

- [x] 4.1 Add init command that creates Factoryfile with defaults
- [x] 4.2 Init command fails if Factoryfile or Factoryfile.yaml exists
- [x] 4.3 Rename start command to run
- [x] 4.4 Add optional agent name positional argument to run command
- [x] 4.5 Run command uses default_agent when no agent specified
- [x] 4.6 Run command fails with clear error for unknown agent name
- [x] 4.7 Update main.go to create agent from config instead of hardcoding

## 5. Integration

- [x] 5.1 Verify end-to-end: init creates valid Factoryfile
- [x] 5.2 Verify end-to-end: run with default agent works
- [x] 5.3 Verify end-to-end: run with explicit agent name works
- [x] 5.4 Update any existing tests that depend on old behavior
