package autofix

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubMetaCreator(err error) MetaCreator {
	return func(MetaType, typedpath.RawPath) error {
		return err
	}
}
