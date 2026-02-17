package agent

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gbrindisi/littlefactory/internal/config"
)

func skipCI(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" {
		t.Skip("skipping PTY test in CI")
	}
}

// TestConfigurableAgentImplementsInterface verifies ConfigurableAgent implements Agent interface.
func TestConfigurableAgentImplementsInterface(t *testing.T) {
	var _ Agent = (*ConfigurableAgent)(nil)
}

// TestNewConfigurableAgent verifies constructor creates agent with command.
func TestNewConfigurableAgent(t *testing.T) {
	agent := NewConfigurableAgent("echo hello", nil)
	if agent == nil {
		t.Fatal("NewConfigurableAgent returned nil")
	}
	if agent.command != "echo hello" {
		t.Errorf("expected command 'echo hello', got %q", agent.command)
	}
}

// TestConfigurableAgentRun_Echo verifies running a simple echo command.
func TestConfigurableAgentRun_Echo(t *testing.T) {
	skipCI(t)
	ag := NewConfigurableAgent("cat", nil)
	result, err := ag.Run(context.Background(), "hello world", io.Discard)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
	if result.Output != "hello world" {
		t.Errorf("Output = %q, want %q", result.Output, "hello world")
	}
	if result.OutputBytes != 11 {
		t.Errorf("OutputBytes = %d, want 11", result.OutputBytes)
	}
	if result.OutputLines != 1 {
		t.Errorf("OutputLines = %d, want 1", result.OutputLines)
	}
}

// TestConfigurableAgentRun_WithArgs verifies command with arguments.
func TestConfigurableAgentRun_WithArgs(t *testing.T) {
	skipCI(t)
	ag := NewConfigurableAgent("head -n 1", nil)
	result, err := ag.Run(context.Background(), "line1\nline2\nline3", io.Discard)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
	// PTY converts LF to CRLF (terminal behavior)
	if result.Output != "line1\r\n" {
		t.Errorf("Output = %q, want %q", result.Output, "line1\r\n")
	}
}

// TestConfigurableAgentRun_NonZeroExit verifies non-zero exit code handling.
func TestConfigurableAgentRun_NonZeroExit(t *testing.T) {
	ag := NewConfigurableAgent("false", nil)
	result, err := ag.Run(context.Background(), "", io.Discard)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code")
	}
}

// TestConfigurableAgentRun_CommandNotFound verifies error for non-existent command.
func TestConfigurableAgentRun_CommandNotFound(t *testing.T) {
	ag := NewConfigurableAgent("nonexistent-command-12345", nil)
	_, err := ag.Run(context.Background(), "", io.Discard)

	if err == nil {
		t.Error("expected error for non-existent command")
	}
}

// TestConfigurableAgentRun_EmptyCommand verifies handling of empty command.
func TestConfigurableAgentRun_EmptyCommand(t *testing.T) {
	ag := NewConfigurableAgent("", nil)
	result, err := ag.Run(context.Background(), "", io.Discard)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != -1 {
		t.Errorf("ExitCode = %d, want -1 for empty command", result.ExitCode)
	}
}

// TestConfigurableAgentRun_ContextCancellation verifies context cancellation.
func TestConfigurableAgentRun_ContextCancellation(t *testing.T) {
	ag := NewConfigurableAgent("sleep 10", nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := ag.Run(ctx, "", io.Discard)

	// Context cancellation can happen at different points (start vs wait)
	if err == nil {
		t.Error("expected error for canceled context")
	}
	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context canceled error, got %v", err)
	}
}

// TestParseCommand verifies command string parsing.
func TestParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    []string
	}{
		{
			name:    "simple command",
			command: "echo hello",
			want:    []string{"echo", "hello"},
		},
		{
			name:    "command with multiple args",
			command: "git commit -m message",
			want:    []string{"git", "commit", "-m", "message"},
		},
		{
			name:    "double quoted arg",
			command: `echo "hello world"`,
			want:    []string{"echo", "hello world"},
		},
		{
			name:    "single quoted arg",
			command: `echo 'hello world'`,
			want:    []string{"echo", "hello world"},
		},
		{
			name:    "mixed quotes",
			command: `echo "hello" 'world'`,
			want:    []string{"echo", "hello", "world"},
		},
		{
			name:    "empty string",
			command: "",
			want:    []string{},
		},
		{
			name:    "only spaces",
			command: "   ",
			want:    []string{},
		},
		{
			name:    "quoted with spaces",
			command: `git commit -m "fix: some bug"`,
			want:    []string{"git", "commit", "-m", "fix: some bug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommand(tt.command)
			if len(got) != len(tt.want) {
				t.Errorf("parseCommand(%q) = %v, want %v", tt.command, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseCommand(%q)[%d] = %q, want %q", tt.command, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestResolveEnv_StaticValues verifies static environment variable resolution.
func TestResolveEnv_StaticValues(t *testing.T) {
	envConfig := map[string]config.EnvValue{
		"FOO": {Static: "bar"},
		"BAZ": {Static: "qux"},
	}

	ag := NewConfigurableAgent("echo", envConfig)
	env, err := ag.resolveEnv()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that our custom vars are present
	envMap := envSliceToMap(env)
	if envMap["FOO"] != "bar" {
		t.Errorf("FOO = %q, want %q", envMap["FOO"], "bar")
	}
	if envMap["BAZ"] != "qux" {
		t.Errorf("BAZ = %q, want %q", envMap["BAZ"], "qux")
	}

	// Check that parent env is inherited
	if _, exists := envMap["PATH"]; !exists {
		t.Error("expected PATH to be inherited from parent environment")
	}
}

// TestResolveEnv_ShellCommands verifies dynamic env var resolution via shell.
func TestResolveEnv_ShellCommands(t *testing.T) {
	envConfig := map[string]config.EnvValue{
		"COMPUTED":     {Shell: "printf hello"},
		"WITH_NEWLINE": {Shell: "echo world"}, // echo adds newline
	}

	ag := NewConfigurableAgent("echo", envConfig)
	env, err := ag.resolveEnv()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	envMap := envSliceToMap(env)
	if envMap["COMPUTED"] != "hello" {
		t.Errorf("COMPUTED = %q, want %q", envMap["COMPUTED"], "hello")
	}
	// Verify trailing newline is trimmed
	if envMap["WITH_NEWLINE"] != "world" {
		t.Errorf("WITH_NEWLINE = %q, want %q (newline should be trimmed)", envMap["WITH_NEWLINE"], "world")
	}
}

// TestResolveEnv_ShellCommandFailure verifies error handling for failed shell commands.
func TestResolveEnv_ShellCommandFailure(t *testing.T) {
	envConfig := map[string]config.EnvValue{
		"FAIL": {Shell: "exit 1"},
	}

	ag := NewConfigurableAgent("echo", envConfig)
	_, err := ag.resolveEnv()

	if err == nil {
		t.Error("expected error for failed shell command")
	}
	if !strings.Contains(err.Error(), "FAIL") {
		t.Errorf("error message should mention variable name FAIL, got: %v", err)
	}
}

// TestResolveEnv_MixedStaticAndShell verifies mixing static and shell env vars.
func TestResolveEnv_MixedStaticAndShell(t *testing.T) {
	envConfig := map[string]config.EnvValue{
		"STATIC_VAR": {Static: "static_value"},
		"SHELL_VAR":  {Shell: "printf dynamic_value"},
	}

	ag := NewConfigurableAgent("echo", envConfig)
	env, err := ag.resolveEnv()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	envMap := envSliceToMap(env)
	if envMap["STATIC_VAR"] != "static_value" {
		t.Errorf("STATIC_VAR = %q, want %q", envMap["STATIC_VAR"], "static_value")
	}
	if envMap["SHELL_VAR"] != "dynamic_value" {
		t.Errorf("SHELL_VAR = %q, want %q", envMap["SHELL_VAR"], "dynamic_value")
	}
}

// TestResolveEnv_OverrideParentEnv verifies config overrides parent environment.
func TestResolveEnv_OverrideParentEnv(t *testing.T) {
	t.Setenv("TEST_OVERRIDE_VAR", "original")

	envConfig := map[string]config.EnvValue{
		"TEST_OVERRIDE_VAR": {Static: "overridden"},
	}

	ag := NewConfigurableAgent("echo", envConfig)
	env, err := ag.resolveEnv()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	envMap := envSliceToMap(env)
	if envMap["TEST_OVERRIDE_VAR"] != "overridden" {
		t.Errorf("TEST_OVERRIDE_VAR = %q, want %q (should override parent)", envMap["TEST_OVERRIDE_VAR"], "overridden")
	}
}

// TestResolveEnv_EmptyConfig verifies nil/empty env config works.
func TestResolveEnv_EmptyConfig(t *testing.T) {
	ag := NewConfigurableAgent("echo", nil)
	env, err := ag.resolveEnv()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should just return parent env
	envMap := envSliceToMap(env)
	if _, exists := envMap["PATH"]; !exists {
		t.Error("expected PATH to be present when no env config provided")
	}
}

// TestConfigurableAgentRun_WithEnv verifies env vars are passed to subprocess.
func TestConfigurableAgentRun_WithEnv(t *testing.T) {
	skipCI(t)
	envConfig := map[string]config.EnvValue{
		"TEST_VAR": {Static: "test_value"},
	}

	ag := NewConfigurableAgent("sh -c 'echo $TEST_VAR'", envConfig)
	result, err := ag.Run(context.Background(), "", io.Discard)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
	if strings.TrimSpace(result.Output) != "test_value" {
		t.Errorf("Output = %q, want %q", result.Output, "test_value\n")
	}
}

// TestEvalShellCommand verifies shell command evaluation and newline trimming.
func TestEvalShellCommand(t *testing.T) {
	ag := NewConfigurableAgent("echo", nil)

	tests := []struct {
		name    string
		cmd     string
		want    string
		wantErr bool
	}{
		{
			name: "simple echo",
			cmd:  "echo hello",
			want: "hello",
		},
		{
			name: "printf without newline",
			cmd:  "printf world",
			want: "world",
		},
		{
			name: "multiline output",
			cmd:  "printf 'line1\\nline2\\nline3\\n'",
			want: "line1\nline2\nline3",
		},
		{
			name:    "command failure",
			cmd:     "exit 1",
			wantErr: true,
		},
		{
			name: "command with pipe",
			cmd:  "echo foo | tr a-z A-Z",
			want: "FOO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ag.evalShellCommand(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("evalShellCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("evalShellCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestConfigurableAgentRun_WithMixedEnv verifies both static and shell env vars work together.
func TestConfigurableAgentRun_WithMixedEnv(t *testing.T) {
	skipCI(t)
	envConfig := map[string]config.EnvValue{
		"STATIC_VAR": {Static: "static_value"},
		"SHELL_VAR":  {Shell: "printf shell_value"},
	}

	ag := NewConfigurableAgent("sh -c 'echo STATIC=$STATIC_VAR SHELL=$SHELL_VAR'", envConfig)
	result, err := ag.Run(context.Background(), "", io.Discard)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
	expected := "STATIC=static_value SHELL=shell_value"
	if strings.TrimSpace(result.Output) != expected {
		t.Errorf("Output = %q, want %q", strings.TrimSpace(result.Output), expected)
	}
}

// Helper: convert env slice to map for easy testing.
func envSliceToMap(env []string) map[string]string {
	m := make(map[string]string)
	for _, e := range env {
		if idx := strings.Index(e, "="); idx > 0 {
			m[e[:idx]] = e[idx+1:]
		}
	}
	return m
}
