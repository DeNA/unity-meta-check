package pathutil

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"sort"
	"testing"
)

func TestAllAncestorsAndSelf(t *testing.T) {
	cases := []struct{
		Path typedpath.SlashPath
		Expected []typedpath.SlashPath
	} {
		{
			Path: "",
			Expected: []typedpath.SlashPath{},
		},
		{
			Path: "path",
			Expected: []typedpath.SlashPath{
				"path",
			},
		},
		{
			Path: "path/",
			Expected: []typedpath.SlashPath{
				"path",
			},
		},
		{
			Path: "path/to/file",
			Expected: []typedpath.SlashPath{
				"path",
				"path/to",
				"path/to/file",
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%q -> %v", c.Path, c.Expected), func(t *testing.T) {
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
