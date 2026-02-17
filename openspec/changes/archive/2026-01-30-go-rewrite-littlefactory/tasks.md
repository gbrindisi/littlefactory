## 1. Project Setup

- [x] 1.1 Initialize Go module with `go mod init github.com/yourusername/littlefactory`
- [x] 1.2 Add cobra dependency (`github.com/spf13/cobra`)
- [x] 1.3 Add uuid dependency (`github.com/google/uuid`)
- [x] 1.4 Add yaml dependency (`gopkg.in/yaml.v3`)
- [x] 1.5 Create directory structure (cmd/, internal/, templates/)
- [x] 1.6 Copy CLAUDE.md template from ciccio to templates/CLAUDE.md

## 2. Core Types and Interfaces

- [x] 2.1 Define Task struct in internal/tasks/source.go
- [x] 2.2 Define TaskSource interface in internal/tasks/source.go
- [x] 2.3 Define Agent interface in internal/agent/agent.go
- [x] 2.4 Define AgentResult struct in internal/agent/agent.go
- [x] 2.5 Define RunMetadata struct in internal/driver/metadata.go
- [x] 2.6 Define IterationMetadata struct in internal/driver/metadata.go
- [x] 2.7 Define status enums (RunStatus, IterationStatus) in internal/driver/metadata.go

## 3. Configuration Management

- [x] 3.1 Create Config struct in internal/config/config.go
- [x] 3.2 Implement LoadConfig() with default values
- [x] 3.3 Implement Factoryfile loading from project root
- [x] 3.4 Implement CLI flag override logic
- [x] 3.5 Handle missing/invalid Factoryfile gracefully

## 4. Project Detection

- [x] 4.1 Implement FindProjectRoot() in internal/config/project.go
- [x] 4.2 Add logic to check current directory for .beads/
- [x] 4.3 Add logic to walk up parent directories searching for .beads/
- [x] 4.4 Add fallback to current directory if .beads/ not found
- [x] 4.5 Implement tasks directory path resolution (<project-root>/tasks/)

## 5. Beads Client Implementation

- [x] 5.1 Create BeadsClient struct in internal/tasks/beads.go
- [x] 5.2 Implement CheckBdCLI() to validate bd binary exists
- [x] 5.3 Implement Ready() using `bd ready --json`
- [x] 5.4 Implement Show(id) using `bd show <id> --json`
- [x] 5.5 Implement Close(id, reason) using `bd close <id> --reason <reason>`
- [x] 5.6 Implement Sync() using `bd sync`
- [x] 5.7 Add JSON parsing with proper error handling
- [x] 5.8 Handle bd show array response (extract first element)

## 6. Claude Code Agent Implementation

- [x] 6.1 Create ClaudeAgent struct in internal/agent/claude.go
- [x] 6.2 Implement Run(ctx, prompt, sessionID) method
- [x] 6.3 Build command: `claude --dangerously-skip-permissions --print --session-id <uuid>`
- [x] 6.4 Pass prompt via stdin to command
- [x] 6.5 Capture stdout + stderr combined output
- [x] 6.6 Handle context timeout (check ctx.Err())
- [x] 6.7 Compute session path using project root encoding
- [x] 6.8 Return AgentResult with exit code, output, line count, byte count
- [x] 6.9 Add sessionPath() helper function for path computation

## 7. Template System

- [x] 7.1 Create template.go in internal/template/
- [x] 7.2 Add go:embed directive for templates/CLAUDE.md
- [x] 7.3 Implement Load(projectRoot) with local override check
- [x] 7.4 Implement Render(tmpl, task) with string replacement
- [x] 7.5 Handle nil task (return template unchanged)
- [x] 7.6 Replace {task_id}, {task_title}, {task_description} placeholders

## 8. Metadata Tracking

- [x] 8.1 Implement RunMetadata JSON marshaling with ISO8601 timestamps
- [x] 8.2 Implement IterationMetadata JSON marshaling with ISO8601 timestamps
- [x] 8.3 Add custom MarshalJSON() for time.Time to ISO8601 format
- [x] 8.4 Implement SaveMetadata() to write tasks/run_metadata.json
- [x] 8.5 Add aggregate stats calculation (avg duration, success/fail counts)
- [x] 8.6 Implement generateRunID() with YYYYMMDD-HHMMSS format

## 9. Progress File Handling

- [x] 9.1 Create progress.go in internal/driver/
- [x] 9.2 Implement InitProgressFile() to create/initialize tasks/progress.txt
- [x] 9.3 Implement AppendSessionToProgress() with iteration format
- [x] 9.4 Handle missing session path gracefully (skip line if nil)
- [x] 9.5 Ensure append-only semantics (never truncate)

## 10. Driver Implementation

- [x] 10.1 Create Driver struct in internal/driver/driver.go
- [x] 10.2 Implement NewDriver() constructor with dependencies
- [x] 10.3 Implement Run(ctx) main loop method
- [x] 10.4 Add run initialization (metadata, progress file)
- [x] 10.5 Implement iteration loop (1 to maxIterations)
- [x] 10.6 Implement IsComplete() check (no ready tasks)
- [x] 10.7 Implement RunIteration(ctx, num) method
- [x] 10.8 Add task retrieval (GetNextTask, GetTaskDetails)
- [x] 10.9 Add template rendering with task injection
- [x] 10.10 Add agent execution with timeout context
- [x] 10.11 Add iteration metadata tracking (start, end, duration)
- [x] 10.12 Add status determination (completed/failed/timeout)
- [x] 10.13 Implement FinalizeRun() to compute final stats
- [x] 10.14 Add metadata save after each iteration

## 11. Visual Output

- [x] 11.1 Implement PrintStartBanner() with max iterations and ready count
- [x] 11.2 Implement PrintIterationBanner() with iteration number and target task
- [x] 11.3 Implement PrintSummary() with run stats and metadata path
- [x] 11.4 Format durations as human-readable (Xh Ym Zs)
- [x] 11.5 Print agent output to stdout during iteration

## 12. CLI with Cobra

- [x] 12.1 Create main.go in cmd/littlefactory/
- [x] 12.2 Initialize cobra root command
- [x] 12.3 Create "start" subcommand
- [x] 12.4 Add --max-iterations flag (default: 10)
- [x] 12.5 Add --timeout flag (default: 600)
- [x] 12.6 Wire up project detection, config loading, and driver
- [x] 12.7 Add context with signal handling for SIGINT
- [x] 12.8 Map RunStatus to exit codes (0, 130, 1)
- [x] 12.9 Add version command showing version info

## 13. Error Handling and Edge Cases

- [x] 13.1 Handle bd CLI not found with helpful error message
- [x] 13.2 Handle invalid Factoryfile with clear error
- [x] 13.3 Handle no ready tasks at start (exit early)
- [x] 13.4 Handle task source errors during iteration (continue)
- [x] 13.5 Handle agent execution errors during iteration (continue)
- [x] 13.6 Handle timeout expiration (mark as timeout, continue)
- [x] 13.7 Handle SIGINT gracefully (finalize metadata, exit 130)

## 14. Testing

- [x] 14.1 Add unit tests for template rendering (with/without task)
- [x] 14.2 Add unit tests for metadata JSON marshaling
- [x] 14.3 Add unit tests for project root detection
- [x] 14.4 Add unit tests for config loading (defaults, Factoryfile, flags)
- [x] 14.5 Create mock Agent implementation for driver tests
- [x] 14.6 Create mock TaskSource implementation for driver tests
- [x] 14.7 Add unit tests for driver iteration logic
- [x] 14.8 Add unit tests for session path computation

## 15. Documentation

- [x] 15.1 Create README.md with installation instructions
- [x] 15.2 Document Factoryfile format and options
- [x] 15.3 Document CLI flags and commands
- [x] 15.4 Add inline code comments for interfaces
- [x] 15.5 Document beads integration requirements

## 16. Final Verification

- [x] 16.1 Build binary: `go build -o littlefactory ./cmd/littlefactory`
- [x] 16.2 Test against ciccio project with real beads tasks
- [x] 16.3 Verify metadata JSON format matches Python output
- [x] 16.4 Verify progress.txt format matches Python output
- [x] 16.5 Verify visual output matches Python banners
- [x] 16.6 Test with no Factoryfile (defaults)
- [x] 16.7 Test with Factoryfile override
- [x] 16.8 Test with CLI flag override
- [x] 16.9 Test SIGINT handling (Ctrl+C during run)
- [x] 16.10 Test max iterations reached scenario
