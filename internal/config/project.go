package config

import (
	"os"
	"path/filepath"
)

// FactoryfileMarker is the file name used to identify a project root
const FactoryfileMarker = "Factoryfile"

// FindProjectRoot locates the project root by searching for Factoryfile.
// It starts from the current working directory and walks up the directory tree.
// Returns an error if no Factoryfile is found.
// Returns an absolute path.
func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return findProjectRootFrom(cwd)
}

// findProjectRootFrom locates the project root starting from the given directory.
// Exported for testing purposes.
func findProjectRootFrom(startDir string) (string, error) {
	// Ensure we have an absolute path
	absDir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	dir := absDir
	for {
		// Check if Factoryfile exists in current directory
		factoryfilePath := filepath.Join(dir, FactoryfileMarker)
		if info, err := os.Stat(factoryfilePath); err == nil && !info.IsDir() {
			return dir, nil
		}

		// Get parent directory
		parent := filepath.Dir(dir)

		// Check if we've reached the filesystem root
		if parent == dir {
			// Factoryfile not found, return error
			return "", os.ErrNotExist
		}

		dir = parent
	}
}

// TasksDir returns the tasks directory path for the given project root.
// The tasks directory is always <project-root>/tasks/
// Returns an absolute path.
func TasksDir(projectRoot string) string {
	return filepath.Join(projectRoot, "tasks")
}

// EnsureTasksDir ensures the tasks directory exists, creating it if necessary.
// Uses mkdir -p behavior (creates parent directories as needed).
func EnsureTasksDir(projectRoot string) error {
	tasksPath := TasksDir(projectRoot)
	return os.MkdirAll(tasksPath, 0755)
}
