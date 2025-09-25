package util

import (
	"testing"
)

func FuzzRepoPath(f *testing.F) {
	f.Add("base", "path1")
	f.Add("", "")
	f.Add("/", "../../etc/passwd")
	f.Add("", "./relative/path")
	f.Fuzz(func(t *testing.T, base string, path string) {
		_, err := RepoPath(base, path)
		if err != nil {
			t.Fatalf("err %v", err)
		}
	})
}

func TestRepoPathCases(t *testing.T) {
	cases := []struct {
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

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := RepoPath(c.base, c.paths...)
			if err != nil {
				t.Fatalf("err %v", err)
			}
			if actual != c.expected {
				t.Errorf("RepoPath(%v, %v) = %v; expected %v", c.base, c.paths, actual, c.expected)
			}
		})
	}
}
