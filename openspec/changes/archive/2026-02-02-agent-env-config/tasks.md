## 1. Config Layer

- [x] 1.1 Add EnvValue type to config.go with custom YAML unmarshaling (string or {shell: string})
- [x] 1.2 Add Env field to AgentConfig struct as map[string]EnvValue
- [x] 1.3 Add unit tests for EnvValue unmarshaling (static and shell variants)

## 2. Agent Layer

- [x] 2.1 Update ConfigurableAgent to accept env config in constructor
- [x] 2.2 Implement resolveEnv() to evaluate shell commands and build env slice
- [x] 2.3 Set cmd.Env to os.Environ() merged with resolved env before execution
- [x] 2.4 Add unit tests for env resolution (static, shell, errors)

## 3. Integration

- [x] 3.1 Update main.go to pass AgentConfig.Env to NewConfigurableAgent
- [x] 3.2 Update Factoryfile example in defaultFactoryfile const (optional env comment)
- [x] 3.3 Manual test with real Factoryfile using shell env var
