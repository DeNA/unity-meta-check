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
		rootPaths []PathPair[string]
		expected  PathTree[string]
	}{
		{
			[]PathPair[string]{},
			PathTree[string]{},
		},
		{
			[]PathPair[string]{
				{"", ""},
			},
			PathTree[string]{},
		},
		{
			[]PathPair[string]{
				{"a", "A"},
			},
			PathTree[string]{
				"a": &PathTreeEntry[string]{
					Value:   ptrString("A"),
					Subtree: PathTree[string]{},
				},
			},
		},
		{
			[]PathPair[string]{
				{"a/b", "a/b"},
			},
			PathTree[string]{
				"a": &PathTreeEntry[string]{
					nil,
					PathTree[string]{
						"b": &PathTreeEntry[string]{ptrString("a/b"), PathTree[string]{}},
					},
				},
			},
		},
		{
			[]PathPair[string]{
				{"a", "A"},
				{"b", "B"},
			},
			PathTree[string]{
				"a": &PathTreeEntry[string]{ptrString("A"), PathTree[string]{}},
				"b": &PathTreeEntry[string]{ptrString("B"), PathTree[string]{}},
			},
		},
		{
			[]PathPair[string]{
				{"a/1", "A/1"},
				{"a/2", "A/2"},
			},
			PathTree[string]{
				"a": &PathTreeEntry[string]{
					nil,
					PathTree[string]{
						"1": &PathTreeEntry[string]{ptrString("A/1"), PathTree[string]{}},
						"2": &PathTreeEntry[string]{ptrString("A/2"), PathTree[string]{}},
					},
				},
			},
		},
		{
			[]PathPair[string]{
				{"a", "A"},
				{"a/b", "A/B"},
			},
			PathTree[string]{
				"a": &PathTreeEntry[string]{
					ptrString("A"),
					PathTree[string]{
						"b": &PathTreeEntry[string]{ptrString("A/B"), PathTree[string]{}},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.rootPaths), func(t *testing.T) {
			actual := NewPathTreeWithValues[string](c.rootPaths...)

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

func TestPathTree_Postorder(t *testing.T) {
	tree := NewPathTree(
		"a/b",
		"a/b/d",
		"a/c",
	)

	trace := make([]typedpath.SlashPath, 0)
	_ = tree.Postorder(func(b typedpath.SlashPath, p PathTreeEntry[struct{}]) error {
		trace = append(trace, b)
		return nil
	})

	expected := []typedpath.SlashPath{
		"a/b/d",
		"a/b",
		"a/c",
		"a",
	}
	if !reflect.DeepEqual(expected, trace) {
		t.Error(cmp.Diff(expected, trace))
	}
}

func ptrString(s string) *string {
	return &s
}
