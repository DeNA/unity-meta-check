package cstrset

import (
	"github.com/scylladb/go-set/strset"
	"strings"
)

type CaseInsensitiveSet map[string]string

var _ Set = &CaseInsensitiveSet{}

func NewCaseInsensitive(items... string) *CaseInsensitiveSet {
	s := CaseInsensitiveSet(make(map[string]string, len(items)))
	for _, item := range items {
		s.Add(item)
	}
	return &s
}

func NewCaseInsensitiveWithSize(size int) *CaseInsensitiveSet {
	s := CaseInsensitiveSet(make(map[string]string, size))
	return &s
}

func (s CaseInsensitiveSet) Has(e string) bool {
	_, ok := s[strings.ToLower(e)]
	return ok
}

func (s *CaseInsensitiveSet) Add(e string) {
	(*s)[strings.ToLower(e)] = e
}

func (s *CaseInsensitiveSet) Difference(other Set) *strset.Set {
	return difference(s, other.(*CaseInsensitiveSet))
}

func (s CaseInsensitiveSet) Len() int {
	return len(s)
}

func difference(a, b *CaseInsensitiveSet) *strset.Set {
	values := strset.NewWithSize(len(*a))
	strset.Difference(a.keys(), b.keys()).Each(func(item string) bool {
		values.Add((*a)[item])
		return true
	})
	return values
}

func (s *CaseInsensitiveSet) keys() *strset.Set {
	result := strset.NewWithSize(len(*s))
	for key := range *s {
		result.Add(key)
	}
	return result
}
