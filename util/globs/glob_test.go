package globs

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"testing"
)

func TestMatchAny(t *testing.T) {
	cases := []struct {
		Path typedpath.SlashPath
		Globs []Glob
		Expected bool
	} {
		{
			Path: "path/to/file",
			Globs : []Glob{},
			Expected: false,
		},
		{
			Path: "path/to/file",
			Globs : []Glob{
				"path",
			},
			Expected: true,
		},
		{
			Path: "path/to/file",
			Globs : []Glob{
				"path/*",
			},
			Expected: true,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%q, %v -> %t", c.Path, c.Globs, c.Expected), func(t *testing.T) {
			actual, _, err := MatchAny(c.Path, c.Globs)
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
