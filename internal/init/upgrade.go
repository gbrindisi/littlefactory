package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/littlefactory/internal/init/openspec"
)

// Upgrade applies littlefactory configuration improvements to an existing project.
// It requires a Factoryfile to exist and idempotently sets up AGENTS.md,
// .gitignore entries, skills, and OpenSpec configuration.
func Upgrade(projectRoot string) error {
	if err := openspec.CheckInstalled(); err != nil {
		return err
	}

	factoryfilePath := filepath.Join(projectRoot, "Factoryfile")
	factoryfileYAMLPath := filepath.Join(projectRoot, "Factoryfile.yaml")

	_, errPlain := os.Stat(factoryfilePath)
	_, errYAML := os.Stat(factoryfileYAMLPath)

	if errPlain != nil && errYAML != nil {
		return fmt.Errorf("No Factoryfile found. Run 'littlefactory init' first.")
	}

	log := newLogger(4)

	if err := setupAgentsMD(log, projectRoot); err != nil {
		return err
	}

	if err := ensureGitignore(log, projectRoot); err != nil {
		return err
	}

	if err := installSkills(log, projectRoot); err != nil {
		return err
	}

	if err := setupOpenSpec(log, projectRoot); err != nil {
		return err
	}

	fmt.Println("\nUpgrade complete.")
	return nil
}
