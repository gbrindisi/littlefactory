// Package config provides configuration loading and management for littlefactory.
// It supports hierarchical configuration from defaults, Factoryfile, and CLI flags.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Default configuration values
const (
	DefaultMaxIterations = 10
	DefaultTimeout       = 600
	DefaultStateDir      = ".littlefactory"
	DefaultWorktreesDir  = ".."
)

// EnvValue represents an environment variable value that can be either
// a static string or a shell command to be evaluated.
type EnvValue struct {
	Static string // Static value (used when YAML value is a plain string)
	Shell  string // Shell command to evaluate (used when YAML value is {shell: "cmd"})
}

// UnmarshalYAML implements custom YAML unmarshaling for EnvValue.
// It supports two forms:
// - String: "literal value" -> Static field is set
// - Object: { shell: "command" } -> Shell field is set
func (e *EnvValue) UnmarshalYAML(node *yaml.Node) error {
	// Try to unmarshal as a plain string first
	var str string
	if err := node.Decode(&str); err == nil {
		e.Static = str
		return nil
	}

	// Try to unmarshal as an object with shell field
	var obj struct {
		Shell string `yaml:"shell"`
	}
	if err := node.Decode(&obj); err != nil {
		return fmt.Errorf("env value must be either a string or {shell: \"command\"}: %w", err)
	}

	if obj.Shell == "" {
		return errors.New("env value with shell form must have non-empty shell field")
	}

	e.Shell = obj.Shell
	return nil
}

// AgentConfig holds configuration for a named agent
type AgentConfig struct {
	Command string              `yaml:"command"`
	Env     map[string]EnvValue `yaml:"env,omitempty"`
}

// Config holds the littlefactory configuration
type Config struct {
	MaxIterations int                    `yaml:"max_iterations"`
	Timeout       int                    `yaml:"timeout"`
	StateDir      string                 `yaml:"state_dir"`
	WorktreesDir  string                 `yaml:"worktrees_dir"`
	SpecsDir      string                 `yaml:"specs_dir"`
	UseWorktree   bool                   `yaml:"use_worktree"`
	DefaultAgent  string                 `yaml:"default_agent"`
	Agents        map[string]AgentConfig `yaml:"agents"`
}

// CLIFlags holds CLI flag values for overriding config
type CLIFlags struct {
	MaxIterations *int
	Timeout       *int
}

// LoadConfig loads configuration with the following precedence:
// 1. Hardcoded defaults
// 2. Factoryfile (if exists) - overrides defaults
// 3. CLI flags (if provided) - overrides everything
func LoadConfig(projectRoot string, flags CLIFlags) (*Config, error) {
	// Start with defaults
	cfg := &Config{
		MaxIterations: DefaultMaxIterations,
		Timeout:       DefaultTimeout,
		StateDir:      DefaultStateDir,
		WorktreesDir:  DefaultWorktreesDir,
	}

	// Try to load Factoryfile
	if err := cfg.loadFactoryfile(projectRoot); err != nil {
		return nil, err
	}

	// Apply CLI flag overrides
	cfg.applyFlags(flags)

	// Validate final config
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// loadFactoryfile attempts to load configuration from Factoryfile at project root
// Returns nil if Factoryfile doesn't exist (not an error)
// Returns error if Factoryfile exists but is invalid
func (c *Config) loadFactoryfile(projectRoot string) error {
	// Check for "Factoryfile" first, then "Factoryfile.yaml"
	candidates := []string{
		filepath.Join(projectRoot, "Factoryfile"),
		filepath.Join(projectRoot, "Factoryfile.yaml"),
	}

	var factoryfilePath string
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			factoryfilePath = path
			break
		}
	}

	// No Factoryfile found - not an error, continue with defaults
	if factoryfilePath == "" {
		return nil
	}

	// Read and parse Factoryfile
	data, err := os.ReadFile(factoryfilePath)
	if err != nil {
		return fmt.Errorf("failed to read Factoryfile: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("invalid Factoryfile syntax: %w", err)
	}

	return nil
}

// applyFlags applies CLI flag overrides to the config
func (c *Config) applyFlags(flags CLIFlags) {
	if flags.MaxIterations != nil {
		c.MaxIterations = *flags.MaxIterations
	}
	if flags.Timeout != nil {
		c.Timeout = *flags.Timeout
	}
}

// validate checks that config values are valid
func (c *Config) validate() error {
	if c.MaxIterations <= 0 {
		return errors.New("max_iterations must be greater than 0")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}
	if c.StateDir == "" {
		return errors.New("state_dir must not be empty")
	}
	if len(c.Agents) == 0 {
		return errors.New("agents map must not be empty")
	}
	if c.DefaultAgent == "" {
		return errors.New("default_agent must be specified")
	}
	if _, exists := c.Agents[c.DefaultAgent]; !exists {
		return fmt.Errorf("default_agent %q does not exist in agents map", c.DefaultAgent)
	}
	return nil
}
