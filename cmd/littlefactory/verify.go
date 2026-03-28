package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gbrindisi/littlefactory/internal/agent"
	"github.com/gbrindisi/littlefactory/internal/config"
	"github.com/gbrindisi/littlefactory/internal/template"
	"github.com/gbrindisi/littlefactory/internal/worktree"
	"github.com/spf13/cobra"
)

// Verify command flag variables
var (
	verifyChangeName string
)

// verifyCmd runs agent-driven spec verification for a change.
var verifyCmd = &cobra.Command{
	Use:   "verify [agent]",
	Short: "Verify implementation against change specs",
	Long: `Run the verification agent to check implementation against change specs.
The agent is invoked with the VERIFIER.md template rendered with change context.

Exit codes:
  0  All specs satisfied (agent exited 0)
  1  Drift detected (agent exited non-zero)

Examples:
  littlefactory verify -c feature-a         # Verify with default agent
  littlefactory verify -c feature-a claude   # Verify with specific agent`,
	Args: cobra.MaximumNArgs(1),
	Run:  runVerify,
}

func init() {
	verifyCmd.Flags().StringVarP(&verifyChangeName, "change", "c", "",
		"Change name to verify (required)")
	_ = verifyCmd.MarkFlagRequired("change")

	rootCmd.AddCommand(verifyCmd)
}

// buildChangeContext scans the change directory and builds a ChangeContext.
func buildChangeContext(projectRoot, changeName string) *template.ChangeContext {
	changePath := filepath.Join(".littlefactory", "changes", changeName)
	changeDir := filepath.Join(projectRoot, changePath)

	ctx := &template.ChangeContext{
		ChangeName: changeName,
		ChangePath: changePath,
	}

	// Proposal
	proposalPath := filepath.Join(changePath, "proposal.md")
	if _, err := os.Stat(filepath.Join(projectRoot, proposalPath)); err == nil {
		ctx.ProposalPath = proposalPath
	}

	// Design
	designPath := filepath.Join(changePath, "design.md")
	if _, err := os.Stat(filepath.Join(projectRoot, designPath)); err == nil {
		ctx.DesignPath = designPath
	}

	// Tasks
	tasksPath := filepath.Join(changePath, "tasks.json")
	if _, err := os.Stat(filepath.Join(projectRoot, tasksPath)); err == nil {
		ctx.TasksPath = tasksPath
	}

	// Specs - collect all spec.md files under specs/
	specsDir := filepath.Join(changeDir, "specs")
	var specPaths []string
	entries, err := os.ReadDir(specsDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				specFile := filepath.Join(changePath, "specs", entry.Name(), "spec.md")
				if _, err := os.Stat(filepath.Join(projectRoot, specFile)); err == nil {
					specPaths = append(specPaths, specFile)
				}
			}
		}
	}
	if len(specPaths) > 0 {
		ctx.SpecsPaths = strings.Join(specPaths, "\n")
	}

	return ctx
}

func runVerify(cmd *cobra.Command, args []string) {
	// Find project root
	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting project root: %v\n", err)
		os.Exit(1)
	}

	// Load config
	cfg, err := config.LoadConfig(projectRoot, config.CLIFlags{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Validate change exists
	changeDir := filepath.Join(projectRoot, ".littlefactory", "changes", verifyChangeName)
	if _, err := os.Stat(changeDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: change %q not found at .littlefactory/changes/%s/\n", verifyChangeName, verifyChangeName)
		os.Exit(1)
	}

	// Resolve agent
	agentName := cfg.DefaultAgent
	if len(args) > 0 {
		agentName = args[0]
	}

	agentConfig, ok := cfg.Agents[agentName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: agent %q is not configured\n", agentName)
		os.Exit(1)
	}

	// Determine working directory: worktree if exists, else project root
	workDir := projectRoot
	exists, wtPath, err := worktree.WorktreeExists(projectRoot, verifyChangeName)
	if err == nil && exists {
		workDir = wtPath
		fmt.Printf("Running verification in worktree: %s\n", wtPath)
	}

	// Build change context
	changeCtx := buildChangeContext(workDir, verifyChangeName)

	// Load and render verifier template
	stateDir := filepath.Join(workDir, cfg.StateDir)
	tmpl, err := template.LoadVerifier(stateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading verifier template: %v\n", err)
		os.Exit(1)
	}

	prompt := template.RenderVerifier(tmpl, changeCtx)

	// Create agent
	ag := agent.NewConfigurableAgent(agentConfig.Command, agentConfig.Env)

	// Create context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Execute agent
	result, err := ag.Run(ctx, prompt, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running verification agent: %v\n", err)
		os.Exit(1)
	}

	// Map exit code: 0 -> 0 (pass), non-zero -> 1 (drift)
	if result.ExitCode != 0 {
		os.Exit(1)
	}
}
