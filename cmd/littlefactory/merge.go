package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gbrindisi/littlefactory/internal/config"
	"github.com/gbrindisi/littlefactory/internal/tasks"
	"github.com/gbrindisi/littlefactory/internal/worktree"
	"github.com/spf13/cobra"
)

// Merge command flag variables
var (
	mergeChangeName string
	mergeForce      bool
	mergeMaxRetries int
)

// mergeCmd orchestrates verify-fix-rebase-merge-cleanup for a change.
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge a completed change into main",
	Long: `Merge a completed change from its worktree branch into main.

Runs a verify-fix loop, rebases onto main if needed, merges with --no-ff,
and cleans up the worktree and branch.

Examples:
  littlefactory merge -c feature-a              # Merge with default settings
  littlefactory merge -c feature-a --force       # Skip task completion check
  littlefactory merge -c feature-a --max-verify-retries 5`,
	Run: runMerge,
}

func init() {
	mergeCmd.Flags().StringVarP(&mergeChangeName, "change", "c", "",
		"Change name to merge (required)")
	_ = mergeCmd.MarkFlagRequired("change")

	mergeCmd.Flags().BoolVarP(&mergeForce, "force", "f", false,
		"Bypass task completion check")
	mergeCmd.Flags().IntVar(&mergeMaxRetries, "max-verify-retries", 3,
		"Maximum number of verify-fix loop retries")

	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) {
	// Find project root
	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting project root: %v\n", err)
		os.Exit(1)
	}

	// Check worktree exists
	exists, wtPath, err := worktree.WorktreeExists(projectRoot, mergeChangeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking worktree: %v\n", err)
		os.Exit(1)
	}
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: no worktree found for change %q\n", mergeChangeName)
		os.Exit(1)
	}

	// Check all tasks done (unless --force)
	if !mergeForce {
		if err := checkAllTasksDone(wtPath, mergeChangeName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Verify-fix loop
	if err := verifyFixLoop(mergeChangeName, mergeMaxRetries); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Rebase if main has advanced
	branchName := mergeChangeName
	if err := rebaseIfNeeded(projectRoot, wtPath, branchName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Merge into main
	if err := mergeIntoMain(projectRoot, branchName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully merged %s into main\n", branchName)

	// Cleanup worktree and branch (warn on failure, don't revert merge)
	cleanupWorktreeAndBranch(projectRoot, wtPath, branchName)
}

// checkAllTasksDone reads tasks.json from the worktree and verifies all tasks are done.
func checkAllTasksDone(wtPath, changeName string) error {
	tasksPath := filepath.Join(wtPath, ".littlefactory", "changes", changeName, "tasks.json")
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no tasks.json found for change %q", changeName)
		}
		return fmt.Errorf("reading tasks.json: %w", err)
	}

	var tf struct {
		Tasks []tasks.Task `json:"tasks"`
	}
	if err := json.Unmarshal(data, &tf); err != nil {
		return fmt.Errorf("parsing tasks.json: %w", err)
	}

	var incomplete []string
	for _, t := range tf.Tasks {
		if t.Status != "done" {
			incomplete = append(incomplete, fmt.Sprintf("  - %s (%s)", t.Title, t.Status))
		}
	}

	if len(incomplete) > 0 {
		return fmt.Errorf("incomplete tasks (use --force to bypass):\n%s", strings.Join(incomplete, "\n"))
	}

	return nil
}

// selfExe returns the path to the currently running littlefactory binary.
// This ensures merge shells out to the same binary, not a different installed version.
func selfExe() string {
	exe, err := os.Executable()
	if err != nil {
		return "littlefactory" // fallback to PATH
	}
	return exe
}

// verifyFixLoop runs the verify-fix cycle up to maxRetries times.
func verifyFixLoop(changeName string, maxRetries int) error {
	self := selfExe()

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Verify attempt %d/%d...\n", i+1, maxRetries)

		// Run verify
		verifyCmd := exec.Command(self, "verify", "-c", changeName)
		verifyCmd.Stdout = os.Stdout
		verifyCmd.Stderr = os.Stderr
		if err := verifyCmd.Run(); err == nil {
			fmt.Println("Verification passed")
			return nil
		}

		fmt.Println("Verification detected drift, running fix...")

		// If this is the last attempt, don't run fix
		if i == maxRetries-1 {
			break
		}

		// Run fix (littlefactory run in the worktree)
		runCmd := exec.Command(self, "run", "-c", changeName, "-w")
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Run(); err != nil {
			return fmt.Errorf("fix run failed: %w", err)
		}
	}

	return fmt.Errorf("verification failed after %d attempts", maxRetries)
}

// rebaseIfNeeded checks if main has advanced and rebases the branch if so.
func rebaseIfNeeded(projectRoot, wtPath, branchName string) error {
	// Check if main is ancestor of branch (i.e., branch is up to date)
	checkCmd := exec.Command("git", "merge-base", "--is-ancestor", "main", branchName)
	checkCmd.Dir = projectRoot
	if err := checkCmd.Run(); err == nil {
		// main is ancestor of branch -- no rebase needed
		return nil
	}

	fmt.Println("Main has advanced, rebasing...")

	// Rebase in the worktree
	rebaseCmd := exec.Command("git", "rebase", "main")
	rebaseCmd.Dir = wtPath
	rebaseCmd.Stdout = os.Stdout
	rebaseCmd.Stderr = os.Stderr
	if err := rebaseCmd.Run(); err != nil {
		// Abort rebase on conflict
		abortCmd := exec.Command("git", "rebase", "--abort")
		abortCmd.Dir = wtPath
		_ = abortCmd.Run()
		return fmt.Errorf("rebase failed (conflicts detected); resolve manually and retry")
	}

	fmt.Println("Rebase successful")
	return nil
}

// mergeIntoMain checks out main and merges the branch with --no-ff.
func mergeIntoMain(projectRoot, branchName string) error {
	// Checkout main
	checkoutCmd := exec.Command("git", "checkout", "main")
	checkoutCmd.Dir = projectRoot
	out, err := checkoutCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("checkout main: %s: %w", strings.TrimSpace(string(out)), err)
	}

	// Merge with --no-ff
	mergeGitCmd := exec.Command("git", "merge", "--no-ff", branchName, "-m",
		fmt.Sprintf("Merge branch '%s'", branchName))
	mergeGitCmd.Dir = projectRoot
	out, err = mergeGitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("merge failed: %s: %w", strings.TrimSpace(string(out)), err)
	}

	return nil
}

// cleanupWorktreeAndBranch removes the worktree and deletes the branch.
// Failures are warned about but do not cause the merge to be reverted.
func cleanupWorktreeAndBranch(projectRoot, wtPath, branchName string) {
	// Remove worktree (--force needed because agent may leave untracked files)
	removeCmd := exec.Command("git", "worktree", "remove", "--force", wtPath)
	removeCmd.Dir = projectRoot
	if out, err := removeCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to remove worktree: %s\n", strings.TrimSpace(string(out)))
	} else {
		fmt.Printf("Removed worktree at %s\n", wtPath)
	}

	// Delete branch
	deleteCmd := exec.Command("git", "branch", "-d", branchName)
	deleteCmd.Dir = projectRoot
	if out, err := deleteCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to delete branch: %s\n", strings.TrimSpace(string(out)))
	} else {
		fmt.Printf("Deleted branch %s\n", branchName)
	}
}
