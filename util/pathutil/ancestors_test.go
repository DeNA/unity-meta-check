package pathutil

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"sort"
	"testing"
)

func TestAllAncestorsAndSelf(t *testing.T) {
	cases := map[string]struct{
		Path typedpath.SlashPath
		Expected []typedpath.SlashPath
	} {
		"empty (boundary case)": {
			Path: "",
			Expected: []typedpath.SlashPath{},
		},
		"only dot (boundary case)": {
			Path: ".",
			Expected: []typedpath.SlashPath{
			},
		},
		"single path (easiest case)": {
			Path: "path",
			Expected: []typedpath.SlashPath{
				"path",
			},
		},
		"trailing slash (edge case)": {
			Path: "path/",
			Expected: []typedpath.SlashPath{
				"path",
			},
		},
		"relative path (normal case)": {
			Path: "path/to/file",
			Expected: []typedpath.SlashPath{
				"path",
				"path/to",
				"path/to/file",
			},
		},
		"dot-prefix path (edge case)": {
			Path: "./path/to/file",
			Expected: []typedpath.SlashPath{
				"path",
				"path/to",
				"path/to/file",
			},
		},
		"dot-dot-prefix path (edge case)": {
			Path: "../path/to/file",
			Expected: []typedpath.SlashPath{
				"..",
				"../path",
				"../path/to",
				"../path/to/file",
			},
		},
		"absolute path (edge case)": {
			Path: "/path/to/file",
			Expected: []typedpath.SlashPath{
				"/",
				"/path",
				"/path/to",
				"/path/to/file",
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			actual := AllAncestorsAndSelf(c.Path)

			sort.Slice(actual, func(i, j int) bool {
				return actual[i] < actual[j]
			})

			if !reflect.DeepEqual(actual, c.Expected) {
				t.Error(cmp.Diff(c.Expected, actual))
			}
		})
	}
}
