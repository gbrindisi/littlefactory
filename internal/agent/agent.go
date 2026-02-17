// Package agent defines the agent interface and related types.
package agent

import (
	"context"
	"io"
)

// AgentResult contains the results from an agent execution.
type AgentResult struct {
	ExitCode    int    `json:"exit_code"`
	Output      string `json:"output"`
	OutputLines int    `json:"output_lines"`
	OutputBytes int    `json:"output_bytes"`
}

// Agent defines the interface for autonomous agent execution.
// Implementations can integrate with Claude Code, GPT, or other agents.
type Agent interface {
	// Run executes the agent with the given prompt.
	// The context is used for timeout enforcement.
	// The output parameter receives streaming output in real-time.
	// Returns AgentResult with execution details, or error if execution fails.
	Run(ctx context.Context, prompt string, output io.Writer) (AgentResult, error)
}
