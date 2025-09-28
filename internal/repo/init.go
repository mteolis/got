package repo

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

func InitRepo(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("fatal: could not %v\n", err)
		return nil
	}

	repo := &GotRepo{
		Worktree: absolutePath,
		Dir:      joinPaths(absolutePath, ".got/"),
	}

	if err := os.MkdirAll(repo.Worktree, 0755); err != nil {
		fmt.Printf("fatal: could not %v\n", err)
		return nil
	}

	entries, err := os.ReadDir(repo.Worktree)
	if err != nil {
		fmt.Printf("fatal: could not %v\n", err)
		return nil
	}

	if len(entries) > 0 {
		gotdir, err := os.Lstat(repo.Dir)
		if err != nil {
			if os.IsNotExist(err) {
				return repo.createGotRepo()
			}
			fmt.Printf("fatal: could not %v\n", err)
			return nil
		}
		if gotdir.IsDir() {
			gotentries, err := os.ReadDir(repo.Dir)
			if err != nil {
				fmt.Printf("fatal: could not %v\n", err)
				return nil
			} else if len(gotentries) == 0 {
				return repo.createGotRepo()
			}
			fmt.Printf("Reinitialized existing Got repository in %s\n", repo.Dir)
			return nil
		} else {
			fmt.Printf("fatal: could not %v\n", err)
			return nil
		}
	}

	return repo.createGotRepo()
}

func (repo GotRepo) createGotRepo() error {
	// reject init if dir is a symbolic link
	isSymlink, err := isSymLink(repo.Worktree)
	if err != nil {
		fmt.Printf("fatal: could not %v\n", err)
		return nil
	} else if isSymlink {
		fmt.Printf("fatal: could not create work tree dir %s: File is symlink\n", repo.Worktree)
		return nil
	}

	// mkdir .got/ and child dirs
	for _, dir := range []string{
		repo.Dir,
		joinPaths(repo.Dir, "branches"),
		joinPaths(repo.Dir, "objects"),
		joinPaths(repo.Dir, "refs"),
		joinPaths(repo.Dir, "refs", "heads"),
		joinPaths(repo.Dir, "refs", "tags"),
	} {
		if err := os.Mkdir(dir, 0755); err != nil {
			fmt.Printf("fatal: could not %v\n", err)
			return nil
		}
	}

	// .got/description
	desc := []byte("Unnamed repository; edit this file 'description' to name the repository.\n")
	if err := os.WriteFile(joinPaths(repo.Dir, "description"), desc, 0644); err != nil {
		fmt.Printf("fatal: could not %v\n", err)
		return nil
	}

	// .got/HEAD
	head := []byte("ref: refs/heads/gotmain\n")
	if err := os.WriteFile(joinPaths(repo.Dir, "HEAD"), head, 0644); err != nil {
		fmt.Printf("fatal: could not %v\n", err)
		return nil
	}

	// .got/config
	cfg := ini.Empty()
	cfg.Section("core").Key("\trepositoryformatversion").SetValue("0")
	cfg.Section("core").Key("\tfilemode").SetValue("false")
	cfg.Section("core").Key("\tbare").SetValue("false")
	cfg.Section("core").Key("\tlogallrefupdates").SetValue("false")

	if err := cfg.SaveTo(joinPaths(repo.Dir, "config")); err != nil {
		fmt.Printf("fatal: could not %v\n", err)
		return nil
	}

	fmt.Printf("Initialized empty Got repository in %s\n", repo.Dir)

	return nil
}
