package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gbrindisi/littlefactory/internal/agent"
	"github.com/gbrindisi/littlefactory/internal/config"
	"github.com/gbrindisi/littlefactory/internal/driver"
	lfinit "github.com/gbrindisi/littlefactory/internal/init"
	"github.com/gbrindisi/littlefactory/internal/tasks"
	"github.com/gbrindisi/littlefactory/internal/worktree"
	"github.com/spf13/cobra"
)

// Version information - set during build via ldflags
var (
	version = "dev"
	commit  = "unknown"
)

// CLI flag variables
var (
	maxIterations int
	timeout       int
	changeName    string
	tasksPath     string
	useWorktree   bool
)

// rootCmd is the base command for the CLI
var rootCmd = &cobra.Command{
	Use:   "littlefactory",
	Short: "An autonomous coding agent orchestrator",
	Long: `littlefactory is an autonomous coding agent that runs iterative loops
to complete software engineering tasks using Claude Code.`,
}

// initCmd initializes a new littlefactory project
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new littlefactory project",
	Long: `Initialize a new littlefactory project in the current directory.
Creates Factoryfile, sets up AGENTS.md, updates .gitignore, and installs skills.
Fails if Factoryfile or Factoryfile.yaml already exists.`,
	Run: runInit,
}

// runCmd is the command to run the agent loop
var runCmd = &cobra.Command{
	Use:   "run [agent]",
	Short: "Run the autonomous agent loop",
	Long: `Run the autonomous agent loop that iteratively processes tasks.
The agent will continue running until all tasks are complete,
max iterations is reached, or it is interrupted.

If agent name is not specified, uses the default_agent from config.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runRun,
}

// upgradeCmd upgrades an existing littlefactory project
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade littlefactory configuration",
	Long: `Upgrade an existing littlefactory project in the current directory.
Applies AGENTS.md setup, .gitignore updates, and skill installation.
Requires an existing Factoryfile (run 'littlefactory init' first for new projects).
All operations are idempotent and safe to run multiple times.`,
	Run: runUpgrade,
}

// versionCmd shows version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("littlefactory %s (commit: %s)\n", version, commit)
	},
}

func init() {
	rootCmd.Version = version

	// Add flags to run command
	runCmd.Flags().IntVar(&maxIterations, "max-iterations", 0,
		"Maximum number of iterations (default: from config or 10)")
	runCmd.Flags().IntVar(&timeout, "timeout", 0,
		"Timeout in seconds per iteration (default: from config or 600)")
	runCmd.Flags().StringVarP(&changeName, "change", "c", "",
		"Change name to use as task source")
	runCmd.Flags().StringVarP(&tasksPath, "tasks", "t", "",
		"Explicit path to tasks.json file")
	runCmd.Flags().BoolVarP(&useWorktree, "worktree", "w", false,
		"Create a new git worktree for the change")

	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// runInit creates a new littlefactory project with all setup steps.
func runInit(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := lfinit.Run(cwd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runUpgrade applies configuration improvements to an existing project.
func runUpgrade(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := lfinit.Upgrade(cwd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// validateChangeFlags validates --tasks, --change, and --worktree flag combinations.
// Returns an error if validation fails.
func validateChangeFlags(projectRoot, change, tasks string, wt bool) error {
	// -w requires -c
	if wt && change == "" {
		return fmt.Errorf("the --worktree flag requires --change to specify the branch name")
	}

	// Validate explicit --tasks path exists
	if tasks != "" {
		if _, err := os.Stat(tasks); os.IsNotExist(err) {
			return fmt.Errorf("tasks file not found: %s", tasks)
		}
		// --tasks takes priority; skip --change validation
		return nil
	}

	// Validate change exists if --change is specified
	if change != "" {
		changeDir := filepath.Join(projectRoot, ".littlefactory", "changes", change)
		if _, err := os.Stat(changeDir); os.IsNotExist(err) {
			return fmt.Errorf("change %q not found at .littlefactory/changes/%s/", change, change)
		}

		tasksPath := filepath.Join(changeDir, "tasks.json")
		if _, err := os.Stat(tasksPath); os.IsNotExist(err) {
			return fmt.Errorf("no tasks.json found for change %q", change)
		}
	}

	return nil
}

// prepareWorktree performs worktree precondition checks and creates the worktree.
// Returns the worktree path on success.
func prepareWorktree(projectRoot, change, worktreesDir string) (string, error) {
	// Check for clean working tree
	clean, err := worktree.IsClean(projectRoot)
	if err != nil {
		return "", fmt.Errorf("checking git status: %w", err)
	}
	if !clean {
		return "", fmt.Errorf("uncommitted changes detected; commit or stash before creating worktree")
	}

	// Check if worktree already exists
	exists, existingPath, err := worktree.WorktreeExists(projectRoot, change)
	if err != nil {
		return "", fmt.Errorf("checking worktree: %w", err)
	}
	if exists {
		return "", fmt.Errorf("worktree for %q already exists at %s; run without -w to use existing worktree", change, existingPath)
	}

	// Resolve worktrees directory
	resolvedDir := worktreesDir
	if !filepath.IsAbs(resolvedDir) {
		resolvedDir = filepath.Join(projectRoot, resolvedDir)
	}

	// Create the worktree
	return worktree.Create(projectRoot, change, resolvedDir)
}

// runRun executes the main agent loop
func runRun(cmd *cobra.Command, args []string) {
	// Detect project root
	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting project root: %v\n", err)
		os.Exit(1)
	}

	// Build CLI flags for config override
	cliFlags := config.CLIFlags{}
	if cmd.Flags().Changed("max-iterations") {
		cliFlags.MaxIterations = &maxIterations
	}
	if cmd.Flags().Changed("timeout") {
		cliFlags.Timeout = &timeout
	}

	// Load configuration
	cfg, err := config.LoadConfig(projectRoot, cliFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Resolve --tasks path: relative paths are resolved against cwd
	if tasksPath != "" && !filepath.IsAbs(tasksPath) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		tasksPath = filepath.Join(cwd, tasksPath)
	}

	// Validate tasks, change, and worktree flags
	if err := validateChangeFlags(projectRoot, changeName, tasksPath, useWorktree); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle worktree creation if -w is set
	var worktreePath string
	if useWorktree {
		wtPath, err := prepareWorktree(projectRoot, changeName, cfg.WorktreesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		worktreePath = wtPath
	}

	// Determine which agent to use
	agentName := cfg.DefaultAgent
	if len(args) > 0 {
		agentName = args[0]
	}

	// Look up agent config
	agentConfig, ok := cfg.Agents[agentName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: agent %q is not configured\n", agentName)
		fmt.Fprintf(os.Stderr, "Available agents: ")
		first := true
		for name := range cfg.Agents {
			if !first {
				fmt.Fprintf(os.Stderr, ", ")
			}
			fmt.Fprintf(os.Stderr, "%s", name)
			first = false
		}
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	// Create task source with priority: --tasks > --change > default
	var taskSource tasks.TaskSource
	if tasksPath != "" {
		// Use explicit --tasks path (highest priority)
		ts, err := tasks.NewJSONTaskSourceWithPath(tasksPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		taskSource = ts
	} else if changeName != "" {
		// Use change-specific tasks.json
		changeTasksPath := filepath.Join(projectRoot, ".littlefactory", "changes", changeName, "tasks.json")
		ts, err := tasks.NewJSONTaskSourceWithPath(changeTasksPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		taskSource = ts
	} else {
		// Default task source
		ts, err := tasks.NewJSONTaskSource(projectRoot, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		taskSource = ts
	}

	// Create agent from config
	ag := agent.NewConfigurableAgent(agentConfig.Command, agentConfig.Env)

	// Create driver
	d := driver.NewDriver(ag, taskSource, cfg, projectRoot)

	// Configure driver with change name and worktree path
	if changeName != "" {
		d.SetChangeName(changeName)
	}
	if worktreePath != "" {
		d.SetWorktreePath(worktreePath)
	}

	// Create context with signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handler to cancel context
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Run driver synchronously
	status := d.Run(ctx)

	// Map status to exit code
	exitCode := mapStatusToExitCode(status)
	os.Exit(exitCode)
}

// mapStatusToExitCode converts a RunStatus to an appropriate exit code.
// - 0: success (completed)
// - 130: interrupted (SIGINT/cancelled)
// - 1: error/failed
func mapStatusToExitCode(status driver.RunStatus) int {
	switch status {
	case driver.RunStatusCompleted:
		return 0
	case driver.RunStatusCancelled:
		return 130
	default:
		return 1
	}
}
