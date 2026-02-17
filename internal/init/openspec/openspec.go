package openspec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// defaultConfig is the minimal config written for new projects.
const defaultConfig = "schema: littlefactory\n"

// CheckInstalled verifies that the openspec binary is available in PATH.
func CheckInstalled() error {
	_, err := exec.LookPath("openspec")
	if err != nil {
		return fmt.Errorf("openspec is not installed or not in PATH: %w\nInstall it before running littlefactory init", err)
	}
	return nil
}

// SetupResult describes what Setup did with the config file.
type SetupResult struct {
	ConfigCreated bool // true if config.yaml was created; false if it already existed
}

// Setup extracts the embedded schema files and creates a default config if missing.
func Setup(projectRoot string) (*SetupResult, error) {
	if err := ExtractSchema(projectRoot); err != nil {
		return nil, fmt.Errorf("extracting openspec schema: %w", err)
	}

	configPath := filepath.Join(projectRoot, "openspec", "config.yaml")

	_, err := os.Stat(configPath)
	if err == nil {
		// Config already exists, preserve it.
		return &SetupResult{ConfigCreated: false}, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("checking openspec config: %w", err)
	}

	// Create the openspec directory if needed (ExtractSchema creates schemas/ but not the parent openspec/ itself in all cases).
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return nil, fmt.Errorf("creating openspec directory: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0o644); err != nil {
		return nil, fmt.Errorf("writing openspec config: %w", err)
	}

	return &SetupResult{ConfigCreated: true}, nil
}
