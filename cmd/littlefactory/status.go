package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gbrindisi/littlefactory/internal/config"
	"github.com/gbrindisi/littlefactory/internal/tasks"
	"github.com/gbrindisi/littlefactory/internal/worktree"
	"github.com/spf13/cobra"
)

// Status command flag variables
var (
	statusChangeName string
	statusAll        bool
	statusVerbose    bool
)

// statusCmd shows task progress for changes
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show task progress for changes",
	Long: `Show task progress for the current directory, a specific change, or all worktrees.

Examples:
  littlefactory status                  # Status for current directory
  littlefactory status -c feature-a     # Status for specific change
  littlefactory status --all            # Status for all worktrees
  littlefactory status -c feature-a -v  # Detailed task list`,
	Run: runStatus,
}

func init() {
	statusCmd.Flags().StringVarP(&statusChangeName, "change", "c", "",
		"OpenSpec change name to show status for")
	statusCmd.Flags().BoolVar(&statusAll, "all", false,
		"Show status for all worktrees")
	statusCmd.Flags().BoolVarP(&statusVerbose, "verbose", "v", false,
		"Show detailed task list")

	rootCmd.AddCommand(statusCmd)
}

// taskSummary holds counted task status for a change.
type taskSummary struct {
	Name       string
	Total      int
	Done       int
	InProgress string // title of in-progress task, empty if none
}

// formatSummary formats a taskSummary as a single line.
func formatSummary(s taskSummary) string {
	line := fmt.Sprintf("%s: %d/%d done", s.Name, s.Done, s.Total)
	if s.Done == s.Total && s.Total > 0 {
		line += " [complete]"
	} else if s.InProgress != "" {
		line += fmt.Sprintf(" (in_progress: %q)", s.InProgress)
	}
	return line
}

// readTasksFromPath reads and parses tasks from a tasks.json file path.
func readTasksFromPath(tasksPath string) ([]tasks.Task, error) {
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, err
	}
	var tf struct {
		Tasks []tasks.Task `json:"tasks"`
	}
	if err := json.Unmarshal(data, &tf); err != nil {
		return nil, fmt.Errorf("failed to parse tasks file: %w", err)
	}
	return tf.Tasks, nil
}

// summarizeTasks computes a taskSummary from a task list.
func summarizeTasks(name string, taskList []tasks.Task) taskSummary {
	s := taskSummary{
		Name:  name,
		Total: len(taskList),
	}
	for _, t := range taskList {
		switch t.Status {
		case "done":
			s.Done++
		case "in_progress":
			if s.InProgress == "" {
				s.InProgress = t.Title
			}
		}
	}
	return s
}

// printVerboseTasks prints detailed task list to stdout.
func printVerboseTasks(taskList []tasks.Task) {
	for _, t := range taskList {
		var indicator string
		switch t.Status {
		case "done":
			indicator = "[done]"
		case "in_progress":
			indicator = "[in_progress]"
		default:
			indicator = "[todo]"
		}
		fmt.Printf("  %s %s\n", indicator, t.Title)
	}
}

func runStatus(cmd *cobra.Command, args []string) {
	// Find project root
	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting project root: %v\n", err)
		os.Exit(1)
	}

	// Load config for state_dir
	cfg, err := config.LoadConfig(projectRoot, config.CLIFlags{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if statusAll {
		runStatusAll(projectRoot, cfg)
		return
	}

	if statusChangeName != "" {
		runStatusChange(projectRoot, statusChangeName)
		return
	}

	// Default: show status for current directory's tasks.json
	runStatusDefault(projectRoot, cfg)
}

// runStatusDefault shows status for the default tasks.json in the state directory.
func runStatusDefault(projectRoot string, cfg *config.Config) {
	tasksPath := filepath.Join(projectRoot, cfg.StateDir, "tasks.json")
	taskList, err := readTasksFromPath(tasksPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No tasks.json found in state directory")
			return
		}
		fmt.Fprintf(os.Stderr, "Error reading tasks: %v\n", err)
		os.Exit(1)
	}

	name := filepath.Base(projectRoot)
	s := summarizeTasks(name, taskList)
	fmt.Println(formatSummary(s))

	if statusVerbose {
		printVerboseTasks(taskList)
	}
}

// runStatusChange shows status for a specific change.
// If a worktree exists for the change, reads from the worktree's tasks.json.
func runStatusChange(projectRoot, change string) {
	// Check if a worktree exists for this change
	root := projectRoot
	exists, wtPath, err := worktree.WorktreeExists(projectRoot, change)
	if err == nil && exists {
		root = wtPath
	}

	tasksPath := filepath.Join(root, ".littlefactory", "changes", change, "tasks.json")
	taskList, err := readTasksFromPath(tasksPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: no tasks.json found for change '%s'\n", change)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error reading tasks: %v\n", err)
		os.Exit(1)
	}

	s := summarizeTasks(change, taskList)
	fmt.Println(formatSummary(s))

	if statusVerbose {
		printVerboseTasks(taskList)
	}
}

// runStatusAll shows status for all worktrees.
func runStatusAll(projectRoot string, cfg *config.Config) {
	worktrees, err := worktree.List(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing worktrees: %v\n", err)
		os.Exit(1)
	}

	found := false
	for _, wt := range worktrees {
		if wt.IsBare {
			continue
		}

		// Check for tasks.json in the worktree's state directory
		tasksPath := filepath.Join(wt.Path, cfg.StateDir, "tasks.json")
		taskList, err := readTasksFromPath(tasksPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Fprintf(os.Stderr, "Error reading tasks in %s: %v\n", wt.Path, err)
			continue
		}

		name := wt.BranchShort()
		if name == "" {
			name = filepath.Base(wt.Path)
		}

		s := summarizeTasks(name, taskList)
		fmt.Println(formatSummary(s))

		if statusVerbose {
			printVerboseTasks(taskList)
		}
		found = true
	}

	if !found {
		fmt.Println("No worktrees with tasks.json found")
	}
}
