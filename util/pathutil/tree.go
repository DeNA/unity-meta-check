package pathutil

import "github.com/DeNA/unity-meta-check/util/typedpath"

type PathTreeMap map[typedpath.BaseName]*PathTree
type PathTree struct {
	End bool
	Map PathTreeMap
}

func NewPathTree(paths ...typedpath.SlashPath) *PathTree {
	root := make(map[typedpath.BaseName]*PathTree)
	includeEmpty := false

	for _, path := range paths {
		tree := root
		elements := SplitPathElements(path)
		if len(elements) == 0 {
			includeEmpty = true
		}

		var ok bool
		var treeNode *PathTree
		for _, element := range elements {
			treeNode, ok = tree[element]
			if !ok {
				treeNode = &PathTree{false, map[typedpath.BaseName]*PathTree{}}
				tree[element] = treeNode
			}
			tree = treeNode.Map
		}
		if treeNode != nil {
			treeNode.End = true
		}
	}
	return &PathTree{includeEmpty, root}
}

// Returns whether the pathElements is a member of the tree.
// Notably, returns false if the pathElements point at the tree.
func (t *PathTree) Member(pathElements []typedpath.BaseName) bool {
	if t.End {
		return len(pathElements) > 0
	}
	if len(pathElements) == 0 {
		return false
	}
	child, ok := t.Map[pathElements[0]]
	if !ok {
		return false
	}
	return child.Member(pathElements[1:])
}
