package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// validFactoryfileContent returns minimal valid Factoryfile content
func validFactoryfileContent() string {
	return `max_iterations: 10
timeout: 600
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Create temp directory with valid Factoryfile (now required)
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	if err := os.WriteFile(factoryfile, []byte(validFactoryfileContent()), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.MaxIterations != DefaultMaxIterations {
		t.Errorf("MaxIterations = %d, want %d", cfg.MaxIterations, DefaultMaxIterations)
	}
	if cfg.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %d, want %d", cfg.Timeout, DefaultTimeout)
	}
}

func TestLoadConfig_Factoryfile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Factoryfile with custom values
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 20
timeout: 1200
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.MaxIterations != 20 {
		t.Errorf("MaxIterations = %d, want 20", cfg.MaxIterations)
	}
	if cfg.Timeout != 1200 {
		t.Errorf("Timeout = %d, want 1200", cfg.Timeout)
	}
}

func TestLoadConfig_FactoryfileYaml(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Factoryfile.yaml (alternative filename)
	factoryfile := filepath.Join(tmpDir, "Factoryfile.yaml")
	content := `max_iterations: 15
timeout: 900
default_agent: aider
agents:
  aider:
    command: "aider --no-check"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile.yaml: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.MaxIterations != 15 {
		t.Errorf("MaxIterations = %d, want 15", cfg.MaxIterations)
	}
	if cfg.Timeout != 900 {
		t.Errorf("Timeout = %d, want 900", cfg.Timeout)
	}
}

func TestLoadConfig_FactoryfileTakesPrecedence(t *testing.T) {
	tmpDir := t.TempDir()

	// Create both Factoryfile and Factoryfile.yaml
	// Factoryfile (without .yaml) should take precedence
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	factoryfileYaml := filepath.Join(tmpDir, "Factoryfile.yaml")

	content1 := `max_iterations: 100
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	content2 := `max_iterations: 200
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content1), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}
	if err := os.WriteFile(factoryfileYaml, []byte(content2), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile.yaml: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Factoryfile (no extension) takes precedence
	if cfg.MaxIterations != 100 {
		t.Errorf("MaxIterations = %d, want 100 (Factoryfile should take precedence)", cfg.MaxIterations)
	}
}

func TestLoadConfig_FlagsOverrideFactoryfile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Factoryfile with values that will be overridden
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 20
timeout: 1200
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	// CLI flags override
	maxIter := 5
	timeout := 300
	flags := CLIFlags{
		MaxIterations: &maxIter,
		Timeout:       &timeout,
	}

	cfg, err := LoadConfig(tmpDir, flags)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Flags should win
	if cfg.MaxIterations != 5 {
		t.Errorf("MaxIterations = %d, want 5 (flags should override)", cfg.MaxIterations)
	}
	if cfg.Timeout != 300 {
		t.Errorf("Timeout = %d, want 300 (flags should override)", cfg.Timeout)
	}
}

func TestLoadConfig_FlagsOverrideDefaults(t *testing.T) {
	tmpDir := t.TempDir()

	// Factoryfile with agents config
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	if err := os.WriteFile(factoryfile, []byte(validFactoryfileContent()), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	// CLI flags provided
	maxIter := 3
	flags := CLIFlags{
		MaxIterations: &maxIter,
		// Timeout not set, should use default
	}

	cfg, err := LoadConfig(tmpDir, flags)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.MaxIterations != 3 {
		t.Errorf("MaxIterations = %d, want 3", cfg.MaxIterations)
	}
	if cfg.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %d, want %d (default)", cfg.Timeout, DefaultTimeout)
	}
}

func TestLoadConfig_PartialFactoryfile(t *testing.T) {
	tmpDir := t.TempDir()

	// Factoryfile with only max_iterations (but must have agents)
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 50
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.MaxIterations != 50 {
		t.Errorf("MaxIterations = %d, want 50", cfg.MaxIterations)
	}
	// Timeout should remain at default since not in Factoryfile
	if cfg.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %d, want %d (default)", cfg.Timeout, DefaultTimeout)
	}
}

func TestLoadConfig_InvalidFactoryfile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid YAML
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `invalid yaml: [
  not closed
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error for invalid YAML")
	}
}

func TestLoadConfig_InvalidMaxIterations(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Factoryfile with invalid max_iterations
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 0
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error for invalid max_iterations")
	}
}

func TestLoadConfig_InvalidTimeout(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Factoryfile with invalid timeout
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `timeout: -1
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error for invalid timeout")
	}
}

func TestLoadConfig_FlagCanSetInvalidValue(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid Factoryfile
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	if err := os.WriteFile(factoryfile, []byte(validFactoryfileContent()), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	// Flag with invalid value should fail validation
	maxIter := 0
	flags := CLIFlags{
		MaxIterations: &maxIter,
	}

	_, err := LoadConfig(tmpDir, flags)
	if err == nil {
		t.Error("LoadConfig() should return error for invalid flag value")
	}
}

func TestLoadConfig_PrecedenceOrder(t *testing.T) {
	tmpDir := t.TempDir()

	// Test full precedence: defaults < Factoryfile < flags
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 25
timeout: 800
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	// Only override max_iterations via flag, timeout should come from Factoryfile
	maxIter := 7
	flags := CLIFlags{
		MaxIterations: &maxIter,
	}

	cfg, err := LoadConfig(tmpDir, flags)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// max_iterations from flag
	if cfg.MaxIterations != 7 {
		t.Errorf("MaxIterations = %d, want 7 (from flag)", cfg.MaxIterations)
	}
	// timeout from Factoryfile
	if cfg.Timeout != 800 {
		t.Errorf("Timeout = %d, want 800 (from Factoryfile)", cfg.Timeout)
	}
}

// New validation tests for agents configuration

func TestLoadConfig_MultipleAgents(t *testing.T) {
	tmpDir := t.TempDir()

	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `default_agent: claude
agents:
  claude:
    command: "claude --print --dangerously-skip-permissions"
  aider:
    command: "aider --no-check-update"
  custom:
    command: "/path/to/custom-agent --flag"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if len(cfg.Agents) != 3 {
		t.Errorf("len(Agents) = %d, want 3", len(cfg.Agents))
	}
	if cfg.DefaultAgent != "claude" {
		t.Errorf("DefaultAgent = %q, want %q", cfg.DefaultAgent, "claude")
	}
	if cfg.Agents["claude"].Command != "claude --print --dangerously-skip-permissions" {
		t.Errorf("claude command = %q, want %q", cfg.Agents["claude"].Command, "claude --print --dangerously-skip-permissions")
	}
	if cfg.Agents["aider"].Command != "aider --no-check-update" {
		t.Errorf("aider command = %q, want %q", cfg.Agents["aider"].Command, "aider --no-check-update")
	}
}

func TestLoadConfig_EmptyAgentsMap(t *testing.T) {
	tmpDir := t.TempDir()

	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 10
default_agent: claude
agents: {}
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error for empty agents map")
	}
	if !strings.Contains(err.Error(), "agents map must not be empty") {
		t.Errorf("error = %q, want error containing 'agents map must not be empty'", err.Error())
	}
}

func TestLoadConfig_MissingDefaultAgent(t *testing.T) {
	tmpDir := t.TempDir()

	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 10
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error for missing default_agent")
	}
	if !strings.Contains(err.Error(), "default_agent must be specified") {
		t.Errorf("error = %q, want error containing 'default_agent must be specified'", err.Error())
	}
}

func TestLoadConfig_DefaultAgentNotInMap(t *testing.T) {
	tmpDir := t.TempDir()

	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 10
default_agent: nonexistent
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error for default_agent not in agents map")
	}
	if !strings.Contains(err.Error(), "does not exist in agents map") {
		t.Errorf("error = %q, want error containing 'does not exist in agents map'", err.Error())
	}
}

func TestLoadConfig_NoFactoryfile(t *testing.T) {
	// Create temp directory without Factoryfile
	tmpDir := t.TempDir()

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error when no Factoryfile exists (agents required)")
	}
	if !strings.Contains(err.Error(), "agents map must not be empty") {
		t.Errorf("error = %q, want error containing 'agents map must not be empty'", err.Error())
	}
}

// Table-driven tests for validation scenarios
func TestValidation(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			content: `max_iterations: 10
timeout: 600
default_agent: claude
agents:
  claude:
    command: "claude --print"
`,
			wantErr: false,
		},
		{
			name: "empty agents map",
			content: `max_iterations: 10
timeout: 600
default_agent: claude
agents: {}
`,
			wantErr:     true,
			errContains: "agents map must not be empty",
		},
		{
			name: "missing agents key",
			content: `max_iterations: 10
timeout: 600
default_agent: claude
`,
			wantErr:     true,
			errContains: "agents map must not be empty",
		},
		{
			name: "empty default_agent",
			content: `max_iterations: 10
timeout: 600
default_agent: ""
agents:
  claude:
    command: "claude --print"
`,
			wantErr:     true,
			errContains: "default_agent must be specified",
		},
		{
			name: "missing default_agent key",
			content: `max_iterations: 10
timeout: 600
agents:
  claude:
    command: "claude --print"
`,
			wantErr:     true,
			errContains: "default_agent must be specified",
		},
		{
			name: "default_agent not in map",
			content: `max_iterations: 10
timeout: 600
default_agent: aider
agents:
  claude:
    command: "claude --print"
`,
			wantErr:     true,
			errContains: "does not exist in agents map",
		},
		{
			name: "max_iterations zero",
			content: `max_iterations: 0
timeout: 600
default_agent: claude
agents:
  claude:
    command: "claude --print"
`,
			wantErr:     true,
			errContains: "max_iterations must be greater than 0",
		},
		{
			name: "timeout negative",
			content: `max_iterations: 10
timeout: -1
default_agent: claude
agents:
  claude:
    command: "claude --print"
`,
			wantErr:     true,
			errContains: "timeout must be greater than 0",
		},
		{
			name: "multiple agents with valid default",
			content: `max_iterations: 10
timeout: 600
default_agent: aider
agents:
  claude:
    command: "claude --print"
  aider:
    command: "aider --no-check"
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			factoryfile := filepath.Join(tmpDir, "Factoryfile")
			if err := os.WriteFile(factoryfile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write Factoryfile: %v", err)
			}

			_, err := LoadConfig(tmpDir, CLIFlags{})
			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadConfig() expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %q, want error containing %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("LoadConfig() unexpected error: %v", err)
				}
			}
		})
	}
}

// Tests for EnvValue unmarshaling

func TestEnvValue_UnmarshalYAML_StaticString(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `default_agent: myagent
agents:
  myagent:
    command: "agent-box run"
    env:
      API_KEY: "secret123"
      DB_HOST: "localhost"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	agent := cfg.Agents["myagent"]
	if len(agent.Env) != 2 {
		t.Errorf("len(Env) = %d, want 2", len(agent.Env))
	}

	apiKey := agent.Env["API_KEY"]
	if apiKey.Static != "secret123" {
		t.Errorf("API_KEY.Static = %q, want %q", apiKey.Static, "secret123")
	}
	if apiKey.Shell != "" {
		t.Errorf("API_KEY.Shell = %q, want empty", apiKey.Shell)
	}

	dbHost := agent.Env["DB_HOST"]
	if dbHost.Static != "localhost" {
		t.Errorf("DB_HOST.Static = %q, want %q", dbHost.Static, "localhost")
	}
	if dbHost.Shell != "" {
		t.Errorf("DB_HOST.Shell = %q, want empty", dbHost.Shell)
	}
}

func TestEnvValue_UnmarshalYAML_ShellCommand(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `default_agent: myagent
agents:
  myagent:
    command: "agent-box run"
    env:
      API_KEY:
        shell: "security find-generic-password -w -s 'Claude Code'"
      TIMESTAMP:
        shell: "date +%s"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	agent := cfg.Agents["myagent"]
	if len(agent.Env) != 2 {
		t.Errorf("len(Env) = %d, want 2", len(agent.Env))
	}

	apiKey := agent.Env["API_KEY"]
	if apiKey.Shell != "security find-generic-password -w -s 'Claude Code'" {
		t.Errorf("API_KEY.Shell = %q, want %q", apiKey.Shell, "security find-generic-password -w -s 'Claude Code'")
	}
	if apiKey.Static != "" {
		t.Errorf("API_KEY.Static = %q, want empty", apiKey.Static)
	}

	timestamp := agent.Env["TIMESTAMP"]
	if timestamp.Shell != "date +%s" {
		t.Errorf("TIMESTAMP.Shell = %q, want %q", timestamp.Shell, "date +%s")
	}
	if timestamp.Static != "" {
		t.Errorf("TIMESTAMP.Static = %q, want empty", timestamp.Static)
	}
}

func TestEnvValue_UnmarshalYAML_MixedStaticAndShell(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `default_agent: myagent
agents:
  myagent:
    command: "agent-box run"
    env:
      STATIC_VAR: "literal value"
      DYNAMIC_VAR:
        shell: "echo dynamic"
      ANOTHER_STATIC: "another literal"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	agent := cfg.Agents["myagent"]
	if len(agent.Env) != 3 {
		t.Errorf("len(Env) = %d, want 3", len(agent.Env))
	}

	// Check static
	staticVar := agent.Env["STATIC_VAR"]
	if staticVar.Static != "literal value" {
		t.Errorf("STATIC_VAR.Static = %q, want %q", staticVar.Static, "literal value")
	}
	if staticVar.Shell != "" {
		t.Errorf("STATIC_VAR.Shell should be empty")
	}

	// Check dynamic
	dynamicVar := agent.Env["DYNAMIC_VAR"]
	if dynamicVar.Shell != "echo dynamic" {
		t.Errorf("DYNAMIC_VAR.Shell = %q, want %q", dynamicVar.Shell, "echo dynamic")
	}
	if dynamicVar.Static != "" {
		t.Errorf("DYNAMIC_VAR.Static should be empty")
	}

	// Check another static
	anotherStatic := agent.Env["ANOTHER_STATIC"]
	if anotherStatic.Static != "another literal" {
		t.Errorf("ANOTHER_STATIC.Static = %q, want %q", anotherStatic.Static, "another literal")
	}
	if anotherStatic.Shell != "" {
		t.Errorf("ANOTHER_STATIC.Shell should be empty")
	}
}

func TestEnvValue_UnmarshalYAML_EmptyShellField(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `default_agent: myagent
agents:
  myagent:
    command: "agent-box run"
    env:
      BAD_VAR:
        shell: ""
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error for empty shell field")
	}
	if !strings.Contains(err.Error(), "non-empty shell field") {
		t.Errorf("error = %q, want error containing 'non-empty shell field'", err.Error())
	}
}

func TestEnvValue_UnmarshalYAML_InvalidObjectForm(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `default_agent: myagent
agents:
  myagent:
    command: "agent-box run"
    env:
      BAD_VAR:
        invalid_key: "value"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	_, err := LoadConfig(tmpDir, CLIFlags{})
	if err == nil {
		t.Error("LoadConfig() should return error for invalid env object form")
	}
	if !strings.Contains(err.Error(), "non-empty shell field") {
		t.Errorf("error = %q, want error containing 'non-empty shell field'", err.Error())
	}
}

func TestEnvValue_UnmarshalYAML_NoEnvField(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `default_agent: myagent
agents:
  myagent:
    command: "agent-box run"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	agent := cfg.Agents["myagent"]
	if agent.Env == nil {
		// nil is OK - optional field
		return
	}
	if len(agent.Env) != 0 {
		t.Errorf("len(Env) = %d, want 0 (or nil)", len(agent.Env))
	}
}

func TestEnvValue_UnmarshalYAML_EmptyEnvMap(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `default_agent: myagent
agents:
  myagent:
    command: "agent-box run"
    env: {}
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	agent := cfg.Agents["myagent"]
	if len(agent.Env) != 0 {
		t.Errorf("len(Env) = %d, want 0", len(agent.Env))
	}
}

// Tests for WorktreesDir config field

func TestLoadConfig_WorktreesDir_Default(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	if err := os.WriteFile(factoryfile, []byte(validFactoryfileContent()), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.WorktreesDir != DefaultWorktreesDir {
		t.Errorf("WorktreesDir = %q, want %q (default)", cfg.WorktreesDir, DefaultWorktreesDir)
	}
}

func TestLoadConfig_WorktreesDir_CustomRelative(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 10
timeout: 600
worktrees_dir: ../worktrees
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.WorktreesDir != "../worktrees" {
		t.Errorf("WorktreesDir = %q, want %q", cfg.WorktreesDir, "../worktrees")
	}
}

func TestLoadConfig_WorktreesDir_CustomAbsolute(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 10
timeout: 600
worktrees_dir: /tmp/my-worktrees
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.WorktreesDir != "/tmp/my-worktrees" {
		t.Errorf("WorktreesDir = %q, want %q", cfg.WorktreesDir, "/tmp/my-worktrees")
	}
}

func TestLoadConfig_WorktreesDir_NotOverriddenByOtherFields(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	// Factoryfile sets other fields but not worktrees_dir
	content := `max_iterations: 20
timeout: 1200
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// worktrees_dir should keep its default even when other fields are overridden
	if cfg.WorktreesDir != DefaultWorktreesDir {
		t.Errorf("WorktreesDir = %q, want %q (default)", cfg.WorktreesDir, DefaultWorktreesDir)
	}
	if cfg.MaxIterations != 20 {
		t.Errorf("MaxIterations = %d, want 20", cfg.MaxIterations)
	}
}

// Tests for SpecsDir config field

// Tests for UseWorktree config field

func TestLoadConfig_UseWorktree_DefaultFalse(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	if err := os.WriteFile(factoryfile, []byte(validFactoryfileContent()), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.UseWorktree {
		t.Error("UseWorktree should default to false")
	}
}

func TestLoadConfig_UseWorktree_Configured(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 10
timeout: 600
use_worktree: true
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if !cfg.UseWorktree {
		t.Error("UseWorktree should be true when configured")
	}
}

func TestLoadConfig_SpecsDir_DefaultEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	if err := os.WriteFile(factoryfile, []byte(validFactoryfileContent()), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.SpecsDir != "" {
		t.Errorf("SpecsDir = %q, want empty string (opt-in field)", cfg.SpecsDir)
	}
}

func TestLoadConfig_SpecsDir_Configured(t *testing.T) {
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	content := `max_iterations: 10
timeout: 600
specs_dir: "specs/"
default_agent: claude
agents:
  claude:
    command: "claude --print"
`
	if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Factoryfile: %v", err)
	}

	cfg, err := LoadConfig(tmpDir, CLIFlags{})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.SpecsDir != "specs/" {
		t.Errorf("SpecsDir = %q, want %q", cfg.SpecsDir, "specs/")
	}
}

// Table-driven tests for EnvValue edge cases
func TestEnvValue_UnmarshalYAML_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		envContent  string
		wantErr     bool
		errContains string
		validate    func(*testing.T, AgentConfig)
	}{
		{
			name: "numeric string value",
			envContent: `      PORT: "8080"
      COUNT: "42"`,
			wantErr: false,
			validate: func(t *testing.T, agent AgentConfig) {
				if agent.Env["PORT"].Static != "8080" {
					t.Errorf("PORT.Static = %q, want %q", agent.Env["PORT"].Static, "8080")
				}
				if agent.Env["COUNT"].Static != "42" {
					t.Errorf("COUNT.Static = %q, want %q", agent.Env["COUNT"].Static, "42")
				}
			},
		},
		{
			name:       "empty string value",
			envContent: `      EMPTY_VAR: ""`,
			wantErr:    false,
			validate: func(t *testing.T, agent AgentConfig) {
				if agent.Env["EMPTY_VAR"].Static != "" {
					t.Errorf("EMPTY_VAR.Static = %q, want empty string", agent.Env["EMPTY_VAR"].Static)
				}
			},
		},
		{
			name:       "multiline string value",
			envContent: `      MULTILINE: "line1\nline2\nline3"`,
			wantErr:    false,
			validate: func(t *testing.T, agent AgentConfig) {
				want := "line1\nline2\nline3"
				if agent.Env["MULTILINE"].Static != want {
					t.Errorf("MULTILINE.Static = %q, want %q", agent.Env["MULTILINE"].Static, want)
				}
			},
		},
		{
			name: "shell command with quotes",
			envContent: `      VAR:
        shell: "echo 'hello world'"`,
			wantErr: false,
			validate: func(t *testing.T, agent AgentConfig) {
				want := "echo 'hello world'"
				if agent.Env["VAR"].Shell != want {
					t.Errorf("VAR.Shell = %q, want %q", agent.Env["VAR"].Shell, want)
				}
			},
		},
		{
			name: "shell command with pipe",
			envContent: `      VAR:
        shell: "cat file.txt | grep pattern"`,
			wantErr: false,
			validate: func(t *testing.T, agent AgentConfig) {
				want := "cat file.txt | grep pattern"
				if agent.Env["VAR"].Shell != want {
					t.Errorf("VAR.Shell = %q, want %q", agent.Env["VAR"].Shell, want)
				}
			},
		},
		{
			name:       "special characters in static value",
			envContent: `      SPECIAL: "!@#$%^&*()"`,
			wantErr:    false,
			validate: func(t *testing.T, agent AgentConfig) {
				want := "!@#$%^&*()"
				if agent.Env["SPECIAL"].Static != want {
					t.Errorf("SPECIAL.Static = %q, want %q", agent.Env["SPECIAL"].Static, want)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			factoryfile := filepath.Join(tmpDir, "Factoryfile")
			content := `default_agent: myagent
agents:
  myagent:
    command: "agent-box run"
    env:
` + tt.envContent + "\n"
			if err := os.WriteFile(factoryfile, []byte(content), 0644); err != nil {
				t.Fatalf("failed to write Factoryfile: %v", err)
			}

			cfg, err := LoadConfig(tmpDir, CLIFlags{})
			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadConfig() expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %q, want error containing %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("LoadConfig() unexpected error: %v", err)
					return
				}
				if tt.validate != nil {
					tt.validate(t, cfg.Agents["myagent"])
				}
			}
		})
	}
}
