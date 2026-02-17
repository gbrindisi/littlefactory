package skills

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:embedded/skills
var embeddedSkills embed.FS

// skillsSubdir is the prefix inside embed.FS that contains the skill directories.
const skillsSubdir = "embedded/skills"

// ExtractSkills copies all embedded skills to <projectRoot>/.littlefactory/skills/,
// preserving the directory structure. Each skill is a subdirectory containing its files.
func ExtractSkills(projectRoot string) error {
	destRoot := filepath.Join(projectRoot, ".littlefactory", "skills")

	skillsFS, err := fs.Sub(embeddedSkills, skillsSubdir)
	if err != nil {
		return err
	}

	return fs.WalkDir(skillsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip dotfiles (e.g. .gitkeep used to keep the embed directory non-empty).
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		destPath := filepath.Join(destRoot, path)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		data, err := fs.ReadFile(skillsFS, path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0o644)
	})
}
