package pathutil

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestSplitPathComponents(t *testing.T) {
	cases := []struct{
		path typedpath.SlashPath
		expected []typedpath.BaseName
	}{
		{ "", []typedpath.BaseName{} },
		{ "Foo", []typedpath.BaseName{"Foo"} },
		{ "Foo/", []typedpath.BaseName{"Foo"} },
		{ "Foo/Bar", []typedpath.BaseName{"Foo", "Bar"} },
		{ "Foo/Bar/Baz", []typedpath.BaseName{"Foo", "Bar", "Baz"} },
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%q", c.path), func(t *testing.T) {
			actual := SplitPathElements(c.path)

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf(cmp.Diff(c.expected, actual))
				return
			}
		})
	}
}
