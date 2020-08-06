package autofix

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubMetaTypeDetector(result MetaType, err error) MetaTypeDetector {
	return func(typedpath.RawPath) (MetaType, error) {
		return result, err
	}
}
