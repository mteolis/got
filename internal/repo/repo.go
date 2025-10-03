package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func findGot() (*GotRepo, error) {
	return findGotFrom(".")
}

func findGotFrom(path string) (*GotRepo, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	pathInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if pathInfo.IsDir() {
		gotName := ".got"
		if strings.Contains(path, gotName) {
			return nil, fmt.Errorf("fatal: this operation must be run in a work tree")
		} else {
			gotPath := joinPaths(path, gotName)
			gotInfo, err := os.Lstat(gotPath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					parent, err := filepath.Abs(joinPaths(path, ".."))
					if err != nil {
						return nil, err
					}

					if parent == path {
						return nil, fmt.Errorf("fatal: not a got repository (or any of the parent directories): .got")
					}

					return findGotFrom(parent)
				}
			}

			if gotInfo.IsDir() {
				return &GotRepo{
					Worktree: path,
					Dir:      gotPath,
				}, nil
			} else {
				return nil, fmt.Errorf("fatal: invalid gotfile format: %s", gotPath)
			}
		}
	}

	return findGotFrom(joinPaths(path, ".."))
}
