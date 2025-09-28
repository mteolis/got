package repo

import (
	"os"
	"path/filepath"
)

type GotRepo struct {
	Worktree string
	Dir      string
	Config   string
}

func joinPaths(base string, path ...string) string {
	paths := append([]string{base}, path...)
	return filepath.Join(paths...)
}

func isSymLink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}
