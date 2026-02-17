## ADDED Requirements

### Requirement: Worktrees directory configuration
The system SHALL support a `worktrees_dir` option in Factoryfile.

#### Scenario: Parse worktrees_dir
- **WHEN** Factoryfile contains `worktrees_dir: ../worktrees`
- **THEN** System parses and stores the worktrees directory path

#### Scenario: Missing worktrees_dir uses default
- **WHEN** Factoryfile does not specify `worktrees_dir`
- **THEN** System uses `..` (sibling to repo) as default
