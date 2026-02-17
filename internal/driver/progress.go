package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/littlefactory/internal/config"
)

const (
	// ProgressFileName is the name of the progress file
	ProgressFileName = "progress.md"
)

// InitProgressFile creates or opens the progress file at <state_dir>/progress.md.
// If the file doesn't exist, it creates it with a header.
// If it exists, it preserves the existing content (append-only semantics).
func InitProgressFile(projectRoot string, cfg *config.Config) error {
	stateDir := filepath.Join(projectRoot, cfg.StateDir)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	filePath := filepath.Join(stateDir, ProgressFileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create new file with header
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to create progress file: %w", err)
		}
		defer f.Close()

		header := fmt.Sprintf("# Little Factory Progress Log\n\n**Started:** %s\n\n---\n\n",
			time.Now().Format(time.RFC3339))
		if _, err := f.WriteString(header); err != nil {
			return fmt.Errorf("failed to write progress header: %w", err)
		}
	}
	// If file exists, do nothing - append-only semantics

	return nil
}

// AppendSessionToProgress appends iteration info to the progress file.
// Format: "## Iteration N", task ID, status, and separator with markdown bold labels.
func AppendSessionToProgress(projectRoot string, cfg *config.Config, iteration int, taskID string, status string) error {
	stateDir := filepath.Join(projectRoot, cfg.StateDir)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	filePath := filepath.Join(stateDir, ProgressFileName)

	// Open file for appending (create if not exists)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open progress file: %w", err)
	}
	defer f.Close()

	// Build the iteration block with markdown formatting
	var content string
	content += fmt.Sprintf("## Iteration %d\n\n", iteration)
	content += fmt.Sprintf("- **Task:** %s\n", taskID)
	content += fmt.Sprintf("- **Status:** %s\n\n", status)
	content += "---\n\n"

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to progress file: %w", err)
	}

	return nil
}

// ProgressFilePath returns the full path to the progress file.
func ProgressFilePath(projectRoot string, cfg *config.Config) string {
	return filepath.Join(projectRoot, cfg.StateDir, ProgressFileName)
}
