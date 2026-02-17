package worktree

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// initGitRepo creates a git repository in the given directory with an initial commit.
func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	run(t, dir, "git", "init")
	run(t, dir, "git", "config", "user.email", "test@test.com")
	run(t, dir, "git", "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("init"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(t, dir, "git", "add", ".")
	run(t, dir, "git", "commit", "-m", "initial commit")
}

// run executes a command in the given directory and fails the test on error.
func run(t *testing.T, dir string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command %q failed in %s: %v\noutput: %s", name+" "+strings.Join(args, " "), dir, err, out)
	}
	return strings.TrimSpace(string(out))
}

// --- GetCommonDir tests ---

func TestGetCommonDir_NormalRepo(t *testing.T) {
	tmpDir := t.TempDir()
	initGitRepo(t, tmpDir)

	commonDir, err := GetCommonDir(tmpDir)
	if err != nil {
		t.Fatalf("GetCommonDir failed: %v", err)
	}

	// For a normal repo, common dir is <repo>/.git
	expectedDir := filepath.Join(tmpDir, ".git")

	// Resolve symlinks on both sides (macOS /var -> /private/var).
	resolved, err := filepath.EvalSymlinks(commonDir)
	if err != nil {
		t.Fatalf("could not resolve common dir: %v", err)
	}
	expectedResolved, err := filepath.EvalSymlinks(expectedDir)
	if err != nil {
		t.Fatalf("could not resolve expected dir: %v", err)
	}

	if resolved != expectedResolved {
		t.Fatalf("expected common dir %q, got %q", expectedResolved, resolved)
	}
}

func TestGetCommonDir_FromWorktree(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, repoDir)

	// Create a worktree.
	wtDir := filepath.Join(tmpDir, "wt-branch")
	run(t, repoDir, "git", "worktree", "add", wtDir, "-b", "wt-branch")

	// GetCommonDir from within the worktree should point to the main repo's .git.
	commonDir, err := GetCommonDir(wtDir)
	if err != nil {
		t.Fatalf("GetCommonDir from worktree failed: %v", err)
	}

	expectedDir := filepath.Join(repoDir, ".git")
	resolved, err := filepath.EvalSymlinks(commonDir)
	if err != nil {
		t.Fatalf("could not resolve common dir: %v", err)
	}
	expectedResolved, err := filepath.EvalSymlinks(expectedDir)
	if err != nil {
		t.Fatalf("could not resolve expected dir: %v", err)
	}

	if resolved != expectedResolved {
		t.Fatalf("expected common dir %q, got %q", expectedResolved, resolved)
	}
}

func TestGetCommonDir_NotGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := GetCommonDir(tmpDir)
	if err == nil {
		t.Fatal("expected error for non-git directory")
	}
}

// --- HasWorktrees tests ---

func TestHasWorktrees_NoWorktrees(t *testing.T) {
	tmpDir := t.TempDir()
	initGitRepo(t, tmpDir)

	has, err := HasWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("HasWorktrees failed: %v", err)
	}
	if has {
		t.Fatal("expected no worktrees in fresh repo")
	}
}

func TestHasWorktrees_WithWorktree(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, repoDir)

	// Create a worktree.
	wtDir := filepath.Join(tmpDir, "wt-feature")
	run(t, repoDir, "git", "worktree", "add", wtDir, "-b", "feature")

	has, err := HasWorktrees(repoDir)
	if err != nil {
		t.Fatalf("HasWorktrees failed: %v", err)
	}
	if !has {
		t.Fatal("expected worktrees to be detected")
	}
}

// --- parseWorktreeList tests ---

func TestParseWorktreeList_SingleMainWorktree(t *testing.T) {
	output := "worktree /home/user/repo\nHEAD abc123\nbranch refs/heads/main\n\n"

	worktrees := parseWorktreeList(output)
	if len(worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(worktrees))
	}

	wt := worktrees[0]
	if wt.Path != "/home/user/repo" {
		t.Fatalf("expected path /home/user/repo, got %q", wt.Path)
	}
	if wt.Commit != "abc123" {
		t.Fatalf("expected commit abc123, got %q", wt.Commit)
	}
	if wt.Branch != "refs/heads/main" {
		t.Fatalf("expected branch refs/heads/main, got %q", wt.Branch)
	}
	if wt.BranchShort() != "main" {
		t.Fatalf("expected short branch main, got %q", wt.BranchShort())
	}
	if wt.IsDetached {
		t.Fatal("expected not detached")
	}
	if wt.IsBare {
		t.Fatal("expected not bare")
	}
}

func TestParseWorktreeList_MultipleWorktrees(t *testing.T) {
	output := `worktree /home/user/repo
HEAD abc123
branch refs/heads/main

worktree /home/user/wt-feature
HEAD def456
branch refs/heads/feature

`

	worktrees := parseWorktreeList(output)
	if len(worktrees) != 2 {
		t.Fatalf("expected 2 worktrees, got %d", len(worktrees))
	}

	if worktrees[0].BranchShort() != "main" {
		t.Fatalf("expected first worktree branch main, got %q", worktrees[0].BranchShort())
	}
	if worktrees[1].BranchShort() != "feature" {
		t.Fatalf("expected second worktree branch feature, got %q", worktrees[1].BranchShort())
	}
}

func TestParseWorktreeList_DetachedHead(t *testing.T) {
	output := "worktree /home/user/repo\nHEAD abc123\ndetached\n\n"

	worktrees := parseWorktreeList(output)
	if len(worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(worktrees))
	}

	if !worktrees[0].IsDetached {
		t.Fatal("expected detached HEAD")
	}
	if worktrees[0].Branch != "" {
		t.Fatalf("expected empty branch for detached, got %q", worktrees[0].Branch)
	}
}

func TestParseWorktreeList_BareRepo(t *testing.T) {
	output := "worktree /home/user/repo.git\nbare\n\n"

	worktrees := parseWorktreeList(output)
	if len(worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(worktrees))
	}

	if !worktrees[0].IsBare {
		t.Fatal("expected bare repo entry")
	}
}

func TestParseWorktreeList_EmptyOutput(t *testing.T) {
	worktrees := parseWorktreeList("")
	if len(worktrees) != 0 {
		t.Fatalf("expected 0 worktrees, got %d", len(worktrees))
	}
}

func TestParseWorktreeList_NoTrailingNewline(t *testing.T) {
	output := "worktree /home/user/repo\nHEAD abc123\nbranch refs/heads/main"

	worktrees := parseWorktreeList(output)
	if len(worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(worktrees))
	}
	if worktrees[0].BranchShort() != "main" {
		t.Fatalf("expected branch main, got %q", worktrees[0].BranchShort())
	}
}

// --- List integration test ---

func TestList_ReturnsMainWorktree(t *testing.T) {
	tmpDir := t.TempDir()
	initGitRepo(t, tmpDir)

	worktrees, err := List(tmpDir)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(worktrees) < 1 {
		t.Fatal("expected at least 1 worktree (main)")
	}

	// The first entry should be the main working tree.
	resolved, err := filepath.EvalSymlinks(worktrees[0].Path)
	if err != nil {
		t.Fatalf("could not resolve worktree path: %v", err)
	}
	expectedResolved, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("could not resolve tmpDir: %v", err)
	}

	if resolved != expectedResolved {
		t.Fatalf("expected first worktree path %q, got %q", expectedResolved, resolved)
	}
}

func TestList_IncludesCreatedWorktrees(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, repoDir)

	// Create a worktree.
	wtDir := filepath.Join(tmpDir, "wt-feature")
	run(t, repoDir, "git", "worktree", "add", wtDir, "-b", "feature")

	worktrees, err := List(repoDir)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(worktrees) != 2 {
		t.Fatalf("expected 2 worktrees, got %d", len(worktrees))
	}

	// Find the feature worktree.
	var found bool
	for _, wt := range worktrees {
		if wt.BranchShort() == "feature" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to find worktree for branch 'feature'")
	}
}

// --- BranchShort tests ---

func TestBranchShort_FullRef(t *testing.T) {
	wt := Worktree{Branch: "refs/heads/feature-x"}
	if got := wt.BranchShort(); got != "feature-x" {
		t.Fatalf("expected feature-x, got %q", got)
	}
}

func TestBranchShort_Empty(t *testing.T) {
	wt := Worktree{}
	if got := wt.BranchShort(); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

// --- IsClean tests ---

func TestIsClean_CleanRepo(t *testing.T) {
	tmpDir := t.TempDir()
	initGitRepo(t, tmpDir)

	clean, err := IsClean(tmpDir)
	if err != nil {
		t.Fatalf("IsClean failed: %v", err)
	}
	if !clean {
		t.Fatal("expected clean repo after initial commit")
	}
}

func TestIsClean_UnstagedChanges(t *testing.T) {
	tmpDir := t.TempDir()
	initGitRepo(t, tmpDir)

	// Create an untracked file.
	if err := os.WriteFile(filepath.Join(tmpDir, "dirty.txt"), []byte("dirty"), 0o644); err != nil {
		t.Fatal(err)
	}

	clean, err := IsClean(tmpDir)
	if err != nil {
		t.Fatalf("IsClean failed: %v", err)
	}
	if clean {
		t.Fatal("expected dirty repo with untracked file")
	}
}

func TestIsClean_StagedChanges(t *testing.T) {
	tmpDir := t.TempDir()
	initGitRepo(t, tmpDir)

	// Modify a tracked file and stage it.
	if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("modified"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(t, tmpDir, "git", "add", "README.md")

	clean, err := IsClean(tmpDir)
	if err != nil {
		t.Fatalf("IsClean failed: %v", err)
	}
	if clean {
		t.Fatal("expected dirty repo with staged changes")
	}
}

// --- WorktreeExists tests ---

func TestWorktreeExists_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	initGitRepo(t, tmpDir)

	exists, path, err := WorktreeExists(tmpDir, "feature-x")
	if err != nil {
		t.Fatalf("WorktreeExists failed: %v", err)
	}
	if exists {
		t.Fatalf("expected no worktree for feature-x, found at %q", path)
	}
}

func TestWorktreeExists_Found(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, repoDir)

	// Create a worktree.
	wtDir := filepath.Join(tmpDir, "wt-feature-a")
	run(t, repoDir, "git", "worktree", "add", wtDir, "-b", "feature-a")

	exists, path, err := WorktreeExists(repoDir, "feature-a")
	if err != nil {
		t.Fatalf("WorktreeExists failed: %v", err)
	}
	if !exists {
		t.Fatal("expected worktree to exist for feature-a")
	}

	// Resolve symlinks for comparison.
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("could not resolve returned path: %v", err)
	}
	expectedResolved, err := filepath.EvalSymlinks(wtDir)
	if err != nil {
		t.Fatalf("could not resolve expected path: %v", err)
	}

	if resolvedPath != expectedResolved {
		t.Fatalf("expected worktree path %q, got %q", expectedResolved, resolvedPath)
	}
}

// --- Create tests ---

func TestCreate_CreatesWorktree(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, repoDir)

	wtPath, err := Create(repoDir, "new-feature", tmpDir)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "new-feature")
	if wtPath != expectedPath {
		t.Fatalf("expected worktree path %q, got %q", expectedPath, wtPath)
	}

	// Verify the worktree directory exists.
	info, err := os.Stat(wtPath)
	if err != nil {
		t.Fatalf("worktree directory does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected worktree path to be a directory")
	}

	// Verify the branch was created.
	branch := run(t, wtPath, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if branch != "new-feature" {
		t.Fatalf("expected branch new-feature, got %q", branch)
	}
}

func TestCreate_DuplicateBranch(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, repoDir)

	// Create first worktree.
	_, err := Create(repoDir, "dup-branch", tmpDir)
	if err != nil {
		t.Fatalf("first Create failed: %v", err)
	}

	// Second create with same branch should fail.
	_, err = Create(repoDir, "dup-branch", filepath.Join(tmpDir, "other"))
	if err == nil {
		t.Fatal("expected error when creating worktree with existing branch name")
	}
}

func TestCreate_WorktreeIsListable(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, repoDir)

	_, err := Create(repoDir, "listed-feature", tmpDir)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify it appears in List.
	worktrees, err := List(repoDir)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	var found bool
	for _, wt := range worktrees {
		if wt.BranchShort() == "listed-feature" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected created worktree to appear in List output")
	}

	// Verify HasWorktrees returns true.
	has, err := HasWorktrees(repoDir)
	if err != nil {
		t.Fatalf("HasWorktrees failed: %v", err)
	}
	if !has {
		t.Fatal("expected HasWorktrees to return true after creating worktree")
	}

	// Verify WorktreeExists finds it.
	exists, _, err := WorktreeExists(repoDir, "listed-feature")
	if err != nil {
		t.Fatalf("WorktreeExists failed: %v", err)
	}
	if !exists {
		t.Fatal("expected WorktreeExists to find created worktree")
	}
}
