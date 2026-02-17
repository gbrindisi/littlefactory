## 1. Worktree Package

- [x] 1.1 Create `internal/worktree/detect.go` with `GetCommonDir()` function using `git rev-parse --git-common-dir`
- [x] 1.2 Add `HasWorktrees()` function that checks if `<common-dir>/worktrees` exists and is non-empty
- [x] 1.3 Create `internal/worktree/list.go` with `List()` function that parses `git worktree list` output
- [x] 1.4 Create `internal/worktree/create.go` with `Create(branchName, worktreesDir)` function
- [x] 1.5 Add `IsClean()` function that checks `git status --porcelain` for uncommitted changes
- [x] 1.6 Add `WorktreeExists(branchName)` function to check if a worktree for a branch already exists
- [x] 1.7 Write tests for worktree package

## 2. Config Updates

- [x] 2.1 Add `WorktreesDir` field to config struct in `internal/config/config.go`
- [x] 2.2 Update config parsing to read `worktrees_dir` from Factoryfile
- [x] 2.3 Set default value for `WorktreesDir` to `..` (sibling to repo)
- [x] 2.4 Write tests for worktrees_dir config parsing

## 3. Run Command Changes

- [x] 3.1 Add `--change` / `-c` flag to run command
- [x] 3.2 Add `--worktree` / `-w` flag to run command
- [x] 3.3 Validate that `-w` requires `-c` flag
- [x] 3.4 Implement change validation (check `openspec/changes/<name>/tasks.json` exists)
- [x] 3.5 Implement worktree creation flow with clean-state and exists checks
- [x] 3.6 Pass change name and worktree path to driver
- [x] 3.7 Write tests for run command flags

## 4. Driver Updates

- [x] 4.1 Add `ChangeName` option to driver configuration
- [x] 4.2 Add `WorktreePath` option to driver configuration
- [x] 4.3 Update task source path resolution based on change name
- [x] 4.4 Implement workspace directory switching when worktree path is set
- [x] 4.5 Write tests for driver with change and worktree options

## 5. Status Command

- [x] 5.1 Create new `status` subcommand in `cmd/littlefactory/`
- [x] 5.2 Add `--change` / `-c` flag to status command
- [x] 5.3 Add `--all` flag to status command
- [x] 5.4 Add `--verbose` flag for detailed task list
- [x] 5.5 Implement status summary format: `<name>: X/Y done`
- [x] 5.6 Implement worktree discovery for `--all` flag using `git worktree list`
- [x] 5.7 Write tests for status command
