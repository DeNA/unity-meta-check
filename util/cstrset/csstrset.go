package cstrset

import "github.com/scylladb/go-set/strset"

type CaseSensitiveSet struct {
	s *strset.Set
}

var _ Set = &CaseSensitiveSet{}

func NewCaseSensitive(items... string) *CaseSensitiveSet {
	return &CaseSensitiveSet{strset.New(items...)}
}

func NewCaseSensitiveWithSize(size int) *CaseSensitiveSet {
	return &CaseSensitiveSet{strset.NewWithSize(size)}
}

func (s *CaseSensitiveSet) Has(e string) bool {
	return s.s.Has(e)
}

func (s *CaseSensitiveSet) Add(e string) {
	s.s.Add(e)
}

func (s *CaseSensitiveSet) Difference(other Set) *strset.Set {
	return strset.Difference(s.s, other.(*CaseSensitiveSet).s)
}

func (s CaseSensitiveSet) Len() int {
	return s.s.Size()
}
