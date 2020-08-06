package pathutil

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestNewPathTree(t *testing.T) {
	cases := []struct {
		rootPaths []typedpath.SlashPath
		expected  *PathTree
	}{
		{
			[]typedpath.SlashPath{},
			&PathTree{false, PathTreeMap{}},
		},
		{
			[]typedpath.SlashPath{""},
			&PathTree{true, PathTreeMap{}},
		},
		{
			[]typedpath.SlashPath{"a"},
			&PathTree{
				false,
				PathTreeMap{
					"a": &PathTree{true, PathTreeMap{}},
				},
			},
		},
		{
			[]typedpath.SlashPath{"a/b"},
			&PathTree{
				false,
				PathTreeMap{
					"a": &PathTree{
						false,
						PathTreeMap{
							"b": &PathTree{true, PathTreeMap{}},
						},
					},
				},
			},
		},
		{
			[]typedpath.SlashPath{"a", "b"},
			&PathTree{
				false,
				PathTreeMap{
					"a": &PathTree{true, PathTreeMap{}},
					"b": &PathTree{true, PathTreeMap{}},
				},
			},
		},
		{
			[]typedpath.SlashPath{"a/1", "a/2"},
			&PathTree{
				false,
				PathTreeMap{
					"a": &PathTree{
						false,
						PathTreeMap{
							"1": &PathTree{true, PathTreeMap{}},
							"2": &PathTree{true, PathTreeMap{}},
						},
					},
				},
			},
		},
		{
			[]typedpath.SlashPath{"a", "a/b"},
			&PathTree{
				false,
				PathTreeMap{
					"a": &PathTree{
						true,
						PathTreeMap{
							"b": &PathTree{true, PathTreeMap{}},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.rootPaths), func(t *testing.T) {
			actual := NewPathTree(c.rootPaths...)

			if !reflect.DeepEqual(actual, c.expected) {
				t.Error(cmp.Diff(c.expected, actual))
				return
			}
		})
	}
}

func TestPathTree_Member(t *testing.T) {
	cases := []struct {
		paths    []typedpath.SlashPath
		path     []typedpath.BaseName
		expected bool
	}{
		{
			[]typedpath.SlashPath{},
			[]typedpath.BaseName{},
			false,
		},
		{
			[]typedpath.SlashPath{},
			[]typedpath.BaseName{"foo"},
			false,
		},
		{
			[]typedpath.SlashPath{"foo"},
			[]typedpath.BaseName{"foo"},
			false,
		},
		{
			[]typedpath.SlashPath{"foo"},
			[]typedpath.BaseName{"bar"},
			false,
		},
		{
			[]typedpath.SlashPath{"foo"},
			[]typedpath.BaseName{"foo", "bar"},
			true,
		},
		{
			[]typedpath.SlashPath{"foo/bar"},
			[]typedpath.BaseName{"foo"},
			false,
		},
		{
			[]typedpath.SlashPath{"foo/bar"},
			[]typedpath.BaseName{"foo", "bar"},
			false,
		},
		{
			[]typedpath.SlashPath{"foo/bar"},
			[]typedpath.BaseName{"foo", "baz"},
			false,
		},
		{
			[]typedpath.SlashPath{"foo/bar/baz"},
			[]typedpath.BaseName{"bar", "baz"},
			false,
		},
		{
			[]typedpath.SlashPath{"foo/bar/baz"},
			[]typedpath.BaseName{"foo", "baz"},
			false,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v %v -> %t", c.paths, c.path, c.expected), func(t *testing.T) {
			tree := NewPathTree(c.paths...)

			actual := tree.Member(c.path)
			if actual != c.expected {
				t.Errorf("want %t, got %t", c.expected, actual)
				return
			}
		})
	}
}
