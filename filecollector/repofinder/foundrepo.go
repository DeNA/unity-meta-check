package repofinder

import "github.com/DeNA/unity-meta-check/util/typedpath"

type FoundRepo struct {
	Type    RepositoryType
	RawPath typedpath.RawPath
}
