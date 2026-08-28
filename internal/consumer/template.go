package consumer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Files is the complete framework-owned surface of a newly created consumer
// project. Keep this intentionally tiny: application code, language choices,
// architecture, tests, and product documentation must come from the user's
// project, not from the framework source repository.
var Files = map[string]string{
	".apf/project.yaml": "version: 1\n",
	".gitignore":       ".apf/cache/\n",
	"AGENTS.md": `# Project Development

This repository uses AI Project Framework.

Work on the user's project normally. Keep changes local, modular, and easy to verify.
Use APF automatically when available; users should not need to operate the framework themselves.
Do not copy framework implementation files or framework development documentation into this repository.
`,
}

// NewProject creates an empty consumer project. The target may not already
// contain files; APF must never overwrite or mix itself into an existing tree
// through the `new` command.
func NewProject(target string) error {
	if target == "" {
		return errors.New("project directory is required")
	}

	abs, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	if err := requireEmptyOrMissing(abs); err != nil {
		return err
	}

	paths := make([]string, 0, len(Files))
	for path := range Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, rel := range paths {
		path := filepath.Join(abs, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(Files[rel]), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func requireEmptyOrMissing(path string) error {
	entries, err := os.ReadDir(path)
	if err == nil {
		if len(entries) != 0 {
			return fmt.Errorf("target directory is not empty: %s", path)
		}
		return nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.MkdirAll(path, 0o755)
}
