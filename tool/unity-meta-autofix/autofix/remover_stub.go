package autofix

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubMetaRemover(err error) MetaRemover {
	return func(typedpath.RawPath) error {
		return err
	}
}
