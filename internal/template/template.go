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

//go:embed embedded/VERIFIER.md
var embeddedVerifierTemplate string

// ChangeContext holds the context needed to render a verifier template.
type ChangeContext struct {
	ChangeName   string
	ChangePath   string
	ProposalPath string
	SpecsPaths   string
	DesignPath   string
	TasksPath    string
}

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

// LoadVerifier loads the VERIFIER.md template, checking for local override first.
// If a file exists at <stateDir>/agents/VERIFIER.md, it is used.
// Otherwise, the embedded template is returned.
func LoadVerifier(stateDir string) (string, error) {
	localPath := filepath.Join(stateDir, "agents", "VERIFIER.md")
	content, err := os.ReadFile(localPath)
	if err == nil {
		return string(content), nil
	}
	return embeddedVerifierTemplate, nil
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

// RenderVerifier replaces template placeholders with change context values.
// If ctx is nil, the template is returned unchanged.
func RenderVerifier(tmpl string, ctx *ChangeContext) string {
	if ctx == nil {
		return tmpl
	}
	result := strings.ReplaceAll(tmpl, "{change_name}", ctx.ChangeName)
	result = strings.ReplaceAll(result, "{change_path}", ctx.ChangePath)
	result = strings.ReplaceAll(result, "{proposal_path}", ctx.ProposalPath)
	result = strings.ReplaceAll(result, "{specs_paths}", ctx.SpecsPaths)
	result = strings.ReplaceAll(result, "{design_path}", ctx.DesignPath)
	result = strings.ReplaceAll(result, "{tasks_path}", ctx.TasksPath)
	return result
}
