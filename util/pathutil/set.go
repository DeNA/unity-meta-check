package pathutil

import (
	"github.com/DeNA/unity-meta-check/util/cstrset"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type PathSet struct {
	set cstrset.Set
}

func NewPathSet(ignoreCase bool, ss ...typedpath.SlashPath) *PathSet {
	items := make([]string, len(ss))
	for i, s := range ss {
		items[i] = string(s)
	}
	if ignoreCase {
		return &PathSet{cstrset.NewCaseInsensitive(items...)}
	}
	return &PathSet{cstrset.NewCaseSensitive(items...)}
}

func NewPathSetWithSize(ignoreCase bool, size int) *PathSet {
	if ignoreCase {
		return &PathSet{cstrset.NewCaseInsensitiveWithSize(size)}
	}
	return &PathSet{cstrset.NewCaseSensitiveWithSize(size)}
}

func (s *PathSet) Add(path typedpath.SlashPath) {
	s.set.Add(string(path))
}

func (s *PathSet) Has(path typedpath.SlashPath) bool {
	return s.set.Has(string(path))
}

func (s *PathSet) Difference(other *PathSet) []typedpath.SlashPath {
	diff := s.set.Difference(other.set)
	result := make([]typedpath.SlashPath, diff.Size())
	i := 0
	diff.Each(func(item string) bool {
		result[i] = typedpath.SlashPath(item)
		i++
		return true
	})
	return result
}

func (s *PathSet) Len() int {
	return s.set.Len()
}
