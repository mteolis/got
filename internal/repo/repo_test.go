package repo

import (
	"reflect"
	"runtime"
	"testing"
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
			actual := joinPaths(c.base, c.paths...)
			if actual != c.expected {
				functionName := runtime.FuncForPC(reflect.ValueOf(joinPaths).Pointer()).Name()
				t.Errorf("%v(%v, %v) = %v; expected %v", functionName, c.base, c.paths, actual, c.expected)
			}
		})
	}
}
