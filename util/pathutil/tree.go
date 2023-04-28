package pathutil

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"sort"
)

type PathTree[T interface{}] map[typedpath.BaseName]*PathTreeEntry[T]
type PathTreeEntry[T interface{}] struct {
	Value   *T
	Subtree PathTree[T]
}
type PathPair[T interface{}] struct {
	Path  typedpath.SlashPath
	Value T
}
type KeyValuePair[T interface{}] struct {
	BaseName typedpath.BaseName
	Entry    *PathTreeEntry[T]
}

func NewPathTree(paths ...typedpath.SlashPath) PathTree[struct{}] {
	pairs := make([]PathPair[struct{}], len(paths))
	for i, p := range paths {
		pairs[i] = struct {
			Path  typedpath.SlashPath
			Value struct{}
		}{p, struct{}{}}
	}
	return NewPathTreeWithValues(pairs...)
}

func NewPathTreeWithValues[T interface{}](pairs ...PathPair[T]) PathTree[T] {
	root := make(PathTree[T])

	for _, pair := range pairs {
		tree := root
		elements := SplitPathElements(pair.Path)
		if len(elements) == 0 {
			continue
		}

		var ok bool
		var treeNode *PathTreeEntry[T]
		for _, element := range elements {
			treeNode, ok = tree[element]
			if !ok {
				treeNode = &PathTreeEntry[T]{Value: nil, Subtree: map[typedpath.BaseName]*PathTreeEntry[T]{}}
				tree[element] = treeNode
			}
			tree = treeNode.Subtree
		}
		// NOTE: Use copied pointer instead of pointer for the loop variable.
		val := pair.Value
		treeNode.Value = &val
	}
	return root
}

// Member returns whether the pathElements is a member of the tree.
// Notably, returns false if the pathElements point at the tree.
func (t PathTree[T]) Member(pathElements []typedpath.BaseName) bool {
	if len(pathElements) == 0 {
		return false
	}
	child, ok := t[pathElements[0]]
	if !ok {
		return false
	}
	return child.member(pathElements[1:])
}

func (t PathTree[T]) Postorder(f func(typedpath.SlashPath, PathTreeEntry[T]) error) error {
	for _, kv := range sortDict(t) {
		if err := kv.Entry.postorder(".", kv.BaseName, f); err != nil {
			return err
		}
	}
	return nil
}

func (e PathTreeEntry[T]) member(pathElements []typedpath.BaseName) bool {
	if len(pathElements) == 0 {
		return false
	}
	if len(e.Subtree) == 0 {
		return true
	}
	child, ok := e.Subtree[pathElements[0]]
	if !ok {
		return false
	}
	return child.member(pathElements[1:])
}

func (e PathTreeEntry[T]) postorder(relPath typedpath.SlashPath, baseName typedpath.BaseName, f func(typedpath.SlashPath, PathTreeEntry[T]) error) error {
	var path typedpath.SlashPath
	if relPath == "." {
		path = typedpath.SlashPath(baseName)
	} else {
		path = relPath.JoinBaseName(baseName)
	}

	for _, kv := range sortDict(e.Subtree) {
		if err := kv.Entry.postorder(path, kv.BaseName, f); err != nil {
			return err
		}
	}
	return f(path, e)
}

func sortDict[T interface{}](tree PathTree[T]) []KeyValuePair[T] {
	kvs := make([]KeyValuePair[T], 0, len(tree))
	for baseName, entry := range tree {
		kvs = append(kvs, KeyValuePair[T]{baseName, entry})
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].BaseName < kvs[j].BaseName
	})
	return kvs
}
