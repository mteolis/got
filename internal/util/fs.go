package util

import "path/filepath"

func RepoPath(base string, path ...string) (string, error) {
	paths := append([]string{base}, path...)
	return filepath.Join(paths...), nil
}
