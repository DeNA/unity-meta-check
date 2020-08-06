package ostestable

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubIsDir(isDir bool, err error) IsDir {
	return func(typedpath.RawPath) (bool, error) {
		return isDir, err
	}
}
