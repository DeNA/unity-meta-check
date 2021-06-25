package globs

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"testing"
)

func TestMatchAny(t *testing.T) {
	cases := map[string]struct {
		Path     typedpath.SlashPath
		Globs    []Glob
		Cwd      typedpath.SlashPath
		Expected bool
	}{
		"empty globs (boundary case)": {
			Path:     "path/to/file",
			Cwd:      "/cwd",
			Globs:    []Glob{},
			Expected: false,
		},
		"easiest case": {
			Path: "path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				"path",
			},
			Expected: true,
		},
		"asterisk pattern": {
			Path: "path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				"path/*",
			},
			Expected: true,
		},
		"only dot (edge case)": {
			Path: "path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				".",
			},
			Expected: true,
		},
		"only asterisk (edge case)": {
			Path: "path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				"*",
			},
			Expected: true,
		},
		"dot asterisk (edge case)": {
			Path: "path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				"./*",
			},
			Expected: true,
		},
		"empty glob (edge case)": {
			Path: "path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				"",
			},
			Expected: true,
		},
		"only asterisk not match absolute path because the glob based on relative": {
			Path: "/path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				"*",
			},
			Expected: false,
		},
		"empty glob not match absolute path because the glob based on relative": {
			Path: "/path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				"",
			},
			Expected: false,
		},
		"only dot not match absolute path": {
			Path: "/path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				".",
			},
			Expected: false,
		},
		"slash asterisk match absolute path": {
			Path: "/path/to/file",
			Cwd:  "/cwd",
			Globs: []Glob{
				"/*",
			},
			Expected: true,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			actual, _, err := MatchAny(c.Path, c.Globs, c.Cwd)
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if actual != c.Expected {
				t.Errorf("want %t, got %t", c.Expected, actual)
				return
			}
		})
	}
}
