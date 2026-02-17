package openspec

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:embedded/schema
var embeddedSchema embed.FS

// schemaSubdir is the prefix inside embed.FS that contains the schema files.
const schemaSubdir = "embedded/schema"

// ExtractSchema copies all embedded schema files to <projectRoot>/openspec/schemas/littlefactory/,
// preserving the directory structure.
func ExtractSchema(projectRoot string) error {
	destRoot := filepath.Join(projectRoot, "openspec", "schemas", "littlefactory")

	schemaFS, err := fs.Sub(embeddedSchema, schemaSubdir)
	if err != nil {
		return err
	}

	return fs.WalkDir(schemaFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := filepath.Join(destRoot, path)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		data, err := fs.ReadFile(schemaFS, path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0o644)
	})
}
