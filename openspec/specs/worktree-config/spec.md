## ADDED Requirements

### Requirement: Worktrees directory configuration
The system SHALL support a `worktrees_dir` option in Factoryfile to specify where worktrees are created.

#### Scenario: Default worktrees directory
- **WHEN** Factoryfile does not specify `worktrees_dir`
- **THEN** System creates worktrees as siblings to the repository (git's default behavior)

#### Scenario: Custom worktrees directory
- **WHEN** Factoryfile specifies `worktrees_dir: ../worktrees`
- **THEN** System creates worktrees at `../worktrees/<change-name>`

#### Scenario: Relative worktrees directory
- **WHEN** Factoryfile specifies a relative `worktrees_dir` path
- **THEN** System resolves path relative to repository root

#### Scenario: Absolute worktrees directory
- **WHEN** Factoryfile specifies an absolute `worktrees_dir` path
- **THEN** System uses the absolute path directly
