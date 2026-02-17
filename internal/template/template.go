// Package template provides template loading and rendering for agent prompts.
package template

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/gbrindisi/littlefactory/internal/tasks"
)

//go:embed embedded/WORKER.md
var embeddedTemplate string

// Load loads the WORKER.md template, checking for local override first.
// If a file exists at <stateDir>/agents/WORKER.md, it is used.
// Otherwise, the embedded template is returned.
func Load(stateDir string) (string, error) {
	localPath := filepath.Join(stateDir, "agents", "WORKER.md")
	content, err := os.ReadFile(localPath)
	if err == nil {
		return string(content), nil
	}
	// Local override not found, use embedded template
	return embeddedTemplate, nil
}

// Render replaces template placeholders with task values.
// If task is nil, the template is returned unchanged.
func Render(tmpl string, task *tasks.Task) string {
	if task == nil {
		return tmpl
	}
	result := strings.ReplaceAll(tmpl, "{task_id}", task.ID)
	result = strings.ReplaceAll(result, "{task_title}", task.Title)
	result = strings.ReplaceAll(result, "{task_description}", task.Description)
	return result
}
