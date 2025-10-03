package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var (
	joinPathsCases = []struct {
		name     string
		base     string
		paths    []string
		expected string
	}{
		{
			"base_only",
			"~",
			nil,
			"~",
		},
		{
			"one_path",
			"~",
			[]string{"one"},
			"~/one",
		},
		{
			"two_path",
			"~",
			[]string{"one", "two"},
			"~/one/two",
		},
		{
			"empty_base",
			"",
			[]string{"one", "two"},
			"one/two",
		},
		{
			"empty_nil",
			"",
			nil,
			"",
		},
		{
			"invalid_base",
			"//../etc",
			nil,
			"/etc",
		},
		{
			"empty_path",
			"/",
			[]string{""},
			"/",
		},
		{
			"traverse_path",
			"/",
			[]string{"boot", "../usr//..///bin", "../etc"},
			"/etc",
		},
		{
			"invalid_path",
			"/home/",
			[]string{"..", "..", "/////..//./"},
			"/",
		},
		{
			"long_base_path_traverse",
			"/home/to/some/long/folder/for/testing/path/traverse",
			[]string{"../..", "..", ".."},
			"/home/to/some/long/folder",
		},
	}

	findGotCases = []struct {
		name          string
		path          string
		subfolders    []string
		subfiles      []string
		expectedValue *GotRepo
		expectedError error
	}{
		{
			"valid_gotrepo",
			"/repo/",
			[]string{
				"/repo/.got",
			},
			[]string{},
			&GotRepo{
				"/repo/",
				"/repo/.got",
				"",
			},
			nil,
		},
		{
			"valid_gotrepo_nested",
			"/repo/one/two/three",
			[]string{
				"/repo/.got",
				"/repo/one/two/three",
			},
			[]string{},
			&GotRepo{
				"/repo/",
				"/repo/.got",
				"",
			},
			nil,
		},
		{
			"not_gotrepo",
			"/notgot",
			[]string{
				"/notgot",
			},
			[]string{},
			nil,
			fmt.Errorf("fatal: not a got repository (or any of the parent directories): .got"),
		},
		{
			"not_gotrepo_nested",
			"/notgot/one/three",
			[]string{
				"/notgot/one/three",
			},
			[]string{},
			nil,
			fmt.Errorf("fatal: not a got repository (or any of the parent directories): .got"),
		},
		{
			"inside_gotdir",
			"/repo/.got",
			[]string{
				"/repo/.got",
			},
			[]string{},
			nil,
			fmt.Errorf("fatal: this operation must be run in a work tree"),
		},
		{
			"inside_gotdir_nested",
			"/repo/.got/one/two/three",
			[]string{
				"/repo/.got/one/two/three",
			},
			[]string{},
			nil,
			fmt.Errorf("fatal: this operation must be run in a work tree"),
		},
		{
			"invalid_gotfile_format",
			"/invalid",
			[]string{
				"/invalid",
			},
			[]string{
				"/invalid/.got",
			},
			nil,
			fmt.Errorf("fatal: invalid gotfile format: "),
		},
		{
			"invalid_gotfile_format_nested",
			"/invalid/one/two/three",
			[]string{
				"/invalid/one/two/three",
			},
			[]string{
				"/invalid/.got",
			},
			nil,
			fmt.Errorf("fatal: invalid gotfile format: "),
		},
	}
)

func FuzzJoinPaths(f *testing.F) {
	f.Add("base", "path1")
	f.Add("", "")
	f.Add("/", "../../etc/passwd")
	f.Add("", "./relative/path")
	f.Fuzz(func(t *testing.T, base string, path string) {
		joinPaths(base, path)
	})
}

func TestJoinPathsCases(t *testing.T) {
	functionName := runtime.FuncForPC(reflect.ValueOf(joinPaths).Pointer()).Name()
	for _, c := range joinPathsCases {
		t.Run(c.name, func(t *testing.T) {
			actual := joinPaths(c.base, c.paths...)
			if actual != c.expected {
				t.Errorf("%v(%v, %v) = %v; expected %v", functionName, c.base, c.paths, actual, c.expected)
			}
		})
	}
}

func TestFindGot(t *testing.T) {
	for _, c := range findGotCases {
		t.Run(c.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for _, dir := range c.subfolders {
				if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
					t.Errorf("error in %q: %v", c.name, err)
				}
			}
			for _, file := range c.subfiles {
				f, err := os.Create(filepath.Join(tmpDir, file))
				if err != nil {
					t.Errorf("error in %q: %v", c.name, err)
				}
				f.Close()
			}
			if c.expectedValue != nil {
				c.expectedValue.Worktree = filepath.Join(tmpDir, c.expectedValue.Worktree)
				c.expectedValue.Dir = filepath.Join(tmpDir, c.expectedValue.Dir)
			}

			t.Chdir(filepath.Join(tmpDir, c.path))
			actualGotRepo, actualError := findGot()
			if c.expectedError != nil {
				if actualError == nil || !strings.Contains(actualError.Error(), c.expectedError.Error()) {
					t.Errorf("actualError %v != expectedError %v", actualError, c.expectedError)
				}
			}
			if c.expectedValue != nil {
				if actualGotRepo == nil || actualGotRepo.Worktree != c.expectedValue.Worktree || actualGotRepo.Dir != c.expectedValue.Dir {
					t.Errorf("actualGotRepo %v != expectedGotRepo %v", actualError, c.expectedValue)
				}
			}
		})
	}
}
