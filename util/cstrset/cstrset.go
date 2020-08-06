package cstrset

import "github.com/scylladb/go-set/strset"

type Set interface {
	Add(string)
	Has(string) bool
	Difference(Set) *strset.Set
	Len() int
}
