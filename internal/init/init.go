// Package init orchestrates the littlefactory init and upgrade workflows,
// coordinating AGENTS.md setup, .gitignore management, and skill installation.
package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gbrindisi/littlefactory/internal/init/agentsmd"
	"github.com/gbrindisi/littlefactory/internal/init/gitignore"
	"github.com/gbrindisi/littlefactory/internal/init/skills"
)

// DefaultFactoryfile is the default content for a new Factoryfile.
const DefaultFactoryfile = `max_iterations: 10
timeout: 600
default_agent: claude

agents:
  claude:
    command: "claude --dangerously-skip-permissions --print"
    # Optional: configure environment variables for this agent
    # env:
    #   STATIC_VAR: "literal value"
    #   DYNAMIC_VAR:
    #     shell: "command to evaluate"
`

// logger provides formatted step and sub-operation logging for init output.
type logger struct {
	step  int
	total int
}

func newLogger(total int) *logger {
	return &logger{total: total}
}

// Step prints a numbered step header like "[1/4] Creating Factoryfile".
func (l *logger) Step(msg string) {
	l.step++
	fmt.Printf("[%d/%d] %s\n", l.step, l.total, msg)
}

// SubOp prints an indented sub-operation message under the current step.
func (l *logger) SubOp(msg string) {
	fmt.Printf("      %s\n", msg)
}

// Run executes the full init workflow for a new littlefactory project.
// It creates the Factoryfile, sets up AGENTS.md, updates .gitignore,
// and installs skills. Each step is logged with numbered progress output.
func Run(projectRoot string) error {
	log := newLogger(5)

	if err := createFactoryfile(log, projectRoot); err != nil {
		return err
	}

	if err := setupAgentsMD(log, projectRoot); err != nil {
		return err
	}

	if err := ensureGitignore(log, projectRoot); err != nil {
		return err
	}

	if err := installSkills(log, projectRoot); err != nil {
		return err
	}

	if err := setupChangesDir(log, projectRoot); err != nil {
		return err
	}

	fmt.Println("\nDone! littlefactory is ready.")
	return nil
}

func createFactoryfile(log *logger, projectRoot string) error {
	log.Step("Creating Factoryfile")

	factoryfilePath := filepath.Join(projectRoot, "Factoryfile")
	factoryfileYAMLPath := filepath.Join(projectRoot, "Factoryfile.yaml")

	if _, err := os.Stat(factoryfilePath); err == nil {
		return fmt.Errorf("factoryfile already exists")
	}
	if _, err := os.Stat(factoryfileYAMLPath); err == nil {
		return fmt.Errorf("factoryfile.yaml already exists")
	}

	if err := os.WriteFile(factoryfilePath, []byte(DefaultFactoryfile), 0o644); err != nil {
		return fmt.Errorf("creating Factoryfile: %w", err)
	}

	log.SubOp("Created Factoryfile with default configuration")
	return nil
}

func setupAgentsMD(log *logger, projectRoot string) error {
	log.Step("Setting up AGENTS.md")

	result, err := agentsmd.Setup(projectRoot)
	if err != nil {
		return fmt.Errorf("setting up AGENTS.md: %w", err)
	}

	switch result.Action {
	case agentsmd.ActionCreated:
		log.SubOp("Created AGENTS.md with default content")
		log.SubOp("Created symlink CLAUDE.md -> AGENTS.md")
	case agentsmd.ActionMigrated:
		log.SubOp("Renamed CLAUDE.md to AGENTS.md")
		log.SubOp("Created symlink CLAUDE.md -> AGENTS.md")
	case agentsmd.ActionMerged:
		log.SubOp("Merged CLAUDE.md into AGENTS.md")
		log.SubOp("Created symlink CLAUDE.md -> AGENTS.md")
	case agentsmd.ActionSkipped:
		log.SubOp("Already configured (CLAUDE.md is symlink to AGENTS.md)")
	}

	return nil
}

func ensureGitignore(log *logger, projectRoot string) error {
	log.Step("Updating .gitignore")

	result, err := gitignore.EnsureEntries(projectRoot)
	if err != nil {
		return fmt.Errorf("updating .gitignore: %w", err)
	}

	switch result.Action {
	case gitignore.ActionCreated:
		log.SubOp("Created .gitignore with littlefactory entries")
	case gitignore.ActionAdded:
		for _, entry := range result.Added {
			log.SubOp(fmt.Sprintf("Added %s", entry))
		}
	case gitignore.ActionSkipped:
		log.SubOp("All entries already present")
	}

	return nil
}

func installSkills(log *logger, projectRoot string) error {
	log.Step("Installing skills")

	if err := skills.ExtractSkills(projectRoot); err != nil {
		return fmt.Errorf("extracting skills: %w", err)
	}
	log.SubOp("Extracted embedded skills to .littlefactory/skills/")

	result, err := skills.CreateSymlinks(projectRoot)
	if err != nil {
		return fmt.Errorf("creating skill symlinks: %w", err)
	}

	if !result.ClaudeDirExists {
		log.SubOp("Skipped symlinks (.claude/ directory not found)")
		return nil
	}

	for _, entry := range result.Created() {
		log.SubOp(fmt.Sprintf("Linked .claude/skills/%s", entry.Name))
	}
	for _, entry := range result.Skipped() {
		log.SubOp(fmt.Sprintf("Skipped .claude/skills/%s (already exists)", entry.Name))
	}

	return nil
}

func setupChangesDir(log *logger, projectRoot string) error {
	log.Step("Creating changes directory")

	changesDir := filepath.Join(projectRoot, ".littlefactory", "changes")
	if err := os.MkdirAll(changesDir, 0o755); err != nil {
		return fmt.Errorf("creating changes directory: %w", err)
	}

	log.SubOp("Created .littlefactory/changes/")
	return nil
}
