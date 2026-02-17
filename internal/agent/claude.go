// Package agent provides the configurable agent implementation.
package agent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/creack/pty"
	"github.com/yourusername/littlefactory/internal/config"
)

// ConfigurableAgent implements the Agent interface using a configurable command.
// The command is executed with the prompt passed via stdin.
type ConfigurableAgent struct {
	command   string
	envConfig map[string]config.EnvValue
}

// NewConfigurableAgent creates a new ConfigurableAgent with the given command string and environment config.
// The command string is parsed into args and executed when Run is called.
// The envConfig is used to set environment variables for the subprocess.
func NewConfigurableAgent(command string, envConfig map[string]config.EnvValue) *ConfigurableAgent {
	return &ConfigurableAgent{
		command:   command,
		envConfig: envConfig,
	}
}

// resolveEnv resolves the environment variables by evaluating shell commands and building an env slice.
// It starts with os.Environ() and applies overrides from envConfig.
// Shell commands are executed once, and their stdout (trimmed) is used as the value.
// Returns error if any shell command fails (non-zero exit).
func (c *ConfigurableAgent) resolveEnv() ([]string, error) {
	// Start with parent environment
	env := os.Environ()

	// Build a map for easy override
	envMap := make(map[string]string)
	for _, e := range env {
		if idx := strings.Index(e, "="); idx > 0 {
			envMap[e[:idx]] = e[idx+1:]
		}
	}

	// Apply overrides from config
	for key, envValue := range c.envConfig {
		var value string
		var err error

		if envValue.Static != "" {
			// Static value
			value = envValue.Static
		} else if envValue.Shell != "" {
			// Dynamic value - execute shell command
			value, err = c.evalShellCommand(envValue.Shell)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate shell command for %s: %w", key, err)
			}
		}

		envMap[key] = value
	}

	// Convert map back to slice
	result := make([]string, 0, len(envMap))
	for k, v := range envMap {
		result = append(result, k+"="+v)
	}

	return result, nil
}

// evalShellCommand executes a shell command and returns its stdout (trimmed).
// Returns error if the command exits with non-zero status.
func (c *ConfigurableAgent) evalShellCommand(cmd string) (string, error) {
	shellCmd := exec.Command("sh", "-c", cmd)
	output, err := shellCmd.Output()
	if err != nil {
		return "", err
	}

	// Trim trailing newlines (common in command substitution)
	result := strings.TrimRight(string(output), "\n")
	return result, nil
}

// Run executes the configured command with the given prompt via stdin.
// The prompt is passed via stdin, and output is streamed to the provided io.Writer.
// The subprocess runs in a PTY so isatty() returns true, preserving colors and spinners.
// Returns AgentResult with execution details, or error if execution fails.
func (c *ConfigurableAgent) Run(ctx context.Context, prompt string, output io.Writer) (AgentResult, error) {
	// Parse the command string into args
	args := parseCommand(c.command)
	if len(args) == 0 {
		return AgentResult{ExitCode: -1}, nil
	}

	// Build command with context
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// Resolve and set environment
	env, err := c.resolveEnv()
	if err != nil {
		return AgentResult{ExitCode: -1}, err
	}
	cmd.Env = env

	// Create PTY for stdout/stderr (so isatty() returns true)
	ptmx, tty, err := pty.Open()
	if err != nil {
		return AgentResult{ExitCode: -1}, fmt.Errorf("failed to create PTY: %w", err)
	}
	defer ptmx.Close() // Cleanup master
	defer tty.Close()   // Cleanup slave

	// Set stdout and stderr to PTY slave (for TTY detection)
	cmd.Stdout = tty
	cmd.Stderr = tty

	// Create a pipe for stdin so we can close it after writing
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return AgentResult{ExitCode: -1}, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return AgentResult{ExitCode: -1}, fmt.Errorf("failed to start command: %w", err)
	}

	// Close the TTY slave in parent process (child has its copy)
	tty.Close()

	// Write prompt to stdin and close it to signal EOF
	_, err = io.WriteString(stdinPipe, prompt)
	if err != nil {
		return AgentResult{ExitCode: -1}, fmt.Errorf("failed to write prompt: %w", err)
	}
	stdinPipe.Close() // Signal EOF to subprocess

	// Capture output for metrics while also streaming to output writer
	var rawBuf bytes.Buffer
	multiWriter := io.MultiWriter(output, &rawBuf)

	// Stream PTY output to writer and buffer
	// io.Copy will block until the PTY is closed (when process exits)
	_, copyErr := io.Copy(multiWriter, ptmx)

	// Wait for command to complete
	err = cmd.Wait()

	// Get the raw output (with ANSI codes)
	rawOutput := rawBuf.String()

	// Calculate raw bytes (includes ANSI codes)
	outputBytes := len(rawOutput)

	// Strip ANSI for line counting
	strippedOutput := stripansi.Strip(rawOutput)
	outputLines := 0
	if len(strippedOutput) > 0 {
		outputLines = strings.Count(strippedOutput, "\n")
		// Add 1 if output doesn't end with newline (count last line)
		if !strings.HasSuffix(strippedOutput, "\n") {
			outputLines++
		}
	}

	// Determine exit code
	exitCode := 0
	if err != nil {
		// Check for context timeout first
		if ctx.Err() == context.DeadlineExceeded {
			return AgentResult{
				ExitCode:    -1,
				Output:      rawOutput,
				OutputLines: outputLines,
				OutputBytes: outputBytes,
			}, ctx.Err()
		}

		// Check for context cancellation
		if ctx.Err() == context.Canceled {
			return AgentResult{
				ExitCode:    -1,
				Output:      rawOutput,
				OutputLines: outputLines,
				OutputBytes: outputBytes,
			}, ctx.Err()
		}

		// Extract exit code from ExitError
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// Non-exit error (e.g., command not found)
			return AgentResult{
				ExitCode:    -1,
				Output:      rawOutput,
				OutputLines: outputLines,
				OutputBytes: outputBytes,
			}, err
		}
	}

	// Check for copy error (but don't override process exit error)
	if copyErr != nil && err == nil {
		return AgentResult{
			ExitCode:    -1,
			Output:      rawOutput,
			OutputLines: outputLines,
			OutputBytes: outputBytes,
		}, fmt.Errorf("failed to copy PTY output: %w", copyErr)
	}

	return AgentResult{
		ExitCode:    exitCode,
		Output:      rawOutput,
		OutputLines: outputLines,
		OutputBytes: outputBytes,
	}, nil
}

// parseCommand splits a command string into args, handling quoted strings.
// Simple implementation that handles space-separated args and quoted strings.
func parseCommand(command string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range command {
		switch {
		case r == '"' || r == '\'':
			if inQuote && r == quoteChar {
				// End quote
				inQuote = false
				quoteChar = 0
			} else if !inQuote {
				// Start quote
				inQuote = true
				quoteChar = r
			} else {
				// Different quote char inside quotes, keep it
				current.WriteRune(r)
			}
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	// Add last arg if any
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
